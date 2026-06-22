# Subscription Service - REST API для управления подписками пользователей

## СТЕК ТЕХНОЛОГИЙ

- Go 1.21
- Gin (веб-фреймворк)
- PostgreSQL 15
- Swagger (документация API)
- Docker / Podman
- godotenv (конфигурация)


## КРАТКОЕ ОПИСАНИЕ

Сервис предоставляет API для CRUD операций с подписками пользователей. 
Каждая подписка содержит:
- название сервиса
- стоимость
- ID пользователя (UUID)
- дату начала
- дату окончания (опционально)

Также реализован эндпоинт для подсчёта суммарной стоимости подписок за выбранный 
период с фильтрацией по пользователю и названию сервиса.


## API ЭНДПОИНТЫ

POST   /api/v1/subscriptions - создать подписку
GET    /api/v1/subscriptions - список подписок (пагинация)
GET    /api/v1/subscriptions/{id} - получить по ID
PUT    /api/v1/subscriptions/{id} - обновить
DELETE /api/v1/subscriptions/{id} - удалить
GET    /api/v1/subscriptions/total-cost - подсчёт стоимости 


## УСТАНОВКА И ЗАПУСК

Способ 1: Docker Compose (рекомендуется)

``` bash
git clone https://github.com/PhosFactum/effective-mobile-test-go
cd effective-mobile-test-go
docker compose up --build
```

# или для Podman
``` bash
podman compose up --build
```

После запуска:
- API доступен: http://localhost:8080
- Swagger документация: http://localhost:8080/swagger/index.html


Способ 2: Локальный запуск

Требуется PostgreSQL 15+ и Go 1.21+

1. Установить зависимости:
``` bash
go mod tidy
```

2. Сгенерировать Swagger документацию:
``` bash
swag init -g cmd/main.go -o docs
```

3. Создать файл .env:
``` env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
SERVER_PORT=8080
```

4. Запустить приложение:
``` bash
go run cmd/main.go
```


## ТЕСТИРОВАНИЕ API (если не через Swagger)

Примеры запросов с curl:

1. Создать подписку
``` bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

2. Получить все подписки
``` bash
curl http://localhost:8080/api/v1/subscriptions
```

3. Получить с пагинацией
``` bash
curl "http://localhost:8080/api/v1/subscriptions?limit=5&offset=0"
```

4. Получить по ID
``` bash
curl http://localhost:8080/api/v1/subscriptions/ваш-uuid
```

5. Обновить подписку
``` bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/ваш-uuid \
  -H "Content-Type: application/json" \
  -d '{"price": 500}'
```

6. Удалить подписку
``` bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/ваш-uuid
```

7. Подсчитать стоимость за период
``` bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2025&end_date=12-2025"
```

8. С фильтрацией по пользователю
``` bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2025&end_date=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```


9. С фильтрацией по названию сервиса
``` bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2025&end_date=12-2025&service_name=Netflix"
```


## ОЖИДАЕМЫЕ КОДЫ ОТВЕТА

200 OK           - успешный GET запрос
201 Created      - успешное создание
204 No Content   - успешное удаление
400 Bad Request  - ошибка валидации
404 Not Found    - запись не найдена
500 Internal Server Error - ошибка сервера


## ПРИМЕЧАНИЯ

- user_id должен быть в формате UUID v4
  Пример: 60601fee-2bf1-4721-ae6f-7636e79a0cba

- Даты принимаются и возвращаются в формате MM-YYYY
  Пример: 07-2025

- end_date опционален (может быть null)

- Цены указываются в рублях, целые числа (копейки не учитываются)


## ОСТАНОВКА СЕРВИСА

#### Остановить контейнеры
``` bash
docker compose down
```

#### или для Podman
``` bash
podman compose down
```

#### Остановить с удалением томов (очистка БД)
``` bash
docker compose down -v
```
