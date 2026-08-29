FROM golang:1.26 AS build

WORKDIR /backend

COPY --exclude=client . .

RUN mkdir bin

RUN CGO_ENABLED=0 go build -o ./bin/around ./main.go

FROM scratch

COPY --from=build backend/bin/around ./bin/around
COPY --from=build backend/config ./config

CMD ["/bin/around"]