# go_backdev_task
## Инструкция (при наличии доступа к БД):
### Запускаем сервер внутри докера 
```shell
docker-compose up -d --build
```

### Генерируем токены:
```shell
curl -v -X POST -H "Content-Type: application/json" -d @generate.json --cookie-jar cookie.txt http://localhost:4242/api/auth/generate-tokens
```
### Пример generate.json
```json
{"user_id":"c84f18a2-c6c7-4850-be15-93f9cbaef3b3","ip_address":"127.0.0.1"}
```
### Обновляем токены:
```shell
curl -v -X POST -H "Content-Type: application/json" -d @refresh.json -b "RefreshTokenCookie=ZGFhZDQyOGYtNDlkOC00M2JhLTg2MDYtZTUwYjEzNjNhYzJm" --cookie-jar cookie.txt http://localhost:4242/api/auth/refresh-tokens
```
### Пример refresh.json

```json
{"session_id":"6397d5ec-876e-47c6-b92f-139050a3df46","ip_address":"127.0.0.1"}
```

### Получаем доступ к защищенным данным:
```shell
curl -v GET -H "Content-Type: application/json" -d @protected.json --cookie-jar cookie.txt http://localhost:4242/protected
```

### Пример protected.json
```json
{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJcEFkcmVzcyI6IjEyNy4wLjAuMSIsImlzcyI6ImJhY2t0ZXN0Iiwic3ViIjoiYzg0ZjE4YTItYzZjNy00ODUwLWJlMTUtOTNmOWNiYWVmM2IzIiwiYXVkIjpbInVzZXIiXSwiZXhwIjoxNzMzMzE4MzQyLCJpYXQiOjE3MzMzMTY1NDJ9.-D3wy2VLnWkyYHwTIGBfRMmmQXE1047cHuuHR4veB35cNzc8lN1vvCfDx0vG0xdL5RiMfznpsxc-oQ1JJVZT3Q","ip_address":"127.0.0.1"}
```