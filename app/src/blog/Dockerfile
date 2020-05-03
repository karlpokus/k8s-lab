FROM golang:alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o ./blog

EXPOSE 9052

ENTRYPOINT ["/src/blog"]
