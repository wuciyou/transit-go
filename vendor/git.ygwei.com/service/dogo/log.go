package dogo

import (
	"fmt"
	"log"
	"runtime/debug"
	"flag"
	"sync"
	"os"
	"time"
	"path"
	//"bytes"
)

type RunLevel string

var (
	RUN_INFO    RunLevel = "INFO"
	RUN_WARNING RunLevel = "WARN"
	RUN_DEBUG   RunLevel = "DEBUG"
	RUN_ERROR   RunLevel = "ERROR"

	RUN_INFO_FORMAT    = fmt.Sprintf("%c[0,0,%dm %-7s", 0x1B, 32, RUN_INFO)
	RUN_WARNING_FORMAT = fmt.Sprintf("%c[0,0,%dm %-7s", 0x1B, 35, RUN_WARNING)
	RUN_DEBUG_FORMAT   = fmt.Sprintf("%c[0,0,%dm %-7s", 0x1B, 36, RUN_DEBUG)
	RUN_ERROR_FORMAT   = fmt.Sprintf("%c[0,0,%dm %-7s", 0x1B, 31, RUN_ERROR)
	_CLEAR_COLOR       = fmt.Sprintf("%c[0m", 0x1B)
	log_output_dir = flag.String("log_output_dir","","日志输出目录")
	isInitLog = false
)

var Dglog = &dglog{runLevel: RUN_DEBUG}

func initLog() {
	log.Printf("initLog")
	SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	if (*log_output_dir) != ""{
		err := os.Mkdir(*log_output_dir,os.ModePerm)
		if err != nil && !os.IsExist(err){
			panic(err)
		}

		fileLogger := &fileLog{m:sync.Mutex{},o:make(map[string]*os.File)}
		AddLogger("file",log.New(fileLogger, "",log.LstdFlags))
	}

	AddLogger("default",log.New(os.Stderr, "",log.LstdFlags))
}

func (l *dglog) AddLogger(name string, logger *log.Logger) {
	l.m.Lock()
	//log.SetOutput(w)
	if l.logs == nil{
		l.logs = make(map[string]*log.Logger)
	}
	if _,exits := l.logs[name]; !exits{
		l.logs[name] = logger
	}
	l.m.Unlock()
}

type fileLog struct{
	o map[string]*os.File
	m sync.Mutex
}

func(f *fileLog) Write(b []byte)(n int, err error){
	return f.getOutput("output").Write(b)
}

func(f *fileLog) getOutput(level  string)*os.File{
	oname := fmt.Sprintf("%s_%s.log",level,time.Now().Format("2006_01_02"))
	if o,ok := f.o[level];!ok{
		f.m.Lock()

		if checkOutPut,ok := f.o[level];ok{
			f.m.Unlock()
			return checkOutPut;
		}
		oFileName := path.Clean(fmt.Sprintf("%s/%s",*log_output_dir,oname))
		oFile ,err := os.OpenFile(oFileName,os.O_CREATE | os.O_WRONLY | os.O_APPEND,os.ModePerm)
		if err != nil {
			panic(err)
		}
		f.o[level] = oFile
		f.m.Unlock()
		return oFile
	}else{
		if o.Name() != oname{
			f.m.Lock()
			o.Close()
			delete(f.o,level)
			f.m.Unlock()
			return f.getOutput(level);
		}
		return o
	}

	return nil
}


type dglog struct {
	runLevel        RunLevel
	logs map[string]*log.Logger
	m sync.Mutex
}

func (l *dglog) walk(f func( logger *log.Logger)){
	if !isInitLog{
		isInitLog = true
		initLog()
	}
	for _,logg :=range l.logs{
		f(logg)
	}
}

func (l *dglog) Output(s string,level RunLevel) {
	l.walk(func(logger *log.Logger){
		logger.Output(3,  string(level)+" "+s)
	})
}

func (l *dglog) SetPrefix(prefix string) {
	l.walk(func(logger *log.Logger){
		logger.SetPrefix(prefix)
	})
}
func (l *dglog) SetFlags(flag int) {
	l.walk(func(logger *log.Logger){
		logger.SetFlags(flag)
	})
}

func (l *dglog) Info(v ...interface{}) {
	l.SetPrefix(RUN_INFO_FORMAT)
	v = append(v, _CLEAR_COLOR)
	l.Output(fmt.Sprint(v...),RUN_INFO)
}

func (l *dglog) Infof(format string, v ...interface{}) {
	l.SetPrefix(RUN_INFO_FORMAT)
	l.Output(fmt.Sprintf(format+_CLEAR_COLOR, v...),RUN_INFO)
}

func (l *dglog) Debug(v ...interface{}) {

	if l.runLevel == RUN_DEBUG {
		l.SetPrefix(RUN_DEBUG_FORMAT)
		l.Output(fmt.Sprint(v...) + _CLEAR_COLOR,RUN_DEBUG)
	}
}
func (l *dglog) Level(level RunLevel) bool {
	return l.runLevel == level
}

func (l *dglog) Debugf(format string, v ...interface{}) {

	if l.runLevel == RUN_DEBUG {
		l.SetPrefix(RUN_DEBUG_FORMAT)
		l.Output(fmt.Sprintf(format, v...) + _CLEAR_COLOR,RUN_DEBUG)
	}
}

func (l *dglog) Warning(v ...interface{}) {
	if l.runLevel == RUN_DEBUG {
		l.SetPrefix(RUN_WARNING_FORMAT)
		l.Output(fmt.Sprint(v...) + _CLEAR_COLOR,RUN_WARNING)
	}
}

func (l *dglog) Warningf(format string, v ...interface{}) {
	if l.runLevel == RUN_DEBUG {
		l.SetPrefix(RUN_WARNING_FORMAT)
		l.Output(fmt.Sprintf(format, v...) + _CLEAR_COLOR,RUN_WARNING)
	}
}

func (l *dglog) Error(v ...interface{}) {
	l.SetPrefix(RUN_ERROR_FORMAT)
	s := fmt.Sprint(v...) + string(debug.Stack()) + _CLEAR_COLOR
	l.Output(s,RUN_ERROR)
}

func (l *dglog) Errorf(format string, v ...interface{}) {
	l.SetPrefix(RUN_ERROR_FORMAT)
	s := fmt.Sprintf(format+string(debug.Stack())+_CLEAR_COLOR, v...)
	l.Output(s,RUN_ERROR)

}

func AddLogger(name string, logger *log.Logger) {
	Dglog.AddLogger(name,logger)
}

func Level(level RunLevel) bool {
	return Dglog.Level(level)
}

func SetFlags(flag int) {
	Dglog.SetFlags(flag)
}

func SetPrefix(prefix string) {
	Dglog.SetPrefix(prefix)
}

func Info(v ...interface{}) {
	Dglog.Info(v...)
}

func Infof(format string, v ...interface{}) {
	Dglog.Infof(format, v...)
}

func Debug(v ...interface{}) {
	Dglog.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	Dglog.Debugf(format, v...)
}

func Warning(v ...interface{}) {
	Dglog.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	Dglog.Warningf(format, v...)
}

func Error(v ...interface{}) {
	Dglog.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	Dglog.Errorf(format, v...)
}
