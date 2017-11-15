package logs

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func InitLog(fileName string, level L) {
	InitLogAdapter(fileName)
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(level))
	log.AddHook(LogAdapterInstance)
}

type F log.Fields
type L log.Level

func commonFileds(fields F) *log.Entry {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["0pos"] = fmt.Sprintf("%s:%d", file, line)
	return log.WithFields(log.Fields(fields))
}

func Log(kvs map[string]interface{}) *log.Entry {
	return commonFileds(kvs)
}

func WithField(key string, val interface{}) *log.Entry {
	return commonFileds(F{key: val})
}
