FROM golang:alpine
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./blog

EXPOSE 9052

ENTRYPOINT ["/src/blog"]
