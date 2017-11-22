package main

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"os"
	"sync"
	"time"
)

type authoirty struct {
	keys map[string]string
	sync.Mutex
}

func NewAuthoirty() *authoirty {
	return &authoirty{keys: make(map[string]string)}
}

func (this *authoirty) getKey() string {
	f, err := os.Open("/dev/random")
	var key string
	if err != nil {
		log.Panicf("Can't open /dev/random file, err:%+v", err)
	} else {
		b := make([]byte, 1024)
		f.Read(b)
		var md5Ctx = md5.New()
		md5Ctx.Write(b)
		key = hex.EncodeToString(md5Ctx.Sum(nil))
		log.Printf("auto generate key[%s]", key)
	}
	if this.setKey(key) {
		return key
	} else {
		return this.getKey()
	}
}

func (this *authoirty) setKey(key string) bool {
	var isOk bool
	this.Lock()
	log.Printf("Add key to authoirty from key[%s]", key)
	if _, ok := this.keys[key]; ok {
		isOk = false
	} else {
		this.keys[key] = key
		isOk = true
	}
	this.Unlock()
	return isOk
}

func (this *authoirty) check(key string) bool {
	if _, ok := this.keys[key]; ok {
		this.delKey(key)
		log.Printf("Success key authoirty from key[%s]", key)
		return true
	} else {
		log.Printf("Wrong key authoirty from key[%s]", key)
		return false
	}
}

func (this *authoirty) delKey(key string) bool {

	if _, ok := this.keys[key]; ok {
		this.Lock()
		log.Printf("Delete key[%s] success", key)
		delete(this.keys, key)
		this.Unlock()
		return true
	} else {
		log.Printf("Delete key fail, Because key[%s] not exist", key)
	}
	return false
}

func (this *authoirty) checkPool(session sockjs.Session, callBack func(sessionKey string, session sockjs.Session)) {
	var sessionKey = make(chan string)
	go func() {
		for {
			select {
			// 验证成功
			case key := <-sessionKey:
				callBack(key, session)
				return
			// 超时
			case <-time.After(10 * time.Second):
				log.Printf("Check  session.id[%s] time out for authority", session.ID())
				log.Printf("Close session.id[%s] for session", session.ID())
				session.Close(500, "time out")
				return
			}
		}
	}()

	for {
		if key, err := session.Recv(); err == nil {
			if this.check(key) {
				log.Printf("Check authority success for key[%s] session.id[%s]", key, session.ID())
				sessionKey <- key
				break
			} else {
				log.Printf("Check authority fail, Wrong key[%s]", key)
				session.Close(403, "Check authority fail")
			}
		} else {
			break
		}
	}
}
