#!/bin/bash

BASE_URL="http://localhost:8080/api"
COOKIE_FILE="cookie.txt"

EMAIL="user_$(date +%s)@test.com"
PASS="secret123"

echo " ТЕСТИРОВАНИЕ АВТОРИЗАЦИИ"
echo ""
echo "1. Register ($EMAIL)"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}")
echo "Статус: $STATUS (Ожидается: 201 Created)"
["$STATUS" -ne 201] && echo "Error :( )" && exit 1
echo "Что делаем? Радуемся"
echo ""

echo "2. ЛОГИН (получаем токены(не для нейронок, а Вы что подумали?))"
LOGIN_RES=$(curl -s -c "$COOKIE_FILE" -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RES" | jq -r '.access_token')
if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "Error: Не удалось получить access_token!"
    echo "   Ответ сервера: $LOGIN_RES"
    exit 1
fi
echo "Получен Access Token: ${ACCESS_TOKEN:0:15}..."
echo "Получен Refresh Token (сохранен в $COOKIE_FILE)"
echo "Что делаем? Радуемся"
echo ""

echo "3. ОБНОВЛЕНИЕ ТОКЕНА (Refresh с помощью Cookie)"
REFRESH_RES=$(curl -s -b "$COOKIE_FILE" -X POST "$BASE_URL/auth/refresh")
NEW_TOKEN=$(echo "$REFRESH_RES" | jq -r '.access_token')

if [ "$NEW_TOKEN" == "null" ] || [ -z "$NEW_TOKEN" ]; then
    echo "Error: Не удалось обновить токен!"
    echo "  Ответ сервера: $REFRESH_RES"
    exit 1
fi
echo "Получен НОВЫЙ Access Token: ${NEW_TOKEN:0:15}..."
echo "Что делаем? Радуемся"
echo ""

echo "4. ЗАЩИЩЕННЫЙ ЗАПРОС (Logout с Access Token)"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $NEW_TOKEN")
echo "   Статус: $STATUS (Ожидалось: 204)"
[ "$STATUS" -ne 204 ] && echo "Error" && exit 1
echo "Успех (Сессия удалена из Redis)"
echo "Что делаем? Радуемся"
echo ""

echo "5. ПОВТОРНЫЙ ЗАПРОС С ТЕМ ЖЕ ТОКЕНОМ (Должен быть 401)"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $NEW_TOKEN")
echo "   Статус: $STATUS (Ожидалось: 401)"
[ "$STATUS" -ne 401 ] && echo "error Сессия не удалилась." && exit 1
echo "Успех (Доступ закрыт)"
echo "Что делаем? Радуемся"
echo ""

echo "6. REFRESH ПОСЛЕ LOGOUT (Должен быть 401)"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -b "$COOKIE_FILE" -X POST "$BASE_URL/auth/refresh")
echo "   Статус: $STATUS (Ожидалось: 401)"
if [ "$STATUS" -ne 401 ]; then echo "error Мертвая сессия смогла обновиться."; exit 1; fi
echo " Успех (Система не дала обновить мертвую сессию)"; echo ""
echo "Что делаем? Радуемся"
echo ""echo " 🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!"
