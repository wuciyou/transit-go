package main

import (
	"git.ygwei.com/service/dogo/config"
	"git.ygwei.com/service/dogo"
	"flag"
	"github.com/wuciyou/transit-go/web"
	"net/http"
)

var (
	web_socket_listen_addr = flag.String("web_socket_listen_addr",":8549", "web socket 监听地址，如只允许本地连接可使用127.0.0.1:8549")
)

func main(){

	config.Parse()
	dogo.Dglog.Info("start transit go ...")
	web.RegisterRouter()
	go dogo.Run()
	dogo.Dglog.Infof("listen web socket :%s",*web_socket_listen_addr)
	http.ListenAndServe(*web_socket_listen_addr,nil)
}

