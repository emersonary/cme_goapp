package log

import "github.com/emersonary/go-authentication/log/model"

type tLogHandler struct {
	logList []*model.TLog
}

var LogHandler tLogHandler

func NewLogHandler() *tLogHandler {
	return &tLogHandler{logList: make([]*model.TLog, 16)}
}

func (l *tLogHandler) AddLog(logText string, parentLog *model.TLog) *model.TLog {

	log := model.NewLog(logText, parentLog)

	l.logList = append(l.logList, log)

	return log

}

func init() {

	LogHandler = *NewLogHandler()

}
