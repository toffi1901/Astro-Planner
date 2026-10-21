#!/bin/bash

URL="http://localhost:8080/api/auth/login"
ORIGIN="http://localhost:5500"

echo "=== 1. Preflight (OPTIONS) ==="
curl -s -i -X OPTIONS "$URL" \
  -H "Origin: $ORIGIN" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  | grep -i "access-control\|HTTP/"

echo ""
echo "=== 2. Реальный запрос (POST) ==="
curl -s -i -X POST "$URL" \
  -H "Origin: $ORIGIN" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"secret123"}' \
  | grep -i "access-control\|set-cookie\|HTTP/"

echo ""
echo "=== 3. Чужой Origin (должно быть пусто/без Allow-Origin) ==="
curl -s -i -X OPTIONS "$URL" \
  -H "Origin: http://evil.com" \
  -H "Access-Control-Request-Method: POST" \
  | grep -i "access-control\|HTTP/"