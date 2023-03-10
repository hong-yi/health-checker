FROM golang:latest AS builder
RUN go version
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /healthchecker
CMD ["/healthchecker"]
EXPOSE 8080

FROM amazonlinux:latest
COPY --from=builder /healthchecker .
EXPOSE 8080
CMD ["/healthchecker"]