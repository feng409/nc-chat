# nc-chat

基于`go`的`net`包开发的TCP聊天服务器。
使用 `nc` 命令作为客户端

## Usage
```shell
go run main.go
2020/02/01 17:57:43 listen on  0.0.0.0:8888

# another terminel window
nc 127.0.0.1 8888
hey guy, please tell me your name
$ mike

2020-02-01 17:58:24  [system]: welcome mike
$ hello world

2020-02-01 17:58:59    [mike]: hello world
$
```