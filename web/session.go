package web

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"sync"
	"git.ygwei.com/service/dogo"
)

type userSession struct{
	sockjs.Session
	connectKey string
	groupNames []string
}

type sessionManager struct{
	session map[string]*userSession
	sync.Mutex
}

var DefaultSessionManager = defaultSessionManager()

func defaultSessionManager()*sessionManager{
	sm := &sessionManager{session:make(map[string]*userSession)}

	RegisterEventListener(EVENT_CLOSE_SESSION,sm.listenClose)

	return sm
}

func (sm *sessionManager) listenClose(en eventName, data interface{}){
	if connectKey,ok := data.(string); ok{
		sm.delSession(connectKey)
	}
}

/**
 * 发送消息
 */
func (sm *sessionManager) Send(connectKey string, message string) error {
	message = ToJsonStr(newResponseEntity(REDIRECT_MESSAGE,message,"消息转发"))
	if us,ok := sm.session[connectKey]; ok {
		if err := us.Send(message); err != nil{
			dogo.Dglog.Errorf("发送消息失败，connectKey:%s, message:%s, err:%+v",connectKey,message,err)
			sm.delSession(connectKey)
			return err
		}else{
			dogo.Dglog.Debugf("发送消息成功，connectKey:%s, message:%s",connectKey,message)
		}
	}else{
		dogo.Dglog.Warningf("无法通过connectKey找到session对象，connectKey:%s, message:%s",connectKey,message)
		return ERR_CANNOT_FIND_DATA
	}
	return nil
}


/**
 * 向sessionManager注册用户session
 */
func (sm *sessionManager) addSession(connectKey string, session *userSession) error {

	if us,ok := sm.session[connectKey]; ok{
		if  us.Session == nil{
			us.Session = session
		}else{
			dogo.Dglog.Errorf("sessionManager里不允许插入两个相同的connectKey:%s,userSession:%+v,err:%+v",connectKey,us,ERR_ALREADY_DATA)
		}
		return ERR_ALREADY_DATA
	}else{
		if _,ok := sm.session[connectKey]; !ok{
			sm.Lock()
			sm.session[connectKey] = session
			sm.Unlock()
		}
	}
	return nil
}

/**
 * 将userSessionr从sessionManager中删除
 */
func (sm *sessionManager) delSession(connectKey string) {
	sm.Lock()
	defer sm.Unlock()
	if userSession,ok := sm.session[connectKey]; ok{
		dogo.Dglog.Debugf("删除userSession form key:%s",connectKey)
		delete(sm.session,connectKey)
		userSession.Close(uint32(SESSION_STATUS_SYSTEM_CLOSE),"删除")
	}
}
