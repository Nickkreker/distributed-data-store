FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /lesson1

EXPOSE 8080

CMD ["/lesson1"]