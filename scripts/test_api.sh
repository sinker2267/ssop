#!/bin/bash

# API测试脚本
API_URL=${1:-"http://localhost:8080/api/v1"}
echo "Testing API at: $API_URL"

# 存储token
TOKEN_FILE=".api_test_token"
REFRESH_TOKEN_FILE=".api_test_refresh_token"

# 临时文件
RESPONSE_FILE=".api_test_response"

# 清理函数
cleanup() {
  rm -f $TOKEN_FILE $REFRESH_TOKEN_FILE $RESPONSE_FILE
}

# 设置退出清理
trap cleanup EXIT

# 测试注册
test_register() {
  echo -e "\n=== 测试用户注册 ==="
  curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "testuser",
      "password": "Test@123",
      "email": "test@example.com",
      "fullName": "测试用户",
      "organization": "测试机构"
    }' > $RESPONSE_FILE

  cat $RESPONSE_FILE | jq .
}

# 测试登录
test_login() {
  echo -e "\n=== 测试用户登录 ==="
  curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "testuser",
      "password": "Test@123"
    }' > $RESPONSE_FILE

  # 提取token
  cat $RESPONSE_FILE | jq -r '.data.token' > $TOKEN_FILE
  cat $RESPONSE_FILE | jq -r '.data.refreshToken' > $REFRESH_TOKEN_FILE
  
  cat $RESPONSE_FILE | jq .
}

# 测试游客登录
test_guest_login() {
  echo -e "\n=== 测试游客登录 ==="
  curl -s -X POST "$API_URL/auth/guest-login" \
    -H "Content-Type: application/json" \
    -d '{}' > $RESPONSE_FILE

  # 提取token
  cat $RESPONSE_FILE | jq -r '.data.token' > $TOKEN_FILE
  cat $RESPONSE_FILE | jq -r '.data.refreshToken' > $REFRESH_TOKEN_FILE
  
  cat $RESPONSE_FILE | jq .
}

# 测试刷新token
test_refresh_token() {
  echo -e "\n=== 测试刷新Token ==="
  
  if [ ! -f $REFRESH_TOKEN_FILE ]; then
    echo "No refresh token found. Run login test first."
    return
  fi
  
  REFRESH_TOKEN=$(cat $REFRESH_TOKEN_FILE)
  
  curl -s -X POST "$API_URL/auth/refresh-token" \
    -H "Content-Type: application/json" \
    -d "{
      \"refreshToken\": \"$REFRESH_TOKEN\"
    }" > $RESPONSE_FILE
    
  # 更新token
  NEW_TOKEN=$(cat $RESPONSE_FILE | jq -r '.data.token')
  if [ "$NEW_TOKEN" != "null" ]; then
    echo $NEW_TOKEN > $TOKEN_FILE
  fi
  
  cat $RESPONSE_FILE | jq .
}

# 测试获取当前用户信息
test_current_user() {
  echo -e "\n=== 测试获取当前用户信息 ==="
  
  if [ ! -f $TOKEN_FILE ]; then
    echo "No token found. Run login test first."
    return
  fi
  
  TOKEN=$(cat $TOKEN_FILE)
  
  curl -s -X GET "$API_URL/users/current" \
    -H "Authorization: Bearer $TOKEN" > $RESPONSE_FILE
    
  cat $RESPONSE_FILE | jq .
}

# 测试登出
test_logout() {
  echo -e "\n=== 测试登出 ==="
  
  if [ ! -f $TOKEN_FILE ]; then
    echo "No token found. Run login test first."
    return
  fi
  
  TOKEN=$(cat $TOKEN_FILE)
  
  curl -s -X POST "$API_URL/auth/logout" \
    -H "Authorization: Bearer $TOKEN" > $RESPONSE_FILE
    
  cat $RESPONSE_FILE | jq .
  
  # 清除token
  rm -f $TOKEN_FILE $REFRESH_TOKEN_FILE
}

# 运行所有测试
echo -e "=== 开始API测试 ===\n"

# 测试注册
test_register

# 测试登录
test_login

# 测试获取当前用户信息
test_current_user

# 测试刷新token
test_refresh_token

# 测试登出
test_logout

# 测试游客登录
test_guest_login

# 测试游客获取当前用户信息
test_current_user

echo -e "\n=== API测试完成 ===\n" 