FROM golang:latest
LABEL authors="rahmanrafi"

WORKDIR /usr/src/app

ARG INTERNAL_PORT=8080
ENV PORT ${INTERNAL_PORT}
EXPOSE ${INTERNAL_PORT}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/app

CMD ["app"]
