FROM python:3.11.5-alpine3.18

WORKDIR /app

RUN pip install gunicorn==21.2.0

COPY requirements.txt .
RUN pip install -r requirements.txt


RUN mkdir instance files
COPY flags.csv files/5
COPY static static/
COPY templates templates/
COPY *.py .
COPY entrypoint.sh . 

RUN chown -R nobody:nobody .
USER nobody

ENTRYPOINT [ "sh", "entrypoint.sh" ]
