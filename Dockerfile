FROM golang:1.18

RUN apt update && \ 
    apt install -y \
    wget unzip tree

WORKDIR /app

COPY . .

RUN go build -o exercism-events

CMD [ "./exercism-events", "-sources=mentoring", "-consumers=database" ]
