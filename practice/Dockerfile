
FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod init app && go build -o main .

CMD ["./main"]
