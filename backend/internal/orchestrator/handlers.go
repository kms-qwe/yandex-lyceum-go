// обработчики для оркестратора
package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// обработчик, принимающий новые выражения для вычисления
func handlerNewExp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		log.Println("Не принято на вычисление: Метод не поддерживается")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
		log.Println("Не принято на вычисление: Ошибка при чтении тела запроса")
		return
	}
	defer r.Body.Close()

	var data calculateRequest

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusInternalServerError)
		log.Println("Не принято на вычисление: Ошибка при декодировании JSON")
		return
	}

	polishExpr, err := infixToPostfix(data.Expression)
	if err != nil {
		http.Error(w, "Невалидные данные", http.StatusUnprocessableEntity)
		log.Println("Не принято на вычисление: Невалидные данные")
		return
	}
	if _, ok := orch.db[data.Id]; ok {
		http.Error(w, "Id повторяется", http.StatusUnprocessableEntity)
		log.Println("Не принято на вычисление: Id повторяется")
		return
	}
	orch.mu.Lock()
	orch.db[data.Id] = dbRecord{
		status: "Принято в обработку",
		result: 0.0}
	orch.mu.Unlock()

	log.Println("Принято в обработку", data.Id, data.Expression, polishExpr)
	go eval(polishExpr, data.Id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Выражение принято для вычисления\n"))
}

// Вывод списка выражения
func handlerListExpr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		log.Println("Вывод списка выражений не произведен: Метод не поддерживается")
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
		log.Println("Вывод списка выражений не произведен: Ошибка при кодировании JSON")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	log.Println("Вывод списка выражений произведен")

}

// Вывод выражения с ID
func handlerExprById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}
	sp := strings.Split(r.URL.String(), "/")
	id, err := strconv.Atoi(sp[len(sp)-1])
	if err != nil {
		http.Error(w, "ошибка при str -> int", http.StatusInternalServerError)
		log.Println("Ошибка при выводе выражения по ID: ошибка преобразования str -> int")
		return
	}
	orch.mu.Lock()
	val, ok := orch.db[ID(id)]
	orch.mu.Unlock()
	resp := exprStatus{Id: ID(id), Status: val.status, Result: val.result}
	if !ok {
		http.Error(w, "Нет записи", http.StatusNotFound)
		log.Println("Ошибка при выводе выражения по ID: нет записи")
		return
	}

	jsonData, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		log.Println("Ошибка при выводе выражения по ID: Ошибка при кодировании JSON")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	log.Println("Выражение выведено с ID", id)
}

// Ручка для агента
func handlerAgent(w http.ResponseWriter, r *http.Request) {
	//отдача задачи
	if r.Method == http.MethodGet {
		select {
		case <-time.After(time.Second * 20):
			http.Error(w, "Нет задачи", http.StatusNotFound)
			log.Println("Отдача агенту: нет задачи")
			return
		case task := <-orch.chToAgent:
			jsonData, err := json.MarshalIndent(task, "", "    ")
			if err != nil {
				http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
				log.Println("Отдача агенту: Ошибка при кодировании JSON")
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			log.Println("Отдача агенту: задача отдана", task)
		}
	} else if r.Method == http.MethodPost {
		//получение ответа от агента
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
			log.Println("Получение от агента: Ошибка при чтении тела запроса")
			return
		}
		defer r.Body.Close()

		ans := agAns{}

		err = json.Unmarshal(body, &ans)
		if err != nil {
			http.Error(w, "Ошибка при декодировании JSON", http.StatusInternalServerError)
			log.Println("Получение от агента: Ошибка при декодировании JSON")
			return
		}

		orch.muOp.Lock()
		orch.dbOp[ans.Id] = ans.Result
		orch.muOp.Unlock()
		w.WriteHeader(http.StatusOK)
		log.Println("Получение от агента", ans)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		log.Println("Ручка для агента: Метод не поддерживается")
		return
	}

}
