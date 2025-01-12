FROM golang:1.23.4-alpine3.21 as build

WORKDIR /sapopinguino

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:3.21

COPY --from=build /sapopinguino/main /main

COPY --from=build /sapopinguino/internal/config /config

WORKDIR /

ENTRYPOINT [ "/main" ]

