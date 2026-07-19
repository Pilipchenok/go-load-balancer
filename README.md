# Go Load Balancer

HTTP load balancer на Go, который распределяет входящие запросы между несколькими backend-серверами, поддерживает разные стратегии балансировки и автоматически исключает недоступные backend-ы через health checks.

---

## Возможности

- HTTP reverse proxy для backend-серверов
- Поддержка нескольких стратегий балансировки:
  - Round Robin (Перебирает по очереди все бэкенды)
  - Least Connections (Выбирает бэкенд с наименьшим количеством подключений)
- Периодические health checks backend-ов
- Автоматическое исключение недоступных backend-ов
- YAML-конфигурация
- Graceful shutdown
- Unit и интеграционные тесты
- Проверка на data races через `go test -race`

---

## Архитектура

- backend (Модель backend-сервера и потокобезопасная работа с его состоянием)
- config (Загрузка и валидация YAML-конфига)
- balancer (Производит подключение к одному из бэкендов, проксирует запрос и возвращает ответ)
- middleware (Логирует каждый запрос)
- healthcheck (Периодически проверяет доступность бэкендов и при необходимости отключает их)
- strategy (Реализация алгоритмов выбора бэкенда)
- main (Собирает все части проекта воедино и запускает сервер)

---

## Как работает запрос

```text
Client -> Load Balancer -> Strategy chooses backend -> Request proxied to backend -> Response returned to client
```

---

## Структура проекта

```text
go-load-balancer/
├── cmd/
│   └── balancer/
│       └── main.go
├── configs/
│   └── config.yaml
├── internal/
│   ├── backend/
│   ├── balancer/
│   ├── config/
│   ├── middleware/
│   └── strategy/
└── README.md
```

---

## Конфигурация

Пример `configs/config.yaml`

```yaml
port: 8080
strategy: "round-robin"
health_check_interval: 5s

backends:
  - url: "http://localhost:8081"
  - url: "http://localhost:8082"
  - url: "http://localhost:8083"
```

Параметры:
- `port` — порт запуска балансировщика
- `strategy` — стратегия выбора бэкенда
  - `round-robin`
  - `least-conn`
- `health_check_interval` — интервал проверки доступности бэкендов
- `backends` — бэкенды, доступные балансировщику
- `url` — URL конкретного бэкенда

---

## Работа с балансировщиком

1. Клонировать репозиторий

```bash
git clone https://github.com/Pilipchenok/go-load-balancer.git
cd go-load-balancer
```

2. Запустить бэкенд-серверы

Терминал 1:
```bash
go run cmd/testserver/main.go 8081
```

Терминал 2:
```bash
go run cmd/testserver/main.go 8082
```

И т.д.

3. Запустить балансировщик

```bash
go run ./cmd/balancer
```

4. Отправлять запросы

```bash
curl http://localhost:8080/
```

(Если в конфиге был указан порт 8080)

5. Пример логов:

```text
Balancer started on :8080
2025-01-15 12:00:01 | GET / | 200 | 1.8ms
2025-01-15 12:00:02 | GET / | 200 | 2.1ms
backend0: true -> false
```

---

## Тестирование

```bash
go test ./... -v
go test -race ./... -v
go vet ./...
```

---

## Что было реализовано

- Потокобезопасная модель backend-а
- Стратегия Round Robin с атомарным счётчиком
- Стратегия Least Connections
- Reverse proxy на `net/http`
- Health checker с `context.Context` и `time.Ticker`
- Graceful shutdown HTTP-сервера
- YAML-конфигурация
- Unit/integration тесты

---

## Технологии

- Go
- `net/http`
- `context`
- `sync/atomic`
- `yaml.v3`
- `httptest`

---

## Репозиторий

GitHub: https://github.com/Pilipchenok/go-load-balancer