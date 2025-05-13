FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod init go-mysql-app
RUN go get github.com/go-sql-driver/mysql

RUN go build -o main .

CMD ["./main"]
