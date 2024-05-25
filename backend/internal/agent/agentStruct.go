// структура агента и запуск в init, чтобы при импорте в main сразу запускался агент
package agent

import (
	"time"
)

var (
	COMPUTING_POWER int   = 3
	ag              agent = agent{
		workCh: make(chan agReq),
	}
)

type agent struct {
	workCh chan agReq
}

type agAns struct {
	Id     agentID `json:"id"`
	Result float64 `json:"result"`
}
type agReq struct {
	Task struct {
		Id        agentID       `json:"id"`
		Arg1      float64       `json:"arg1"`
		Arg2      float64       `json:"arg2"`
		Operation string        `json:"operation"`
		TimeOp    time.Duration `json:"operation_time"`
	} `json:"task"`
}
type agentID struct {
	ReqID ID `json:"reqID"`
	OpID  ID `json:"opID"`
}
type ID uint16

func init() {
	go ag.getTask()
	for range COMPUTING_POWER {
		go work()
	}
}
