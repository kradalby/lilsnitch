FROM golang:1.17.5 as builder
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .


FROM alpine:latest  

RUN apk --no-cache add ca-certificates
WORKDIR /

ENV GIN_MODE=release
EXPOSE 8080/tcp

COPY --from=builder /app/app .
CMD ["./app"]
