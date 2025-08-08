FROM alpine:latest

WORKDIR /app

COPY loggerApp /app/loggerApp

CMD ["./loggerApp"]