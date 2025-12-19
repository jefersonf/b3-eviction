FROM python:3.14-slim

WORKDIR /app

# Install dependencies
RUN pip install redis psycopg2-binary

# Copy the script
COPY ./vote_aggregator.py /app/vote_aggregator.py

# Run unbuffered so logs show up immediately
CMD [ "python", "-u", "vote_aggregator.py" ]