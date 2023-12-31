FROM golang:1.21

WORKDIR /app

RUN apt-get update && \
    apt-get install -y \
	poppler-utils \
    imagemagick

COPY . /app

RUN go build -o /pput /app/cmd/pput

CMD ["/pput"]
