# Практическое занятие №1
## Разделение монолита на 2 микросервиса. Взаимодействие через HTTP

**ФИО:** Бакланова Е.С.
**Группа:** ЭФМО-01-25

## Цели работы

- Научиться декомпозировать монолитное приложение на микросервисы
- Освоить организацию синхронного HTTP-взаимодействия между сервисами
- Научиться реализовывать межсервисные вызовы с таймаутами и обработкой ошибок
- Освоить прокидывание correlation ID (request-id) для трассировки запросов
- Закрепить практику использования env-конфигурации и базового логирования

## Теория

### Микросервисная архитектура

**Микросервисы** — это архитектурный стиль, при котором приложение строится как набор небольших, независимо развертываемых сервисов. Каждый сервис:
- Имеет свою зону ответственности
- Может быть написан на разных языках
- Взаимодействует с другими по сети (обычно HTTP/REST или message broker)

### Синхронное взаимодействие через HTTP

Сервисы общаются друг с другом прямыми HTTP-запросами. Важные аспекты:

1. **Таймауты** — чтобы не ждать ответ вечно
2. **Retry** — повтор при временных ошибках (но не в этой работе)
3. **Circuit breaker** — защита от падающего сервиса
4. **Request-ID** — сквозной идентификатор для отслеживания цепочки вызовов

### Содержание проекта

**Auth service** (порт 8081)
- Аутентификация пользователей
- Выдача токенов
- Проверка валидности токенов

**Tasks service** (порт 8082)
- CRUD для задач (TODO-список)
- Проверка доступа через Auth service перед каждой операцией

### Эндпоинты

Auth

 | Метод | Эндпоинт | Описание | Тело запроса | Ответ |
 |-------|----------|----------|--------------|--------|
 | POST | /v1/auth/login | Получение токена | {"username":"student","password":"student"} | 200: {"access_token":"demo-token","token_type":"Bearer"} 400: Неверный JSON 401: Неверные данные |
 | GET | /v1/auth/verify | Проверка токена | Заголовок: Authorization: Bearer <token> | 200: {"valid":true,"subject":"student"} 401: {"valid":false,"error":"unauthorized"} |

Tasks

 | Метод | Эндпоинт | Описание | Тело запроса | Ответ |
 |-------|----------|----------|--------------|--------|
 | POST | /v1/tasks | Создать задачу | {"title":"...","description":"...","due_date":"..."} | 201: Задача создана 400: Неверные данные 401: Неавторизован |
 | GET | /v1/tasks | Все задачи | - | 200: Список задач 401: Неавторизован |
 | GET | /v1/tasks/{id} | Задача по ID | - | 200: Задача 404: Не найдена 401: Неавторизован |
 | PATCH | /v1/tasks/{id} | Обновить задачу | {"title":"...","done":true} | 200: Обновлено 404: Не найдена 401: Неавторизован |
 | DELETE | /v1/tasks/{id} | Удалить задачу | - | 204: Удалено 404: Не найдена 401: Неавторизован |

### Структура

<img width="366" height="853" alt="image" src="https://github.com/user-attachments/assets/8aac112d-3a20-4538-b79d-a1c1dd545c20" />

### Инструкция по запуску

1. Запуск Auth

- cd services/auth
- export AUTH_PORT=8081
- go run ./cmd/auth
  
2. Запуск Tasks

- cd services/tasks
- export TASKS_PORT=8082
- export AUTH_BASE_URL=http://localhost:8081
- go run ./cmd/tasks

### Тестирование

1. Получить токен

```bash
curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}'
```

<img width="974" height="510" alt="image" src="https://github.com/user-attachments/assets/9f92481e-8f80-4fcf-88c4-12a97e10942f" />


2. Проверка токена напрямую

```bash
  curl -i http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-002"
```

<img width="974" height="513" alt="image" src="https://github.com/user-attachments/assets/2f2e1007-1cf2-4cdf-a011-1da319babf91" />


3. Создать задачу через Tasks (с проверкой Auth)

```bash
  curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-003" \
  -d '{"title":"Do PZ17","description":"split services","due_date":"2026-01-10"}'
```

<img width="974" height="604" alt="image" src="https://github.com/user-attachments/assets/797cb68b-57fd-4909-9f7a-0d54464f6e31" />

4. Попробовать без токена

```bash
  curl -i http://localhost:8082/v1/tasks \
  -H "X-Request-ID: req-004"
```

<img width="974" height="481" alt="image" src="https://github.com/user-attachments/assets/8b223143-ca77-4434-b673-4d423f7543c6" />

5. Получить список всех задач

```bash
  curl -X GET http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer demo-token"
```

<img width="821" height="541" alt="image" src="https://github.com/user-attachments/assets/2f935a6e-5ebe-4c11-b12e-9604eb6923cc" />

6. Задача по ID

```bash
curl -X GET http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token"
```

<img width="821" height="493" alt="image" src="https://github.com/user-attachments/assets/6b03c026-ea58-4b92-b9ac-6dca99e416f7" />

7. Редактирование задачи

```bash
  curl -X PATCH http://localhost:8082/v1/tasks/t_001 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{"title": "Обновленное название", "done": true}'
```

<img width="825" height="494" alt="image" src="https://github.com/user-attachments/assets/66a3476d-46bf-49a6-9ac9-3cb72e17ed6a" />

8. Удаление

```bash
  curl -X DELETE http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token"
```

<img width="820" height="378" alt="image" src="https://github.com/user-attachments/assets/fd26dacd-4373-47a0-a21b-4816087a423b" />

### Контрольные вопросы

1.	Почему межсервисный вызов должен иметь таймаут?

  	Таймаут нужен, чтобы сервис не завис, если другой сервис упал или тормозит
  	
3.	Чем request-id помогает при диагностике ошибок?

  	Request-id связывает логи разных сервисов по одному запросу. Можно пройти по цепочке и понять, где именно случилась ошибка
  	
5.	Какие статусы нужно вернуть клиенту при невалидном токене?

  	401
  	
7.	Чем опасно "делить одну БД" между сервисами?

   Сервисы становятся связанными через схему БД. Нельзя поменять один сервис, не сломав другой, те еряется независимость 
