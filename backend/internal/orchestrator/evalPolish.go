// eval запускается для каждого полученного выражения в отдельной горутине и высчитывает результат
package orchestrator

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func eval(initExpr string, rId ID) {
	expr := initExpr
	for {
		log.Printf("Индекс выражения: %d; Начальное выражение: |%s|; Выражение сейчас: |%s|\n", rId, initExpr, expr)

		//Считаем части, которые можно выполнить параллельно
		elementsOfExpr := strings.Fields(expr)
		numOp := 0
		for i := range len(elementsOfExpr) - 2 {
			if isNumber(elementsOfExpr[i]) && isNumber(elementsOfExpr[i+1]) && isOperator([]rune(elementsOfExpr[i+2])[0]) {
				numOp += 1
			}
		}

		// если частей, которые нужно посчитать нет, то выражение готово
		if numOp == 0 {
			orch.mu.Lock()
			r, err := strconv.ParseFloat(elementsOfExpr[0], 64)
			dbR := dbRecord{}
			if err != nil {
				dbR.status = "Ошибка при переводе в число в evalPolish"
				log.Println("Ошика при переводе в число в evalPolish с индексом", rId)
			} else {
				dbR.status = "Ответ готов"
				log.Println("Ответ готов для индекса", rId)
			}
			dbR.result = r
			orch.db[rId] = dbR
			orch.mu.Unlock()
			break
		}

		// Все части отправляем в канал, который отдает выражения агенту, делается в отдельной горутине, чтобы не блокировать получение
		go func() {
			cnt := 0
			for i := range len(elementsOfExpr) - 2 {
				if isNumber(elementsOfExpr[i]) && isNumber(elementsOfExpr[i+1]) && isOperator([]rune(elementsOfExpr[i+2])[0]) {
					cnt += 1
					aR := agReq{}
					switch []rune(elementsOfExpr[i+2])[0] {
					case []rune("+")[0]:
						aR.Task.TimeOp = TIME_ADDITION_MS
					case []rune("-")[0]:
						aR.Task.TimeOp = TIME_ADDITION_MS
					case []rune("*")[0]:
						aR.Task.TimeOp = TIME_ADDITION_MS
					case []rune("/")[0]:
						aR.Task.TimeOp = TIME_ADDITION_MS
					}
					aR.Task.Arg1, _ = strconv.ParseFloat(elementsOfExpr[i], 64)
					aR.Task.Arg2, _ = strconv.ParseFloat(elementsOfExpr[i+1], 64)
					aR.Task.Operation = elementsOfExpr[i+2]
					aR.Task.Id = agentID{rId, ID(cnt)}
					log.Println("Отправлено на вычисление", aR)
					orch.chToAgent <- aR
				}

			}
		}()

		//Получаем результаты из dbOp, полученные результаты удаляем из таблицы, чтобы в следующей итерации evalPolish не дублировались индексы
		resEval := map[agentID]float64{}
		for len(resEval) < numOp {
			orch.muOp.Lock()
			for i := range numOp {
				aID := agentID{rId, ID(i + 1)}
				if _, ok := resEval[aID]; !ok {
					if val, ok := orch.dbOp[aID]; ok {
						resEval[aID] = val
						delete(orch.dbOp, aID)
					}
				}
			}
			orch.muOp.Unlock()
			time.Sleep(1 * time.Second)
		}
		//Заменяем части на их значения и изменяем текущее выражение
		cnt := 0
		for i := range len(elementsOfExpr) - 2 {
			if isNumber(elementsOfExpr[i]) && isNumber(elementsOfExpr[i+1]) && isOperator([]rune(elementsOfExpr[i+2])[0]) {
				cnt += 1
				r := resEval[agentID{rId, ID(cnt)}]
				strRes := strconv.FormatFloat(r, 'f', 6, 64)
				elementsOfExpr[i], elementsOfExpr[i+1], elementsOfExpr[i+2] = strRes, "", ""
			}
		}
		expr = strings.Join(elementsOfExpr, " ")

	}

}
func isNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
