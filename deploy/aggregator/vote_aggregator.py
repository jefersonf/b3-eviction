import os
import time
import redis
import psycopg2
from datetime import datetime
from collections import defaultdict

# --- Configuration ---
# Redis Config
REDIS_HOST = os.getenv("REDIS_ADDR", "redis")
REDIS_PORT = 6379
STREAM_KEY = os.getenv("REDIS_STREAM", "votes")
CONSUMER_GROUP = "aggregator_group"
# Unique name for this worker instance (useful if you scale up replicas)
CONSUMER_NAME = f"worker_{os.getpid()}"

# Postgres Config
PG_HOST = os.getenv("POSTGRES_HOST", "postgres")
PG_USER = os.getenv("POSTGRES_USER", "admin")
PG_PASS = os.getenv("POSTGRES_PASSWORD", "secret")
PG_DB   = os.getenv("POSTGRES_DB", "votings")

def get_pg_connection():
    """Establishes a new connection to Postgres."""
    return psycopg2.connect(
        host=PG_HOST,
        user=PG_USER,
        password=PG_PASS,
        dbname=PG_DB
    )

def main():
    print(f"Starting Aggregator: {CONSUMER_NAME}")
    print(f"Connecting to Redis at {REDIS_HOST}:{REDIS_PORT}...")
    
    # 1. Connect to Redis
    try:
        r = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
        r.ping() # Fail fast if Redis is down
    except Exception as e:
        print(f"Fatal: Could not connect to Redis. {e}")
        return

    # 2. Ensure Consumer Group Exists
    # This allows us to scale workers; Redis distributes messages among them.
    try:
        r.xgroup_create(STREAM_KEY, CONSUMER_GROUP, id="0", mkstream=True)
        print(f"Consumer Group '{CONSUMER_GROUP}' ready.")
    except redis.exceptions.ResponseError as e:
        if "BUSYGROUP" not in str(e):
            raise
        print(f"Consumer Group '{CONSUMER_GROUP}' already exists.")

    # 3. Main Processing Loop
    while True:
        try:
            # A. Read Batch
            # We ask for up to 500 messages. 
            # block=2000 means "sleep for 2 seconds if stream is empty"
            entries = r.xreadgroup(
                groupname=CONSUMER_GROUP,
                consumername=CONSUMER_NAME,
                streams={STREAM_KEY: ">"},
                count=500,
                block=2000
            )

            if not entries:
                continue

            # entries structure: [['stream_name', [('id', {data}), ...]]]
            stream, messages = entries[0]
            if not messages:
                continue

            # B. Aggregation Buffer (The "Micro-Batch" Logic)
            # Key: (bucket_minute, eviction_id, nominee_id) -> Value: count
            db_buffer = defaultdict(int)
            
            # Redis Counter Buffer (For Real-Time Dashboard)
            # Key: "votes:eviction_id:nominee_id" -> Value: count
            redis_counter_buffer = defaultdict(int)
            
            processed_ids = []
            batch_total_votes = 0

            for msg_id, data in messages:
                try:
                    # Extract Data
                    eviction_id = data.get('eviction_id', 'unknown')
                    nominee_id = data.get('nominee_id', 'unknown')
                    # Use event timestamp, default to now if missing
                    ts_raw = int(data.get('timestamp', time.time()))

                    # 1. Prepare DB Aggregate (Minute Bucket)
                    dt = datetime.fromtimestamp(ts_raw)
                    bucket = dt.replace(second=0, microsecond=0)
                    
                    db_key = (bucket, eviction_id, nominee_id)
                    db_buffer[db_key] += 1

                    # 2. Prepare Redis Real-Time Counters
                    # We group these to minimize network calls to Redis later
                    redis_key = f"votes:{eviction_id}:{nominee_id}"
                    redis_counter_buffer[redis_key] += 1
                    
                    processed_ids.append(msg_id)
                    batch_total_votes += 1

                except Exception as e:
                    print(f"Skipping malformed message {msg_id}: {e}")
                    # Ack bad messages so we don't process them forever
                    processed_ids.append(msg_id) 

            # If buffer is empty (e.g., all messages were malformed), skip DB
            if not db_buffer:
                if processed_ids:
                     r.xack(STREAM_KEY, CONSUMER_GROUP, *processed_ids)
                continue

            # C. Write to Infrastructure
            try:
                # 1. Update Real-Time Redis Counters (Instant Dashboard)
                pipeline = r.pipeline()
                pipeline.incrby("global_vote_count", batch_total_votes)
                for key, count in redis_counter_buffer.items():
                    pipeline.incrby(key, count)
                pipeline.execute()

                # 2. Upsert to Postgres (Historical Analytics)
                with get_pg_connection() as conn:
                    with conn.cursor() as cur:
                        for (bucket, eviction, nominee), count in db_buffer.items():
                            query = """
                                INSERT INTO votes_minutely (bucket_minute, eviction_id, nominee_id, votes)
                                VALUES (%s, %s, %s, %s)
                                ON CONFLICT (bucket_minute, eviction_id, nominee_id)
                                DO UPDATE SET votes = votes_minutely.votes + EXCLUDED.votes;
                            """
                            cur.execute(query, (bucket, eviction, nominee, count))
                    conn.commit()
                
                # 3. Acknowledge Batch (Only if DB write succeeded)
                # This ensures "At-Least-Once" processing.
                if processed_ids:
                    r.xack(STREAM_KEY, CONSUMER_GROUP, *processed_ids)
                
                print(f"Synced {len(messages)} votes. (DB Rows: {len(db_buffer)} | Global Count +{batch_total_votes})")

            except psycopg2.Error as db_err:
                print(f"Database Error: {db_err}. Retrying batch in 2s...")
                time.sleep(2)
                # We do NOT ack here. The loop will re-read these messages.
            
            except redis.RedisError as redis_err:
                print(f"Redis Write Error: {redis_err}. Retrying batch in 2s...")
                time.sleep(2)

        except Exception as e:
            print(f"Unexpected Worker Error: {e}")
            time.sleep(2)

if __name__ == "__main__":
    # Wait for DB to be ready (Naive check, can be improved with a retry loop)
    time.sleep(5) 
    main()