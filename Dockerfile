FROM golang:1.17-alpine

WORKDIR /buildsource
COPY . /buildsource

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /usr/local/bin/inform

CMD ["inform", "--in-cluster"]
