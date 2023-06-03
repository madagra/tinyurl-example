FROM golang:1.20

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY ./tinyurl ./

RUN go build -v -o tinyurl

EXPOSE 3000
CMD ["/app/tinyurl"]
