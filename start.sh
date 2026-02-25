#!/bin/bash
# Запуск проекта школьного музея (Docker)

set -e

echo "=================================="
echo "  Школьный музей — Запуск"
echo "=================================="

echo ""
echo "Запуск всех сервисов..."
docker compose up -d --build

echo ""
echo "Ожидание готовности сервера..."
for i in $(seq 1 90); do
    if curl -s 'http://localhost:8081/ping?message=ping' > /dev/null 2>&1; then
        echo "✓ Бэкенд готов!"
        break
    fi
    if [ "$i" -eq 90 ]; then
        echo "⚠ Бэкенд не ответил за 90с. Проверьте: docker compose logs server"
    fi
    sleep 1
done

echo ""
echo "=================================="
echo "  Всё запущено!"
echo "=================================="
echo ""
echo "  Фронтенд:  http://localhost"
echo "  API:        http://localhost:8081"
echo ""
echo "  Для остановки: ./stop.sh"
echo "=================================="
