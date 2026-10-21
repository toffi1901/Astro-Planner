#!/bin/bash

BASE_URL="http://localhost:8080/api"
EMAIL="profile_${RANDOM}@test.com"
PASS="secret123"

if ! command -v jq &> /dev/null; then
    echo "ERROR: Установи jq: sudo apt install jq"
    exit 1
fi

echo " ТЕСТИРОВАНИЕ ПРОФИЛЯ"

echo ""

echo "0. ПОДГОТОВКА (регистрация + логин)"

curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" > /dev/null

LOGIN_RES=$(curl -s -c /tmp/profile_cookies.txt -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}")

TOKEN=$(echo "$LOGIN_RES" | jq -r '.access_token')
if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "Не удалось получить токен"
    echo "   Ответ: $LOGIN_RES"
    exit 1
fi
echo "Токен получен"
echo ""

echo "1. GET /api/profile БЕЗ ТОКЕНА (Ожидается: 401)"
STATUS=$(curl -s -o /tmp/profile_res.json -w "%{http_code}" \
    -X GET "$BASE_URL/profile")
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 401 ]; then
    echo "ERROR сервер пустил без токена"
    cat /tmp/profile_res.json
    exit 1
fi
echo "Успех (доступ закрыт)"
echo ""

echo "2. GET /api/profile С ТОКЕНОМ (Ожидается: 200, координаты = null)"
RES=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/profile" \
    -H "Authorization: Bearer $TOKEN")
STATUS=$(echo "$RES" | tail -n1)
BODY=$(echo "$RES" | head -n-1)

echo "   Статус: $STATUS"
if [ "$STATUS" -ne 200 ]; then
    echo "ERRORОШИБКА"
    echo "   Ответ: $BODY"
    exit 1
fi

LAT=$(echo "$BODY" | jq -r '.latitude')
LON=$(echo "$BODY" | jq -r '.longitude')
if [ "$LAT" != "null" ] || [ "$LON" != "null" ]; then
    echo "Координаты не null (возможно, сохранены с прошлого теста)"
    echo "   latitude=$LAT, longitude=$LON"
else
    echo " Успех (новый юзер, координат нет)"
fi
echo ""

echo "3. PUT /api/profile (Москва: 55.75, 37.62)"
STATUS=$(curl -s -o /tmp/profile_res.json -w "%{http_code}" \
    -X PUT "$BASE_URL/profile" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"latitude": 55.75, "longitude": 37.62}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 200 ]; then
    echo "ERROR"
    cat /tmp/profile_res.json
    exit 1
fi
echo "Успех (координаты сохранены)"
echo ""

echo "4. GET /api/profile (проверка сохранения)"
BODY=$(curl -s -X GET "$BASE_URL/profile" -H "Authorization: Bearer $TOKEN")
LAT=$(echo "$BODY" | jq -r '.latitude')
LON=$(echo "$BODY" | jq -r '.longitude')

if [ "$LAT" != "55.75" ] || [ "$LON" != "37.62" ]; then
    echo "ERROR: координаты не сохранились"
    echo "Ожидалось: 55.75, 37.62"
    echo "Получено: $LAT, $LON"
    exit 1
fi
echo "Успех (координаты: $LAT, $LON)"
echo ""

echo "5. PUT с невалидной широтой (91) (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/profile_res.json -w "%{http_code}" \
    -X PUT "$BASE_URL/profile" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"latitude": 91, "longitude": 37.62}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 400 ]; then
    echo "ERROR: сервер принял широту 91"
    cat /tmp/profile_res.json
    exit 1
fi
echo "Успех (широта отклонена)"
echo ""

echo "6. PUT с невалидной долготой (181) (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/profile_res.json -w "%{http_code}" \
    -X PUT "$BASE_URL/profile" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"latitude": 55.75, "longitude": 181}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 400 ]; then
    echo "ERROR: сервер принял долготу 181"
    cat /tmp/profile_res.json
    exit 1
fi
echo " Успех (долгота отклонена)"
echo ""

echo "7. PUT БЕЗ ТОКЕНА (Ожидается: 401)"
STATUS=$(curl -s -o /tmp/profile_res.json -w "%{http_code}" \
    -X PUT "$BASE_URL/profile" \
    -H "Content-Type: application/json" \
    -d '{"latitude": 55.75, "longitude": 37.62}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 401 ]; then
    echo "ERROR"
    exit 1
fi
echo "Успех (доступ закрыт)"
echo ""

echo " 🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ!"
