package main

import (
	"flag"
	"fmt"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
	"strings"
)

var auth *authoirty
var transitChan *transit

func mains() {

	host := flag.String("host", "", "listen host address")
	port := flag.Int("port", 8549, "listen posrt")
	flag.Parse()
	auth = NewAuthoirty()
	transitChan = NewTransit()

	http.Handle("/wc/", sockjs.NewHandler("/wc", sockjs.DefaultOptions, authorityHandler))
	http.HandleFunc("/key", getKey)

	http.HandleFunc("/sendMsg", sendMsg)
	http.HandleFunc("/closeSession", closeSession)
	http.Handle("/", http.FileServer(http.Dir("web/")))
	listenHost := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Server started on " + listenHost)
	log.Fatal(http.ListenAndServe(listenHost, nil))
}

func getKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("key[%s]", r.Form.Get("key"))
	var key = strings.TrimSpace(r.Form.Get("key"))
	log.Printf("TrimSpace key[%s]", key)
	if key != "" {
		log.Printf("Use Form key[%s]", key)
		if !auth.setKey(key) {
			key = auth.getKey()
		}
	} else {
		key = auth.getKey()
	}

	log.Printf("Share key[%s]", key)

	w.Write([]byte(key))
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var key = r.Form.Get("key")
	var msg = r.Form.Get("msg")
	log.Printf("Have new task send message for keys[%s] message[%s]", key, msg)
	transitChan.transitByKeyString(key, msg)
	w.Write([]byte("ok"))
}

func closeSession(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var key = strings.TrimSpace(r.Form.Get("key"))
	log.Printf("Have new task close  keys[%s] for session", key)
	transitChan.closeSessionByKey(key)
	w.Write([]byte("ok"))
}

func authorityHandler(session sockjs.Session) {
	auth.checkPool(session, func(key string, session sockjs.Session) {
		log.Printf("Have new session key[%s] session.id[%s]", key, session.ID())
		transitChan.addSession(key, session)
		for {
			if msg, err := session.Recv(); err == nil {
				log.Printf("New message from key[%s] session.id[%s] msg:%s", key, session.ID(), msg)
				continue
			} else {
				log.Println(err)
				transitChan.closeSession(session)
			}
			break
		}
	})
}
