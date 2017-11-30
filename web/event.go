package web

import "sync"

type eventName int

const (
	EVENT_NEW_SESSION eventName = 1 << iota
	EVENT_CLOSE_SESSION
	EVENT_NEW_MESSAGE
)
type listenerHandle func(eName eventName, data interface{})

var DefaultEvent = &event{listenerHandleMap:make(map[eventName][]listenerHandle)}

type event struct{
	sync.Mutex
	listenerHandleMap map[eventName] []listenerHandle
}

func (e *event) RegisterListener(eName eventName, h listenerHandle){
	e.Lock()
	e.listenerHandleMap[eName] = append(e.listenerHandleMap[eName],h)
	e.Unlock()
}

func (e *event) sendEvent(eName eventName, data interface{}){
	for _,h := range e.listenerHandleMap[eName]{
		go h(eName,data)
	}
}

func RegisterEventListener(eName eventName, h listenerHandle){
	DefaultEvent.RegisterListener(eName,h)
}

func sendEvent(eName eventName, data interface{}){
	DefaultEvent.sendEvent(eName,data)
}