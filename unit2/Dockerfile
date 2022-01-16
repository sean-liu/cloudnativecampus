FROM golang:latest

WORKDIR /app
#将服务器的go工程代码加入到docker容器中
ADD . /app
#go构建可执行文件
RUN go build .
#暴露端口
EXPOSE 8080
#最终运行docker的命令
ENTRYPOINT  ["./app"]