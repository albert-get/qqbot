FROM golang:alpine3.16

WORKDIR  /qqbot
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN  go mod download
RUN go build
CMD ["./qqbot"]
