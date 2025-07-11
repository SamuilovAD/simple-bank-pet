# Build stage
FROM golang:1.24.0-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
#Run stage (It helps to reduce the size of image)
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration
COPY start.sh .
COPY wait-for.sh .
EXPOSE 8080

CMD [ "/app/main" ]
ENTRYPOINT ["/app/start.sh"]