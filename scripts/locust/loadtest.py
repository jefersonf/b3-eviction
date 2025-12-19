# locust/locustfile.py
import random
from locust import FastHttpUser, task, between

class VotingUser(FastHttpUser):
    # Wait between 0.5 and 1.5 seconds (avg 1s). 
    # With 1000 users, this averages ~1000 req/sec.
    wait_time = between(0.5, 1.5)

    @task
    def vote(self):
        nominee = random.choice(["Ana", "Bob"])
        payload = {
            "nominee_id": nominee,
            "eviction_id": "bbb-load-test-id"
        }
        self.client.post("/api/vote", json=payload, name="/vote")