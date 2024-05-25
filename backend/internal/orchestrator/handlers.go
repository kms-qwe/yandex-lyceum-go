package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	go eval(polishExpr, data.Id)

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
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func handlerExprById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}
	sp := strings.Split(r.URL.String(), "/")
	id, err := strconv.Atoi(sp[len(sp)-1])
	if err != nil {
		http.Error(w, "ошибка при str -> int", http.StatusInternalServerError)
		return
	}
	orch.mu.Lock()
	val, ok := orch.db[ID(id)]
	orch.mu.Unlock()
	resp := exprStatus{Id: ID(id), Status: val.status, Result: val.result}
	if !ok {
		http.Error(w, "Нет записи", http.StatusNotFound)
		return
	}

	jsonData, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
func handlerAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// fmt.Println("GET")
		select {
		case <-time.After(time.Second * 20):
			http.Error(w, "Нет задачи", http.StatusNotFound)
			return
		case task := <-orch.chToAgent:
			// fmt.Println(task)
			jsonData, err := json.MarshalIndent(task, "", "    ")
			if err != nil {
				http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)

		}
	} else if r.Method == http.MethodPost {
		// fmt.Println("POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		ans := agAns{}

		err = json.Unmarshal(body, &ans)
		if err != nil {
			http.Error(w, "Ошибка при декодировании JSON", http.StatusInternalServerError)
			return
		}

		orch.muOp.Lock()
		orch.dbOp[ans.Id] = ans.Result
		orch.muOp.Unlock()
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		// fmt.Println("МЕТОД НЕ ")
		return
	}

}
