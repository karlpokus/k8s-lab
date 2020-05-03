FROM golang:alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o ./user

EXPOSE 9051

ENTRYPOINT ["/src/user"]
