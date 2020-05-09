FROM golang:alpine
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./gw

EXPOSE 9050

ENTRYPOINT ["/src/gw"]
