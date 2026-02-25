#!/bin/bash
# Остановка проекта школьного музея (Docker)

echo "Остановка проекта..."
docker compose down
echo "✓ Все сервисы остановлены"
echo "Готово!"
