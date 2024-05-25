// getTask отправяет запросы дял получения задачи каждые reqTime секунд, если получена задача, она отправляется в канал воркерам
package agent

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	reqTime time.Duration = 500 * time.Millisecond
)

func (a agent) getTask() {
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:8080/internal/task", nil)
		if err != nil {
			log.Println("Error creating request:", err)
			time.Sleep(reqTime)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			time.Sleep(reqTime)
			continue
		}

		// Проверка статуса ответа
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				time.Sleep(reqTime)
				continue
			}
			aR := agReq{}
			err = json.Unmarshal(body, &aR)
			if err != nil {
				log.Println("Ошибка при получении задачи при преобразовании из JSON")
				time.Sleep(reqTime)
				continue
			}
			go func() {
				ag.workCh <- aR
			}()
		} else if resp.StatusCode == http.StatusNotFound {
			log.Println("Нет задачи")
		} else {
			log.Printf("Ошибка: %s\n", resp.Status)
		}
		resp.Body.Close()
		time.Sleep(reqTime)
	}
}
