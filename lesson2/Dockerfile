FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go ./
ADD transaction transaction

RUN CGO_ENABLED=0 GOOS=linux go build -o /lesson2

EXPOSE 8080

CMD ["/lesson2"]