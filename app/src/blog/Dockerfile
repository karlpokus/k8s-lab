FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./blog
RUN go build -o ./blog

#FROM scratch
#WORKDIR /src
#COPY --from=builder /src/blog ./blog

EXPOSE 9050

ENTRYPOINT ["/src/blog"]
