FROM golang:1.25.0-alpine3.22 as build

WORKDIR /sapopinguino

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -tags prod -ldflags="-s -w" -o ./bin/main ./main.go 

FROM alpine:3.22

WORKDIR /sapopinguino

COPY --from=build /sapopinguino/bin/main ./main

COPY --from=build /sapopinguino/config/config.prod.yml ./config/config.prod.yml

ENTRYPOINT [ "./main" ]
