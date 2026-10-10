FROM golang:1.26 AS build

WORKDIR /backend

COPY --exclude=client . .

RUN mkdir bin

RUN CGO_ENABLED=0 go build -o ./bin/around ./main.go

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    tor \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /backend/bin/around ./bin/around
COPY --from=build /backend/config ./config

COPY start.sh ./start.sh
RUN chmod +x ./start.sh

ENTRYPOINT ["./start.sh"]