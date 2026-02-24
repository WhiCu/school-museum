# Фронтенд музея лицея №76

Простой фронтенд для сайта школьного музея.

## Структура

```
frontend/
├── index.html          # Главная страница
├── exhibition.html     # Страница экспозиции
├── css/
│   └── styles.css      # Стили
└── js/
    ├── api.js          # API клиент
    ├── main.js         # Скрипт главной страницы
    └── exhibition.js   # Скрипт страницы экспозиции
```

## Запуск

### Вариант 1: Простой HTTP сервер (Python)

```bash
cd frontend
python -m http.server 3000
```

Откройте http://localhost:3000

### Вариант 2: Live Server (VS Code)

1. Установите расширение "Live Server" в VS Code
2. Откройте `frontend/index.html`
3. Нажмите "Go Live" в правом нижнем углу

### Вариант 3: Node.js (serve)

```bash
npm install -g serve
cd frontend
serve -p 3000
```

## API

Фронтенд ожидает бэкенд на `http://localhost:8080` со следующими эндпоинтами:

- `GET /museum/exhibitions` - список всех экспозиций
- `GET /museum/exhibitions/{id}` - экспозиция с экспонатами
- `GET /museum/news` - список новостей
- `GET /museum/news/{id}` - конкретная новость

## CORS

Для локальной разработки убедитесь, что бэкенд разрешает CORS запросы с `http://localhost:3000`.
