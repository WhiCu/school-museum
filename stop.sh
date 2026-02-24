#!/bin/bash
# Остановка бэкенда и фронтенда

echo "Остановка проекта..."

# Остановка фронтенда
if [ -f .frontend.pid ]; then
    kill $(cat .frontend.pid) 2>/dev/null
    rm .frontend.pid
    echo "✓ Фронтенд остановлен"
fi

# Остановка бэкенда
docker-compose down
echo "✓ Бэкенд остановлен"

echo "Готово!"
