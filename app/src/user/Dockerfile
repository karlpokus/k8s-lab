FROM golang:alpine
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./user

EXPOSE 9051

ENTRYPOINT ["/src/user"]
