FROM golang:latest
RUN mkdir -p /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get github.com/lib/pq
EXPOSE 80
RUN go build srver.go
COPY . . 

ENTRYPOINT ["/server"]