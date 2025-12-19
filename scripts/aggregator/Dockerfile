FROM python:3.14-slim

WORKDIR /voting-aggregator

RUN pip install redis psycopg2-binary

COPY ./vote_aggregator.py /voting-aggregator/vote_aggregator.py

CMD [ "python", "-u", "vote_aggregator.py" ]