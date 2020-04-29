FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./gw
RUN go build -o ./gw

#FROM scratch
#WORKDIR /src
#COPY --from=builder /src/gw ./gw

EXPOSE 9050

ENTRYPOINT ["/src/gw"]
