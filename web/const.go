package web

import(
	"errors"
)

var (
	ERR_CANNOT_FIND_DATA = errors.New("找不到数据")
	ERR_ALREADY_DATA = errors.New("数据已经存在")
)

type sessionStatusCode uint32

const (
	SESSION_STATUS_SUCCESS sessionStatusCode = 1 << iota
	SESSION_STATUS_ALREADY_EXIST
	SESSION_STATUS_SYSTEM_ERROR
	SESSION_STATUS_AUTH_TIMEOUT
	SESSION_STATUS_AUTH_FAIL
	SESSION_STATUS_MULTIPLE_AUTH_FAIL
	SESSION_STATUS_SYSTEM_CLOSE
)