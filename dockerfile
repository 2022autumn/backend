# 考虑暂时不使用docker部署。先跑通再说吧
FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/IShare
COPY . $GOPATH/src/IShare
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./IShare"]