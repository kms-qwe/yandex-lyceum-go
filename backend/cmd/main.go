package main

import (
	"fmt"

	_ "github.com/kms-qwe/yandex-lyceum-go/internal/agent"
	_ "github.com/kms-qwe/yandex-lyceum-go/internal/orchestrator"
)

func main() {
	fmt.Println("Сервер готов")
}
