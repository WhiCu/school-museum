#!/bin/bash
# Скрипт запуска бэкенда и фронтенда школьного музея

set -e

echo "=================================="
echo "  Школьный музей — Запуск проекта"
echo "=================================="

# Запуск бэкенда (docker-compose)
echo ""
echo "[1/2] Запуск бэкенда (docker-compose up)..."
docker-compose up -d --build

echo ""
echo "Ожидание готовности сервера..."
for i in $(seq 1 30); do
    if curl -s http://localhost:8080/ping > /dev/null 2>&1; then
        echo "✓ Бэкенд готов!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "⚠ Бэкенд не ответил за 30 секунд. Проверьте логи: docker-compose logs server"
    fi
    sleep 1
done

# Запуск фронтенда
echo ""
echo "[2/2] Запуск фронтенда (HTTP-сервер на порту 5500)..."
cd frontend
python -m http.server 5500 &
FRONTEND_PID=$!
cd ..

echo ""
echo "=================================="
echo "  Всё запущено!"
echo "=================================="
echo ""
echo "  Фронтенд:  http://localhost:5500"
echo "  API:        http://localhost:8080"
echo "  Jaeger:     http://localhost:16686"
echo "  Umami:      http://localhost:3000"
echo ""
echo "  Для остановки: ./stop.sh"
echo "=================================="

# Сохраняем PID фронтенда
echo $FRONTEND_PID > .frontend.pid

wait
