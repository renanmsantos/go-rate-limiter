FROM golang:latest as builder
WORKDIR /app
COPY . . 

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ratelimiter-go .

FROM scratch
COPY --from=builder ./app/ratelimiter-go .
COPY --from=builder ./app/infra ./infra
CMD ["./ratelimiter-go"]