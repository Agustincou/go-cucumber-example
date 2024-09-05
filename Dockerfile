FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./

RUN ls -la main

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

RUN chmod +x main

EXPOSE 8080

CMD [ "./main" ]
