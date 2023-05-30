FROM golang:1.19

WORKDIR /usr/src/http-benchmarking

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build .

CMD ["/usr/src/http-benchmarking/http-benchmarking"]
