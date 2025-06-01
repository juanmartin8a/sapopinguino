FROM golang:1.24.3-alpine3.22 as build

WORKDIR /sapopinguino

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:3.22

COPY --from=build /sapopinguino/main /main

COPY --from=build /sapopinguino/internal/config /config

COPY --from=build /sapopinguino/assets /assets

WORKDIR /

ENTRYPOINT [ "/main" ]

