package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``

/*
Для запуска тестов с включённой аутентификацией:

1. Запустить сервер:
   go run .

2. Получить токен:
   curl -X POST http://localhost:7540/api/signin \
     -H "Content-Type: application/json" \
     -d '{"password":""}' - тут вставить пароль из .env

3. Скопировать значение token из ответа и временно вставить его в tests/settings.go:
   var Token = `...`

4. Перезапустить сервер и тесты:
   kill $(lsof -ti :7540) 2>/dev/null
   go run . & sleep 1 && go test -count=1 ./...
*/
