FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x postgres.sh

RUN go mod download
RUN go build -o bash_server ./cmd/main.go

CMD ["./bash_server"]