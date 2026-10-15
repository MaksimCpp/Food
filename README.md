# 🍔 Food

Бэкенд сервиса доставки еды, построенный на **микросервисной архитектуре** на Go.
Сервисы общаются между собой по **gRPC**, а внешние клиенты работают с системой через единую точку входа — **API Gateway**.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-protobuf-244c5a?logo=grpc&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Kafka](https://img.shields.io/badge/Apache%20Kafka-3.9-231F20?logo=apachekafka&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)

---

## 📑 Содержание

- [Архитектура](#-архитектура)
- [Сервисы](#-сервисы)
- [Сервис авторизации](#-сервис-авторизации-auth)
- [Технологии](#-технологии)
- [Структура проекта](#-структура-проекта)
- [Быстрый старт](#-быстрый-старт)
- [Переменные окружения](#-переменные-окружения)
- [Roadmap](#-roadmap)

---

## 🏗 Архитектура
![Архитектура](./images/архитектура.png)


- **API Gateway** — единая точка входа: принимает HTTP-запросы от клиентов и проксирует их в нужные сервисы по gRPC.
- **Внутреннее взаимодействие** — только по gRPC, контракты описаны в `.proto`-файлах.
- **Database per service** — у каждого сервиса своя база данных в PostgreSQL.
- **Аутентификация** — каждый сервис проверяет access-токен через общий gRPC-интерцептор из `pkg/interceptor`, который обращается к `AuthService.ValidateAccessToken`.

---

## 🧩 Сервисы

| Сервис | Описание | Статус |
|---|---|---|
| 🔐 **Auth** | Регистрация, вход, выдача и обновление JWT-токенов, роли пользователей | ✅ Готов |
| 📦 **Products** | Каталог товаров: создание, редактирование, получение списка | 🚧 В разработке |
| 🧾 **Orders** | Оформление заказов и отслеживание их статуса | 📝 Запланирован |
| 🌐 **API Gateway** | Единая HTTP-точка входа, маршрутизация запросов в сервисы по gRPC | 📝 Запланирован |

---

## 🔐 Сервис авторизации (Auth)

Отвечает за пользователей и токены. Поддерживает три роли: `user`, `admin`, `courier`.

### gRPC API — `auth.AuthService`

| Метод | Описание | Доступ |
|---|---|---|
| `RegisterUser` | Регистрация покупателя | Публичный |
| `RegisterAdmin` | Регистрация администратора (нужен `admin_code`) | Публичный |
| `RegisterCourier` | Регистрация курьера | 🔒 Только `admin` |
| `LoginUser` | Вход покупателя | Публичный |
| `LoginAdmin` | Вход администратора (нужен `admin_code`) | Публичный |
| `LoginCourier` | Вход курьера | Публичный |
| `Refresh` | Обновление пары токенов по refresh-токену | Публичный |
| `ValidateAccessToken` | Проверка access-токена, возвращает `user_id` и `role` | Публичный (для сервисов) |

Полный контракт: [`services/auth/proto/auth.proto`](services/auth/proto/auth.proto)

### Токены

| Тип | Время жизни |
|---|---|
| Access token | 15 минут |
| Refresh token | 7 дней |

Токен передаётся в gRPC-метаданных:

```
authorization: Bearer <access_token>
```

### Устройство сервиса

Сервис построен по принципам **чистой архитектуры**:

```
services/auth/
├── cmd/server/            # точка входа
├── proto/                 # .proto и сгенерированный код
└── internal/
    ├── config/            # загрузка конфигурации из env
    ├── domain/            # сущности, интерфейсы, доменные ошибки
    ├── dto/               # входные/выходные структуры use-case'ов
    ├── service/           # бизнес-логика
    ├── repository/        # работа с PostgreSQL (pgx)
    ├── infrastructure/    # JWT
    └── delivery/grpc/     # gRPC-хендлеры
```

---

## 🛠 Технологии

- **Go 1.25** + Go workspaces (`go.work`)
- **gRPC / Protocol Buffers** — межсервисное взаимодействие
- **PostgreSQL 17** + [pgx](https://github.com/jackc/pgx) — хранение данных
- **golang-migrate** — миграции БД
- **JWT** ([golang-jwt](https://github.com/golang-jwt/jwt)) — аутентификация
- **bcrypt** — хеширование паролей
- **Redis 7** — кеширование
- **Apache Kafka 3.9** (KRaft) + Kafka UI — асинхронные события
- **Docker / Docker Compose** — локальное окружение

---

## 📁 Структура проекта

```
Food/
├── services/
│   └── auth/              # сервис авторизации
├── pkg/
│   ├── authcontext/       # хелперы для user_id / role в context
│   └── interceptor/       # общие gRPC-интерцепторы (Auth, Role)
├── migrations/
│   └── auth/              # SQL-миграции сервиса auth
├── postgres/
│   └── init.sql           # создание баз данных при старте
├── docker-compose.yml
├── .env.example
└── go.work
```

---

## 🚀 Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)
- `protoc` + плагины `protoc-gen-go`, `protoc-gen-go-grpc` (для генерации кода)

### Запуск

```bash
# 1. Клонировать репозиторий
git clone https://github.com/MaksimCpp/Food.git
cd Food

# 2. Создать файл окружения
cp .env.example .env

# 3. Поднять инфраструктуру и сервисы
docker compose up -d --build
```

После запуска будут доступны:

| Компонент | Адрес |
|---|---|
| Auth gRPC | `localhost:50051` |
| PostgreSQL | `localhost:5433` |
| Redis | `localhost:6379` |
| Kafka | `localhost:9092` |
| Kafka UI | http://localhost:8080 |

### Проверка через grpcurl

```bash
grpcurl -plaintext \
  -import-path services/auth/proto -proto auth.proto \
  -d '{"username":"john","email":"john@example.com","password":"secret"}' \
  localhost:50051 auth.AuthService/RegisterUser
```

### Генерация gRPC-кода

```bash
protoc -I services/auth/proto \
  --go_out=services/auth/proto --go_opt=paths=source_relative \
  --go-grpc_out=services/auth/proto --go-grpc_opt=paths=source_relative \
  services/auth/proto/auth.proto
```

---

## ⚙️ Переменные окружения

Полный список — в [`.env.example`](.env.example). Основные:

| Переменная | Описание |
|---|---|
| `POSTGRES_USER`, `POSTGRES_PASSWORD` | Учётные данные PostgreSQL |
| `POSTGRES_HOST`, `POSTGRES_PORT` | Адрес PostgreSQL |
| `AUTH_DB` | Имя базы данных сервиса auth |
| `JWT_SECRET_KEY` | Секрет для подписи JWT |
| `ADMIN_CODE` | Код для регистрации и входа администраторов |
| `GRPC_AUTH_PORT` | Порт gRPC-сервера auth |
| `GRPC_AUTH_ADDRESS` | Адрес auth-сервиса для других сервисов |
| `REDIS_PORT` | Порт Redis |
| `KAFKA_BROKER` | Адрес брокера Kafka |

---

## 🗺 Roadmap

- [x] Сервис авторизации (JWT, роли, refresh-токены)
- [x] Общие gRPC-интерцепторы аутентификации и ролей
- [ ] Сервис товаров
- [ ] Сервис заказов
- [ ] API Gateway
- [ ] Асинхронные события через Kafka
- [ ] Кеширование через Redis
- [ ] Тесты и CI
