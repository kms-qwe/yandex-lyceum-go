package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handlerNewExp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data calculateRequest

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusInternalServerError)
		return
	}

	polishExpr, err := infixToPostfix(data.Expression)
	if err != nil {
		http.Error(w, "Невалидные данные", http.StatusUnprocessableEntity)
		return
	}
	if _, ok := orch.db[data.Id]; ok {
		http.Error(w, "Id повторяется", http.StatusUnprocessableEntity)
		return
	}
	orch.mu.Lock()
	orch.db[data.Id] = dbRecord{
		status: "Принято в обработку",
		result: 0.0}
	orch.mu.Unlock()
	fmt.Println(data.Expression)
	fmt.Println(polishExpr)
	fmt.Println("ПРИНЯЛ")
	go eval(polishExpr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Выражение принято для вычисления\n"))
}

func handlerListExpr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}

	orch.mu.Lock()
	resp := listExprStatus{}
	for k, v := range orch.db {
		resp.ListOfReq = append(resp.ListOfReq, exprStatus{k, v.status, v.result})
	}
	orch.mu.Unlock()

	jsonData, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
