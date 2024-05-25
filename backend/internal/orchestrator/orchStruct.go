package orchestrator

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	orch *orchestrator = &orchestrator{
		db:          map[ID]dbRecord{},
		chToAgent:   make(chan agReq),
		chFromAgent: make(chan agAns),
		mu:          sync.Mutex{},
		dbOp:        map[agentID]float64{},
		muOp:        sync.RWMutex{}}
	TIME_ADDITION_MS        time.Duration = 10 * time.Second
	TIME_SUBTRACTION_MS     time.Duration = 10 * time.Second
	TIME_MULTIPLICATIONS_MS time.Duration = 10 * time.Second
	TIME_DIVISIONS_MS       time.Duration = 10 * time.Second
)

type calculateRequest struct {
	Id         ID     `json:"id"`
	Expression string `json:"expression"`
}
type exprStatus struct {
	Id     ID      `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
type listExprStatus struct {
	ListOfReq []exprStatus `json:"expressions"`
}
type ID uint16
type agentID struct {
	ReqID ID `json:"reqID"`
	OpID  ID `json:"opID"`
}
type dbRecord struct {
	status string
	result float64
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
type agAns struct {
	Id     agentID `json:"id"`
	Result float64 `json:"result"`
}
type orchestrator struct {
	db          map[ID]dbRecord
	chToAgent   chan agReq
	chFromAgent chan agAns
	mu          sync.Mutex
	dbOp        map[agentID]float64
	muOp        sync.RWMutex
}

func init() {
	log.Println("Оркестратор начал работу")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handlerNewExp)
	mux.HandleFunc("/api/v1/expressions", handlerListExpr)
	mux.HandleFunc("/api/v1/expressions/{id}", handlerExprById)
	mux.HandleFunc("/internal/task", handlerAgent)
	http.ListenAndServe(":8080", mux)
}
