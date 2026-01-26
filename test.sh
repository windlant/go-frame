#!/bin/bash

echo -e "0.获取所有用户"
curl -X POST http://127.0.0.1:8080/users/list

echo -e "\n\n1.批量创建用户"
curl -X POST "http://127.0.0.1:8080/users/create" \
  -H "Content-Type: application/json" \
  -d '[{"name":"Alice","email":"alice@example.com"},{"name":"Bob","email":"bob@example.com"}]'

echo -e "\n\n2. 按 ID 查询（假设 ID 为 1,2）"
curl -X POST "http://127.0.0.1:8080/users/get" \
  -H "Content-Type: application/json" \
  -d '{"ids":[1,2]}'

echo -e "\n\n3. 按 Email 查询"
curl -X POST "http://127.0.0.1:8080/users/get" \
  -H "Content-Type: application/json" \
  -d '{"emails":["alice@example.com","bob@example.com"]}'

echo -e "\n\n4. 批量更新（假设 ID 1,2 存在）"
curl -X POST "http://127.0.0.1:8080/users/update" \
  -H "Content-Type: application/json" \
  -d '[{"id":1,"name":"Alice Updated","email":"alice.updated@example.com"},{"id":2,"name":"Bob Updated","email":"bob.updated@example.com"}]'

# echo -e "\n\n5. 批量删除（ID 1,2）"
# curl -X POST "http://127.0.0.1:8080/users/delete" \
#   -H "Content-Type: application/json" \
#   -d '{"ids":[10,11]}'

echo -e "\n\n6. 错误测试：空用户列表"
curl -X POST "http://127.0.0.1:8080/users/create" \
  -H "Content-Type: application/json" \
  -d '[]'

echo -e "\n\n7. 错误测试：无效 JSON"
curl -X POST "http://127.0.0.1:8080/users/create" \
  -H "Content-Type: application/json" \
  -d 'invalid json'

