FROM golang:alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o ./gw

EXPOSE 9050

ENTRYPOINT ["/src/gw"]
