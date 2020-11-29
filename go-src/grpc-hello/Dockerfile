FROM golang:1.15-buster

# Create and change to the app directory.
WORKDIR /app

COPY . .
RUN cp -pRf helloworld $(go env GOROOT)/src
RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN cd greeter_server && go mod init main && go mod download && go build -o /app/server main.go
ENTRYPOINT ["/app/server"]
