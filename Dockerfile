FROM golang:1.23.4 as build

WORKDIR /sapopinguino

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /sapopinguino/main ./main
ENTRYPOINT [ "./main" ]

