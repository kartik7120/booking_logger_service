FROM alpine:latest

WORKDIR /app

COPY loggerApp /app/loggerApp

RUN chmod +x loggerApp

CMD ["./loggerApp"]