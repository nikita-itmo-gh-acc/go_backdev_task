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
### Пример ответа:
```json
{"session_id":"884c6257-c842-4029-bb07-b7d74fca2eaf","user_id":"c84f18a2-c6c7-4850-be15-93f9cbaef3b3","access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJcEFkcmVzcyI6IjEyNy4wLjAuMSIsImlzcyI6ImJhY2t0ZXN0Iiwic3ViIjoiYzg0ZjE4YTItYzZjNy00ODUwLWJlMTUtOTNmOWNiYWVmM2IzIiwiYXVkIjpbInVzZXIiXSwiZXhwIjoxNzMzMzQ3Mjg0LCJpYXQiOjE3MzMzNDU0ODR9.KGDwSkKdtn6Tsv7RVhbj_MhDCpfSML0o3XJl7L9KVN0nIR8IN0fnnYJJY19zszvoZcIL_TgmlvjMsXwv2Z_rWQ","ip_address":"127.0.0.1"}
```

### Обновляем токены:
```shell
curl -v -X POST -H "Content-Type: application/json" -d @refresh.json -b "RefreshTokenCookie=ZGFhZDQyOGYtNDlkOC00M2JhLTg2MDYtZTUwYjEzNjNhYzJm" --cookie-jar cookie.txt http://localhost:4242/api/auth/refresh-tokens
```
### Пример refresh.json

```json
{"session_id":"6397d5ec-876e-47c6-b92f-139050a3df46","ip_address":"127.0.0.1"}
```

### Пример ответа:
```json
{"session_id":"ade8d122-a18b-4aa9-9e2d-c6e5f3dd1078","user_id":"c84f18a2-c6c7-4850-be15-93f9cbaef3b3","access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJcEFkcmVzcyI6IjEyNy4wLjAuMSIsImlzcyI6ImJhY2t0ZXN0Iiwic3ViIjoiYzg0ZjE4YTItYzZjNy00ODUwLWJlMTUtOTNmOWNiYWVmM2IzIiwiYXVkIjpbInVzZXIiXSwiZXhwIjoxNzMzMzQ3NDgzLCJpYXQiOjE3MzMzNDU2ODN9.o6yPPRUFPKY2yzdjIEqRpV_QgK1u--yPsq-fAiQeZP3M1JYiD0wp-nKwUPtLWDLCz3ULWVrJwja3aJDKSqfTmg","ip_address":"127.0.0.1"}
```

### Получаем доступ к защищенным данным:
```shell
curl -v GET -H "Content-Type: application/json" -d @protected.json --cookie-jar cookie.txt http://localhost:4242/protected
```

### Пример protected.json
```json
{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJcEFkcmVzcyI6IjEyNy4wLjAuMSIsImlzcyI6ImJhY2t0ZXN0Iiwic3ViIjoiYzg0ZjE4YTItYzZjNy00ODUwLWJlMTUtOTNmOWNiYWVmM2IzIiwiYXVkIjpbInVzZXIiXSwiZXhwIjoxNzMzMzE4MzQyLCJpYXQiOjE3MzMzMTY1NDJ9.-D3wy2VLnWkyYHwTIGBfRMmmQXE1047cHuuHR4veB35cNzc8lN1vvCfDx0vG0xdL5RiMfznpsxc-oQ1JJVZT3Q","ip_address":"127.0.0.1"}
```

### Пример ответа:
```json
{"message":"You've successfully visited protected page"}
```

### Пример логов:
```
INFO: 2024/12/04 23:51:18 connection.go:44: Established connection with PostreSQL database - Name: auth_db, port: 5432
INFO: 2024/12/04 23:51:18 main.go:44: Запуск сервера на 0.0.0.0:4242
INFO: 2024/12/04 23:51:24 main.go:184: Session initialization process begin...
INFO: 2024/12/04 23:51:24 session_storage.go:39: CREATED session with ID = 884c6257-c842-4029-bb07-b7d74fca2eaf
INFO: 2024/12/04 23:51:24 main.go:198: Session successfully initialized!
INFO: 2024/12/04 23:54:43 main.go:223: Refresh tokens process begin...
INFO: 2024/12/04 23:54:43 session_storage.go:48: DELETED session with ID = 884c6257-c842-4029-bb07-b7d74fca2eaf
INFO: 2024/12/04 23:54:43 session_storage.go:39: CREATED session with ID = ade8d122-a18b-4aa9-9e2d-c6e5f3dd1078
INFO: 2024/12/04 23:54:43 main.go:271: Tokens successfully updated!
ERROR: 2024/12/04 23:57:04 main.go:79: Can't parse access token token has invalid claims: token is expired
ERROR: 2024/12/04 23:57:04 main.go:142: Access token verification failed..
```