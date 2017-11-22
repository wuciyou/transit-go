package main

import (
	"errors"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"strings"
	"sync"
)

type transit struct {
	sessions map[string][]sockjs.Session
	keys     map[string]string
	sync.Mutex
}

func NewTransit() *transit {
	t := &transit{
		sessions: make(map[string][]sockjs.Session),
		keys:     make(map[string]string),
	}
	return t
}

func (this *transit) transitByKeyString(keys string, msg string) {
	keySlice := strings.Split(keys, ",")
	go this.transitByKeykyeSlice(keySlice, msg)
}

func (this *transit) transitByKeykyeSlice(keys []string, msg string) {
	for _, v := range keys {
		this.sendMsg(v, msg)
	}
}

func (this *transit) sendMsg(key string, msg string) error {
	if sessions, ok := this.sessions[key]; ok {

		for _, session := range sessions {
			err := session.Send(msg)
			if err != nil {
				log.Printf("Send msg fail, err:%+v", err)
				return err
			}
		}
		return nil
	} else {
		return errors.New("Can't find session key:" + key)
	}
}

func (this *transit) addSession(key string, session sockjs.Session) error {
	this.Lock()
	this.sessions[key] = append(this.sessions[key], session)
	this.keys[session.ID()] = key
	this.Unlock()
	return nil
}

func (this *transit) closeSessionByKey(keys string) {
	go func(keys string) {
		keySlice := strings.Split(keys, ",")
		for _, key := range keySlice {
			if sessions, exsit := this.sessions[key]; exsit {
				for _, session := range sessions {
					this.closeSession(session)
				}
			}
		}
	}(keys)
}

func (this *transit) closeSession(session sockjs.Session) error {
	if key, ok := this.keys[session.ID()]; ok {
		this.Lock()
		log.Printf("Close session form key[%s] session.id[%s]", key, session.ID())

		for i, s := range this.sessions[key] {
			if s.ID() == session.ID() {
				this.sessions[key][i] = nil
				this.sessions[key] = append(this.sessions[key][:i], this.sessions[key][i+1:]...)

			}
		}

		if len(this.sessions[key]) <= 0 {
			delete(this.sessions, key)
		}
		delete(this.keys, session.ID())
		session.Close(200, "")
		this.Unlock()
		return nil
	} else {
		return errors.New("Can't find session Keys from session.ID:" + session.ID())
	}
}
