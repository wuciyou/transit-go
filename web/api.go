package web

import (
	"git.ygwei.com/service/dogo"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"sync"
	"time"
	"net/http"
	"strings"
	"os"
)

type StatusCode string

const (
	SUCCESS StatusCode = "SUCCESS"
	REDIRECT_MESSAGE StatusCode = "REDIRECT_MESSAGE"
	PARAM_ERROR StatusCode = "PARAM_ERROR"
	DATA_EXSIT StatusCode = "DATA_EXSIT"
	SYSTEM_ERROR StatusCode = "SYSTEM_ERROR"
	AUTH_FAIL StatusCode = "AUTH_FAIL"
	AUTH_SUCCESS StatusCode = "AUTH_SUCCESS"
)

type responseEntity struct {
	Status StatusCode
	Message string
	Data interface{}
}

type waitConnectUser  struct{
	sync.Mutex
	waitConnectUserQueue map[string]bool
}

/**
 * 获取响应实例
 */
func newResponseEntity( status StatusCode,  data interface {}, messages ...string) *responseEntity {
	var message string
	if len(messages) > 0{
		message = messages[0]
	}
	return &responseEntity{Status:status,Data:data,Message:message}
}

/**
 * 获取响应成功实例
 */
func newSuccessResponseEntity( data interface {}) *responseEntity {
	return &responseEntity{Status:SUCCESS,Data:data,Message:"请求成功"}
}

/**
 * 注册普通用户
 */
func(this *waitConnectUser) registerUser(ctx *dogo.Context){
	sessionKey := ctx.Get("connectKey")
	var response *responseEntity
	if _,ok := this.waitConnectUserQueue[sessionKey];!ok{
		this.Lock()
		if _,exist := this.waitConnectUserQueue[sessionKey];!exist{
			this.waitConnectUserQueue[sessionKey] = false ;
		}
		this.Unlock()
		response = newSuccessResponseEntity(nil)
	}else{
		response = newResponseEntity(DATA_EXSIT,nil,"数据已经存在")
	}
	ctx.W.Json(response)
}


/**
 * 关闭指定用户连接
 */
func(this *waitConnectUser) closeUser(ctx *dogo.Context){

}

/**
 * 关闭指定用户连接
 */
func(this *waitConnectUser) exitApp(ctx *dogo.Context){
	ctx.W.Write([]byte("exit ok"))
	ctx.W.Send()
	os.Exit(1)
}

/**
 * 发送普通消息
 */
func(this *waitConnectUser) sendMessage(ctx *dogo.Context){
	connectKeys := ctx.Get("connectKey")
	message := ctx.Get("message")
	var successNum = 0
	var errStr = ""
	for _,connectKey := range strings.Split(connectKeys,","){
		if err := DefaultSessionManager.Send(connectKey,message); err == nil{
			successNum ++
		}else{
			errStr += err.Error()
		}
	}
	ctx.W.Json(newResponseEntity(SUCCESS,successNum,errStr))
}

func (this *waitConnectUser) webSocketHandle(session sockjs.Session){
	var checkFailNumber = 0
	var sessionConnectKey = make(chan string)
	var sessionClose = make(chan bool)
	var messageChan = make(chan  string,1024)
	var isAuthoritySuccess = false
	go func(){
		var connectKey string
		select{
			case connectKey = <- sessionConnectKey:
				us := &userSession{
					Session:session,
					connectKey:connectKey,
				}
				DefaultSessionManager.addSession(connectKey,us)
				break

		case <-sessionClose:
		case <-time.After(10*time.Second):
				session.Close(uint32(SESSION_STATUS_AUTH_TIMEOUT),"认证超时")
				return
		}

		for{
			select{
				case msg := <-messageChan:
					sendEvent(EVENT_NEW_MESSAGE,msg)
					break
				case  <-sessionClose:
				sendEvent(EVENT_CLOSE_SESSION,connectKey)
				return
			}
		}
	}()

	for{
		if message,err := session.Recv();err == nil{

			if !isAuthoritySuccess{
				connectKey := message
				if isConnect, exist := this.waitConnectUserQueue[connectKey]; exist {
					if !isConnect {
						this.Lock()
						delete(this.waitConnectUserQueue,connectKey)
						this.Unlock()
						sessionConnectKey <- connectKey
						isAuthoritySuccess = true
						data,_ := ToJson(newResponseEntity(AUTH_SUCCESS,connectKey,"认证成功"))
						session.Send(string(data))
					}else{
						sessionClose <- true
						session.Close(uint32(SESSION_STATUS_ALREADY_EXIST),"请勿重复连接")
						break
					}
				} else{
					checkFailNumber ++
					data,_ := ToJson(newResponseEntity(AUTH_FAIL,connectKey,"认证失败"))
					session.Send(string(data))
				}
				// 只有三次密码认证机会
				if checkFailNumber >= 3{
					sessionClose <- true
					session.Close(uint32(SESSION_STATUS_MULTIPLE_AUTH_FAIL),"只有3次认证机会")
					break
				}

			}else{
				// 收到普通消息
				messageChan <- message
			}

		}else{
			sessionClose <- true
			session.Close(uint32(SESSION_STATUS_SYSTEM_ERROR),err.Error())
			break
		}
	}
}

func RegisterRouter(){

	wu := &waitConnectUser{waitConnectUserQueue:make(map[string]bool)}
	r := dogo.Route()
	r.Get("/registerUser",wu.registerUser)
	r.Get("/sendMessage",wu.sendMessage)
	r.Post("/registerUser",wu.registerUser)
	r.Post("/sendMessage",wu.sendMessage)
	r.Get("/exitApp",wu.exitApp)

	http.Handle("/", http.FileServer(http.Dir("web/")))
	http.Handle("/ws/", sockjs.NewHandler("/ws", sockjs.DefaultOptions, wu.webSocketHandle))

}