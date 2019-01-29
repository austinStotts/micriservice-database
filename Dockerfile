FROM golang:latest
RUN mkdir -p /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get github.com/lib/pq
RUN go build server.go
COPY . .
EXPOSE 80
ENTRYPOINT ["./server"]
