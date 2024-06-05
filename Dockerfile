FROM golang:1.22.3

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o pipedrive-test-api

EXPOSE 8080

CMD ["/app/pipedrive-test-api"]



