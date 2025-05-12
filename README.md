# 🏦 Банковский API на Go

REST API Backend-сервис для управления банковскими операциями с JWT-аутентификацией, транзакциями, кредитами и интеграцией с внешними сервисами.

---

## 📋 Основные функции

### Пользователи
- Регистрация новых пользователей с уникальными email и username
- Аутентификация с выдачей JWT-токена сроком на 24 часа

### Управление счетами
  - Создание счетов
  - Переводы средств между счетами
  - Пополнение и списание средств со счета

### Карты
- Выпуск виртуальных карт с безопасным хранением данных:
  - Номер карты генерируется по алгоритму Луна
  - Номер хранится в зашифрованном виде (PGP)
  - CVV хранится в виде bcrypt-хеша
- Просмотр данных карты владельцем

### Кредитные операции
- Оформление кредитов с аннуитетными платежами
- Предоставление графика платежей
- Автоматическое списание платежей через планировщик
- Штрафы за просрочку платежа (+10% к сумме)

### Безопасность
- **Пароли**: хеширование с использованием bcrypt (cost 12+)
- **JWT**: подпись HMAC-SHA256, срок действия 24 часа
- **Данные карт**:
  - Номер карты: PGP-симметричное шифрование
  - CVV: bcrypt-хеш
  - Целостность: HMAC-SHA256
- **Авторизация**: проверка владения ресурсами по userID

### Аналитика
- Анализ доходов и расходов за месяц
- Оценка кредитной нагрузки
- Прогнозирование баланса на срок до 365 дней

- ### Интеграции:
- **ЦБ РФ**: интеграция с SOAP API для получения ключевой ставки
  - URL: `https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx`
  - Парсинг XML-ответов с использованием etree
- **SMTP**: Email-уведомления о регистрации, операциях и просрочках платежей

---

## ⚙️ Технологии

- **Язык программирования**: Go 1.23+
- **Маршрутизация**: `github.com/gorilla/mux`
- **База данных**: PostgreSQL 17
- **Драйвер БД**: `github.com/lib/pq`
- **Аутентификация**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Логирование**: `github.com/sirupsen/logrus`
- **Шифрование/хеширование**: bcrypt, HMAC-SHA256, PGP
- **Email**: `gopkg.in/gomail.v2`
- **Обработка XML**: `github.com/beevik/etree`
- **Планировщик задач**: стандартный `time.Ticker`

---

## 🛠 Установка и запуск

### 1. Подготовка

- Убедитесь, что установлены:
    - Go 1.23+
    - PostgreSQL 17+
    
### 2. Настройка

Клонируйте репозиторий:

```bash
git clone https://github.com/vanhellthing93/sf.mephi.go_homework
```

Сгенерируйте файл pgp ключа и сохраните его в корень проекта ./pgp.key. 
Для генерации можно воспользоваться сервисом https://onlinepgp.com/

Скопируйте и настройте .env файл:
``` env
## Данные для подключения к бд
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bank_service
SSL_MODE=disable

## Данные для подключения к smtp
SMTP_HOST=smtp.example.com
SMTP_PORT=465
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=email_password
SMTP_FROM=user@example.com

##Константа увеличения ставки цб для кредитов
CB_RATE_INCREMENT=2.5

## Переменные безопасности
HMAC_SECRET=HMAC_SECRET
JWT_SECRET=your_strong_jwt_secret

## относительный путь до файла ключа pgp. Если использовали другое название или путь, то здесь необходимо скорректировать
PGP_PRIVATE_KEY=./pgp.key 

```



### 3. Сборка и запуск

Загрузите зависимости и запустите приложение:

```bash
go mod download
go run main.go
```

---

## 📂 Структура проекта

```
go_homework/
├── configs/       # Конфигурации
├── internal/
│   ├── handlers/  # HTTP-обработчики
│   ├── models/    # Сущности БД
│   ├── repos/     # Репозитории (работа с БД)
│   ├── services/  # Бизнес-логика
│   └── utils/     # Вспомогательные утилиты и логгер
├── .env.example   # Шаблон конфига
├── .gitignore     # Игнорируемые файлы
├── go.mod         # Модуль Go
├── go.sum         # Контрольные суммы зависимости
└── README.md      # Описание проекта
```

---

## 🔑 Аутентификация
Токен передается в заголовке:
```
Authorization: Bearer <ваш_jwt_токен>
```
---

## 📖 Структура API

| Метод  | Путь                                  | Описание                         | Доступ    |
|--------|---------------------------------------|----------------------------------|-----------|
| POST   | /register                             | Регистрация пользователя         | Публичный |
| POST   | /login                                | Аутентификация (получение JWT)   | Публичный |
| POST   | /accounts                             | Создание банковского счета       | JWT       |
| GET    | /accounts                             | Получение списка счетов          | JWT       |
| POST   | /accounts/{account_id}/cards          | Создание карты для счета         | JWT       |
| GET    | /accounts/{account_id}/cards          | Получение карт счета             | JWT       |
| GET    | /cards/{card_id}                      | Получение информации о карте     | JWT       |
| DELETE | /cards/{card_id}                      | Удаление карты                   | JWT       |
| POST   | /accounts/{from_account_id}/transfers | Создание перевода между счетами  | JWT       |
| GET    | /accounts/{account_id}/transfers      | Получение переводов счета        | JWT       |
| GET    | /transfers/{transfer_id}              | Получение информации о переводе  | JWT       |
| POST   | /credits                              | Создание кредита                 | JWT       |
| GET    | /credits                              | Получение списка кредитов        | JWT       |
| GET    | /credits/{credit_id}/schedule         | Получение графика платежей       | JWT       |
| POST   | /credits/{credit_id}/payments         | Создание платежа по кредиту      | JWT       |
| GET    | /credits/{credit_id}/payments         | Получение платежей по кредиту    | JWT       |
| GET    | /payments/{payment_id}                | Получение информации о платеже   | JWT       |
| POST   | /accounts/{account_id}/transactions   | Создание операции по счету       | JWT       |
| GET    | /accounts/{account_id}/transactions   | Получение операций счета         | JWT       |
| GET    | /transactions/{transaction_id}        | Получение информации об операции | JWT       |
| PATCH  | /transactions/{transaction_id}        | Обновление операции              | JWT       |
| DELETE | /transactions/{transaction_id}        | Удаление операции                | JWT       |
| GET    | /analytics/income-expense             | Статистика доходов/расходов      | JWT       |
| GET    | /analytics/balance-forecast           | Прогноз баланса                  | JWT       |
| GET    | /analytics/credit-load                | Кредитная нагрузка               | JWT       |
| GET    | /analytics/monthly-stats              | Ежемесячная статистика           | JWT       |

## 📖 Примеры API-запросов

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1", "email":"user@example.com", "password":"qwerty123"}'
```

### Авторизация (получение токена)

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com", "password":"qwerty123"}'
```

### Создание счета  (требует авторизации)

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"currency":"RUB"}'
```

### Получение списка счетов (требует авторизации)
```bash
curl -X GET http://localhost:8080/accounts \
  -H "Authorization: Bearer <токен>"
```

### Создание операции по счету (требует авторизации)
```bash
curl -X POST http://localhost:8080/accounts/<account_id>/transactions \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"amount":100.50, "type":"income", "category":"Salary", "description":"Monthly salary"}'
```

### Получение списка операций счета (требует авторизации)
```bash
curl -X GET http://localhost:8080/accounts/<account_id>/transactions \
  -H "Authorization: Bearer <токен>"
```

### Получение информации об операции (требует авторизации)
```bash
curl -X GET http://localhost:8080/transactions/<transaction_id> \
  -H "Authorization: Bearer <токен>"
```

### Обновление операции (требует авторизации)
```bash
curl -X PATCH http://localhost:8080/transactions/<transaction_id> \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"amount":150.75, "type":"income", "category":"Bonus", "description":"Yearly bonus"}'
```

### Удаление операции (требует авторизации)
```bash
curl -X DELETE http://localhost:8080/transactions/<transaction_id> \
  -H "Authorization: Bearer <токен>"
```

### Создание карты для счета (требует авторизации)
```bash
curl -X POST http://localhost:8080/accounts/<account_id>/cards \
  -H "Authorization: Bearer <токен>"
```

### Получение списка карт счета (требует авторизации)
```bash
curl -X GET http://localhost:8080/accounts/<account_id>/cards \
  -H "Authorization: Bearer <токен>"
```

### Получение информации о карте (требует авторизации)
```bash
curl -X GET http://localhost:8080/cards/<card_id> \
  -H "Authorization: Bearer <токен>"
```

### Удаление карты (требует авторизации)
```bash
curl -X DELETE http://localhost:8080/cards/<card_id> \
  -H "Authorization: Bearer <токен>"
```

### Перевод средств  (требует авторизации)

```bash
curl -X POST http://localhost:8080/accounts/<from_account_id>/transfers \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"to_account":<account_id>, "amount":100.50, "description":"Payment for services"}'
```

### Получение списка переводов счета (требует авторизации)
```bash
curl -X GET http://localhost:8080/accounts/<account_id>/transfers \
  -H "Authorization: Bearer <токен>"
```

### Получение информации о переводе (требует авторизации)
```bash
curl -X GET http://localhost:8080/transfers/<transfer_id> \
  -H "Authorization: Bearer <токен>"
```

### Создание кредита (требует авторизации)
```bash
curl -X POST http://localhost:8080/credits \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"amount":10000, "term":12}'
```

### Получение списка кредитов (требует авторизации)
```bash
curl -X GET http://localhost:8080/credits \
  -H "Authorization: Bearer <токен>"
```

### Получение графика платежей по кредиту (требует авторизации)
```bash
curl -X GET http://localhost:8080/credits/<credit_id>/schedule \
  -H "Authorization: Bearer <токен>"
```

### Создание платежа по кредиту (требует авторизации)
```bash
curl -X POST http://localhost:8080/credits/<credit_id>/payments \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"amount":1000}'
```

### Получение списка платежей по кредиту (требует авторизации)
```bash
curl -X GET http://localhost:8080/credits/<credit_id>/payments \
  -H "Authorization: Bearer <токен>"
```

### Получение информации о платеже (требует авторизации)
```bash
curl -X GET http://localhost:8080/payments/<payment_id> \
  -H "Authorization: Bearer <токен>"
```

### Получение статистики доходов и расходов (требует авторизации)
```bash
curl -X GET "http://localhost:8080/analytics/income-expense?start_date=2025-01-01&end_date=2025-01-31" \
  -H "Authorization: Bearer <токен>"
```

### Получение прогноза баланса (требует авторизации)
```bash
curl -X GET "http://localhost:8080/analytics/balance-forecast?days=30" \
  -H "Authorization: Bearer <токен>"
```

### Получение кредитной нагрузки (требует авторизации)
```bash
curl -X GET http://localhost:8080/analytics/credit-load \
  -H "Authorization: Bearer <токен>"
```

### Получение ежемесячной статистики (требует авторизации)
```bash
curl -X GET "http://localhost:8080/analytics/monthly-stats?year=2025" \
  -H "Authorization: Bearer <токен>"
```

---

## 🧪 Тестирование

Используйте **Postman** или **curl** для проверки API. Убедитесь, что работают:

- регистрация и аутентификация
- создание и получение инфомрации о банковских счетах
- создание и получение инфомрации о банковских картах
- банковские переводы
- получение кредитов и информации о них
- создание платежей по кредиту
- аналитика

Для удобства добавлена Postman коллекция запросов
```
./test/Go Homework.postman_collection.json
```
---

## 🖋 Автор

**Косовский Иван**\
Проект реализован в рамках учебного задания МИФИ\
GitHub: [github.com/vanhellthing93](https://github.com/vanhellthing93)