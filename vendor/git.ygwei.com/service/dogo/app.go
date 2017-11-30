package dogo

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"strings"
)

var (
	web_listen_ip   = flag.String("web_listen_ip", "", "web服务监听IP")
	web_listen_port = flag.Int("web_listen_port", 8108, "web服务监听端口")

)

type appServe struct {
	*route
	serveMux *http.ServeMux
}

var app = &appServe{route: route_entity, serveMux: &http.ServeMux{}}

func Route() *route {
	return app.route
}

func Run() {
	var addr = fmt.Sprintf("%s:%d", *web_listen_ip, *web_listen_port)

	Dglog.Infof("Starting .... %s", addr)

	// Register root route handle func
	app.serveMux.HandleFunc("/", app.do)

	// Listening on addr
	http.ListenAndServe(addr, app.serveMux)
}

func (app *appServe) do(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Dglog.Debugf("pattern:%s", r.Form.Get("name"))

	h := route_entity.checkRoute(r)
	requestUri := strings.TrimSuffix(r.URL.Path, path.Ext(r.URL.Path))

	if h == nil {
		Dglog.Errorf("Not found page :%s", r.RequestURI)
	} else {
		ctx := InitContext(w, r)
		if filter_entity.doFilter(requestUri, ctx) {
			h(ctx)
		}
	}
}
