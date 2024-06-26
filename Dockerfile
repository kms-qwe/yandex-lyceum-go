# Установим базовый образ Go
FROM golang:1.22

# Установим рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY backend/go.mod backend/go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY backend/ ./

# Сборка приложения
RUN go build -o main ./cmd

# Устанавливаем порт на котором будет работать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
