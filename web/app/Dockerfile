FROM python:3.14-slim

WORKDIR /voting-app

COPY . .

RUN python -m pip install requests

CMD [ "python", "-m", "http.server", "4000" ]