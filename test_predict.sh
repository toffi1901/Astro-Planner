#!/bin/bash

BASE_URL="http://localhost:8080/api"
EMAIL="predict_${RANDOM}@test.com"
PASS="secret123"



echo " ТЕСТИРОВАНИЕ PREDICT (пролёты)"
echo ""

echo "0. ПОДГОТОВКА (регистрация, логин)"

curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" > /dev/null

LOGIN_RES=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}")

TOKEN=$(echo "$LOGIN_RES" | jq -r '.access_token')
if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "Токен не получен"
    exit 1
fi
echo "Токен получен"
echo ""

echo "1. PREDICT БЕЗ ТОКЕНА (Ожидается: 401)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-23T00:00:00Z"}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 401 ]; then
    echo "error: пустил без токена"
    cat /tmp/predict_res.json
    exit 1
fi
echo "Успех (доступ закрыт)"
echo ""

echo "2. PREDICT БЕЗ КООРДИНАТ (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-23T00:00:00Z"}')
echo "   Статус: $STATUS"
ERROR_MSG=$(cat /tmp/predict_res.json | jq -r '.error // ""')
echo "   Ошибка: $ERROR_MSG"
if [ "$STATUS" -ne 400 ]; then
    echo "error: ожидался 400 (нет координат)"
    exit 1
fi
if [[ "$ERROR_MSG" != *"coordinates"* ]]; then
    echo "Статус верный, но текст ошибки не про координаты"
fi
echo "Успех (сервер требует координаты)"
echo ""

echo "3. УСТАНОВКА КООРДИНАТ (Москва: 55.75, 37.62)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X PUT "$BASE_URL/profile" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"latitude": 55.75, "longitude": 37.62}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 200 ]; then
    echo "error: координаты не установлены"
    cat /tmp/predict_res.json
    exit 1
fi
echo "Успех"
echo ""

echo "4. PREDICT МКС (25544) на 1 сутки (Ожидается: 200)"
RES=$(curl -s -w "\n%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-23T00:00:00Z"}')
STATUS=$(echo "$RES" | tail -n1)
BODY=$(echo "$RES" | head -n-1)

echo "   Статус: $STATUS"
if [ "$STATUS" -ne 200 ]; then
    echo "error"
    echo "   Ответ: $BODY"
    exit 1
fi

COUNT=$(echo "$BODY" | jq -r '.count')
echo "   Найдено пролётов: $COUNT"
if [ "$COUNT" == "null" ]; then
    echo "Поле count отсутствует в ответе"
fi
echo "Успех (ответ корректный)"
echo ""

echo "5. PREDICT С НЕВАЛИДНЫМ ИНТЕРВАЛОМ (end < start) (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-23T00:00:00Z", "end_time": "2026-10-22T00:00:00Z"}')
echo "   Статус: $STATUS"
ERROR_MSG=$(cat /tmp/predict_res.json | jq -r '.error // ""')
echo "   Ошибка: $ERROR_MSG"
if [ "$STATUS" -ne 400 ]; then
    echo "error: сервер принял end < start"
    exit 1
fi
echo "Успех (сервер отклонил)"
echo ""

echo "6. PREDICT С ИНТЕРВАЛОМ 8 ДНЕЙ (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-30T00:00:00Z"}')
echo "   Статус: $STATUS"
ERROR_MSG=$(cat /tmp/predict_res.json | jq -r '.error // ""')
echo "   Ошибка: $ERROR_MSG"
if [ "$STATUS" -ne 400 ]; then
    echo "error: сервер принял интервал > 7 дней"
    exit 1
fi
echo "Успех (сервер отклонил)"
echo ""

echo "7. PREDICT ДЛЯ НЕСУЩЕСТВУЮЩЕГО СПУТНИКА (norad_id=99999) (Ожидается: 404)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/99999/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-23T00:00:00Z"}')
echo "   Статус: $STATUS"
ERROR_MSG=$(cat /tmp/predict_res.json | jq -r '.error // ""')
echo "   Ошибка: $ERROR_MSG"
if [ "$STATUS" -ne 404 ]; then
    echo "error: сервер не отклонил запрос"
    exit 1
fi
echo "Успех (спутник не найден)"
echo ""

echo "8. PREDICT С НЕВАЛИДНЫМ norad_id (Ожидается: 404)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/abc/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"start_time": "2026-10-22T00:00:00Z", "end_time": "2026-10-23T00:00:00Z"}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 404 ]; then
    echo "error: ожидался 400"
    cat /tmp/predict_res.json
    exit 1
fi
echo "Успех (norad_id отклонён)"
echo ""

echo "9. PREDICT С ПУСТЫМ ТЕЛОМ (Ожидается: 400)"
STATUS=$(curl -s -o /tmp/predict_res.json -w "%{http_code}" \
    -X POST "$BASE_URL/satellites/25544/predict" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}')
echo "   Статус: $STATUS"
if [ "$STATUS" -ne 400 ]; then
    echo "error: сервер принял пустое тело"
    cat /tmp/predict_res.json
    exit 1
fi
echo "Успех (пустые даты отклонены)"
echo ""

echo " 🎉 ВСЕ ТЕСТЫ PREDICT ПРОЙДЕНЫ!"
