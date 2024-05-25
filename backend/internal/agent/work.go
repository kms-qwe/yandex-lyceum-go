package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func work() {
	for {
		aR := <-ag.workCh
		ans := agAns{
			Id: aR.Task.Id,
		}
		switch aR.Task.Operation {
		case "+":
			ans.Result = aR.Task.Arg1 + aR.Task.Arg2
		case "-":
			ans.Result = aR.Task.Arg1 - aR.Task.Arg2
		case "*":
			ans.Result = aR.Task.Arg1 * aR.Task.Arg2
		case "/":
			ans.Result = aR.Task.Arg1 / aR.Task.Arg2
		}
		time.Sleep(aR.Task.TimeOp)
		jsonData, err := json.MarshalIndent(ans, "", "    ")
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		client := &http.Client{}

		req, err := http.NewRequest("POST", "http://localhost:8080/internal/task", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		statusCode := resp.StatusCode
		fmt.Println("Response status code:", statusCode)
		resp.Body.Close()

	}

}
