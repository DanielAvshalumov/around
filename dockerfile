FROM golang:1.26 AS build

WORKDIR /backend

RUN apt-get update && apt-get install -y \
tor \
&& rm -rf /var/lib/apt/lists/*

COPY --exclude=client . .

RUN mkdir bin

RUN CGO_ENABLED=0 go build -o ./bin/around ./main.go

FROM scratch

COPY --from=build backend/bin/around ./bin/around
COPY --from=build backend/config ./config

COPY --from=build /usr/bin/tor /usr/bin/tor
COPY --from=build /etc/tor /etc/tor

ENTRYPOINT ["systemctl", "start", "tor"]

CMD ["/bin/around"]