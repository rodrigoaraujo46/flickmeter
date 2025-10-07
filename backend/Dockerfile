FROM golang:latest
WORKDIR /app
COPY . .
RUN go install github.com/air-verse/air@latest
EXPOSE 1323
ENTRYPOINT ["air"]
