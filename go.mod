module github.com/fregie/gotool

go 1.13

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/brucespang/go-tcpinfo v0.2.0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200518175753-732df43cf85b
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	google.golang.org/grpc v1.29.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/shadowsocks/go-shadowsocks2 => github.com/geewan-rd/go-shadowsocks2 v1.0.1