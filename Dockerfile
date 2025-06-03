FROM golang:1.24.0
WORKDIR /app
ENV CONFIG_PATH=./config/config_container.yaml
COPY . .
RUN go mod download
RUN go build -o files cmd/files/main.go
EXPOSE 8000
CMD ["/app/files"]