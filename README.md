# yandex-lyceum-go

## Что это? 

Распределенный вычислитель арифметических выражений

### Описание

Пользователь хочет считать арифметические выражения. Он вводит строку `2 + 2 * 2` и хочет получить в ответ `6`. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того, вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

### Как работает?

Запускается HTTP сервер (оркестратор), который может обрабатывать выражения. Оркестратор имеет следующие эндпоинты для пользователя:

- Добавление вычисления арифметического выражения
 
```commandline
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "id": <уникалиный идентификатор выраженя>,
      "expression": <строка с выражение>
}'
```
 
- Получение списка выражений
 
```commandline
curl --location 'localhost:8080/api/v1/expressions'

```

- Получение выражения по его идентификатору
 
```commandline
curl --location 'localhost:8080/api/v1/expressions/<интересующий id>'
 
```

Эндпоинты для агента (демона, выполняющего вычисления). Демон и оркестратор сообщаются при помоощи http, что в будущем можно использовать, для горизонтального масштабирования. Агент создает несколько горутин, каждая из которых исполняет вычисления и отправляет ответ на сервер.

- Получение задачи для выполения.
 
```commandline
curl --location 'localhost/internal/task'
 
```
 
- Прием результата обработки данных.
 
```commandline
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
      "id": 1,
      "result": 2.5
}'
 
```

### Параллелизм вычислений

Выражения считается быстрее, если расставить скобки, например 2 + 2 + 2 + 2 будет выполняться последовательно, т.к. операции считаются равноправными. Но (2 + 2) + (2 + 2) будет считаться в два раза быстрее, т.к. выражение распадается на два независимых подвыражения. 

## Деплой

### 1 Клонирования репозитория

```commandline
git clone https://github.com/kms-qwe/yandex-lyceum-go
```
### 2 Сборка докер образа

```commandline
docker compose build
```

### 3 Запуск контейнера

Если хотите видеть логи (в этом случае запросы серверу нужно посылать через новое окно терминала)

```commandline
docker compose up 
```
Если не хотите видеть логи 

```commandline
docker compose up -d
```
### 4 Остановка и удаление контейнера

```commandline
docker compose down
```

## Тесты

- Valid cases
    1. 4 + (-2) + 5 * 6
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 1,
      "expression": "4 + (-2) + 5 * 6"
    }'
    ```
    2. 2 + 2 + 2 + 2
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 2,
      "expression": "(2 + 2) + (2 + 2) + (2 + 2)"
    }'
    ```
    3. 2 + 2 * 4 + 3 - 4 + 5
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 3,
      "expression": "2 + 2 * 4 + 3 - 4 + 5"
    }'
    ```
    4. (23 + 125) - 567 * 23
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 4,
      "expression": "(23 + 125) - 567 * 23"
    }'
    ```
    5. -3 +6
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 5,
      "expression": "-3 +6"
    }'
    ```
- Invalid cases
    1. 4 / 0
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 6,
      "expression": "4 / 0"
    }'
    ```
    2. 45 + x - 5
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 7,
      "expression": "45 + x - 5"
    }'
    ```
    3. 45 + 4*
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 8,
      "expression": "45 + 4*"
    }'
    ```
    4. ---4 + 5
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 9,
      "expression": "---4 + 5"
    }'
    ```
    5. 52 * 3 /
    ```commandline
    curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
      "id": 10,
      "expression": "52 * 3 /"
    }'
    ```

```commandline
curl --location 'localhost:8080/api/v1/expressions'
```















