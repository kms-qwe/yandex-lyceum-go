package orchestrator

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	orch *orchestrator = &orchestrator{
		db:          map[ID]dbRecord{},
		chToAgent:   make(chan []byte),
		chFromAgent: make(chan []byte),
		mu:          sync.Mutex{}}
	TIME_ADDITION_MS        int
	TIME_SUBTRACTION_MS     int
	TIME_MULTIPLICATIONS_MS int
	TIME_DIVISIONS_MS       int
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
type dbRecord struct {
	status string
	result float64
}
type orchestrator struct {
	db          map[ID]dbRecord
	chToAgent   chan []byte
	chFromAgent chan []byte
	mu          sync.Mutex
}

func init() {
	fmt.Println("Init function in orchestrator package is called.")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handlerNewExp)
	mux.HandleFunc("/api/v1/expressions", handlerListExpr)
	http.ListenAndServe(":8080", mux)
}
