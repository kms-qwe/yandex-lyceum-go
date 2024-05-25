package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (a agent) getTask() {
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:8080/internal/task", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		// Проверка статуса ответа
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				time.Sleep(3 * time.Second)
				continue
			}
			aR := agReq{}
			err = json.Unmarshal(body, &aR)
			if err != nil {
				fmt.Println("Ошибка при получении задачи при преобразовании из JSON")
				time.Sleep(3 * time.Second)
				continue
			}
			go func() {
				ag.workCh <- aR
			}()
		} else if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Нет задачи")
		} else {
			fmt.Printf("Ошибка: %s\n", resp.Status)
		}
		resp.Body.Close()
		time.Sleep(time.Millisecond)
	}
}
