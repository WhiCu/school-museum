# Фронтенд музея лицея №86

Простой фронтенд для сайта школьного музея.

## Структура

```
frontend/
├── index.html          # Главная страница
├── exhibition.html     # Страница экспозиции
├── admin.html          # Админ-панель
├── nginx.conf          # Конфигурация nginx для Docker
├── css/
│   ├── styles.css      # Стили сайта
│   └── admin.css       # Стили админки
└── js/
    ├── api.js          # API клиент
    ├── main.js         # Скрипт главной страницы
    ├── exhibition.js   # Скрипт страницы экспозиции
    └── admin.js        # Скрипт админ-панели
```

## Запуск

### Вариант 1: Docker (рекомендуется)

Фронтенд раздаётся через nginx-контейнер в `docker-compose.yml`.

```bash
docker compose up -d
```

Откройте http://localhost

### Вариант 2: Локальная разработка (без Docker)

Запустите Go-сервер и прокси-сервер `server.py`:

```bash
# Терминал 1: бэкенд
go run . -t kdl

# Терминал 2: фронтенд + прокси
py server.py
```

Откройте http://localhost:5500

## API

Фронтенд обращается к бэкенду через следующие эндпоинты:

- `GET /museum/exhibitions` — список всех экспозиций
- `GET /museum/exhibitions/{id}` — экспозиция с экспонатами
- `GET /museum/news` — список новостей
- `GET /museum/news/{id}` — конкретная новость
- `POST /admin/...` — управление контентом (админка)
