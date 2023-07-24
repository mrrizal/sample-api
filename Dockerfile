FROM golang:1.20-alpine as builder

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main


FROM scratch

COPY --from=builder /src/main .

CMD ["./main"]
