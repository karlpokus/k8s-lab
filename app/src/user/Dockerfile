FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./user
RUN go build -o ./user

#FROM scratch
#WORKDIR /src
#COPY --from=builder /src/user ./user

EXPOSE 9051

ENTRYPOINT ["/src/user"]
