#!/bin/bash

CLOUDREVE_URL=${CLOUDREVE_URL}
# 从环境变量获取
CLOUDREVE_USERNAME=${CLOUDREVE_USERNAME}
CLOUDREVE_PASSWORD=${CLOUDREVE_PASSWORD}
# 临时 cookie 文件
COOKIE_FILE="/tmp/cloudreve_cookies.txt"
CLOUDREVE_POLICY_ID=${CLOUDREVE_POLICY_ID}
REMOTE_CONFIG_FILE_ID=${REMOTE_CONFIG_FILE_ID}
# 存储上传文件的信息
UPLOADED_FILES=()
# 存储文件外链
LINUX_DOWNLOAD_URL=""
WINDOWS_DOWNLOAD_URL=""

# 检查是否存在版本文件
if [ -f "version.json" ]; then
    # 读取版本号
    VERSION=$(cat version.json | jq -r '.version')
else
    echo "版本文件不存在"
    exit 1
fi

if [ -z "$CLOUDREVE_USERNAME" ] || [ -z "$CLOUDREVE_PASSWORD" ]; then
    echo "CLOUDREVE_USERNAME 或 CLOUDREVE_PASSWORD 未设置"
    exit 1
fi

# 运行 make package
make package

# 检查是否成功
if [ $? -ne 0 ]; then
    echo "构建失败"
    exit 1
fi

# 登录 Cloudreve
echo "正在登录 Cloudreve..."
LOGIN_RESPONSE=$(curl -s -c "$COOKIE_FILE" -X POST "$CLOUDREVE_URL/api/v3/user/session" \
    -H "Content-Type: application/json" \
    -d "{\"userName\":\"$CLOUDREVE_USERNAME\",\"Password\":\"$CLOUDREVE_PASSWORD\"}")

# 检查登录是否成功
if echo "$LOGIN_RESPONSE" | grep -q "\"code\":0"; then
    echo "登录成功"
else
    echo "登录失败: $LOGIN_RESPONSE"
    exit 1
fi

# 打印 cookie 内容
echo "Cookie 内容:"
cat "$COOKIE_FILE"

# 定义上传文件函数
upload_file() {
    local file_path=$1
    local file_name=$(basename "$file_path")
    local platform=$2  # 新增平台参数，用于区分 Linux 和 Windows
    
    echo "开始上传文件: $file_path (平台: $platform)"
    
    # 获取文件大小
    local file_size=$(stat -c%s "$file_path")
    local current_time=$(date +%s000)
    
    # 生成唯一文件名，避免会话冲突
    local unique_suffix=$(date +%s)
    local unique_file_name="${file_name%.zip}-${unique_suffix}.zip"
    
    echo "使用唯一文件名: $unique_file_name"
    
    # 创建上传会话
    echo "正在创建上传会话..."
    local upload_json="{\"last_modified\":$current_time,\"mime_type\":\"application/zip\",\"name\":\"$unique_file_name\",\"path\":\"/logsnap/v$VERSION\",\"policy_id\":\"$CLOUDREVE_POLICY_ID\",\"size\":$file_size}"
    echo "上传请求: $upload_json"
    
    local upload_response=$(curl -s -b "$COOKIE_FILE" -X PUT "$CLOUDREVE_URL/api/v3/file/upload" \
        -H "Content-Type: application/json" \
        -d "$upload_json")
    
    echo "上传会话响应: $upload_response"
    
    # 检查是否成功创建上传会话
    if ! echo "$upload_response" | grep -q "\"code\":0"; then
        # 尝试清理可能存在的会话
        echo "创建上传会话失败，尝试清理现有会话..."
        
        # 获取现有上传会话列表
        local sessions_response=$(curl -s -b "$COOKIE_FILE" -X GET "$CLOUDREVE_URL/api/v3/file/upload")
        echo "现有会话列表: $sessions_response"
        
        # 重新尝试创建会话，使用更加唯一的文件名
        unique_suffix=$(date +%s)$RANDOM
        unique_file_name="${file_name%.zip}-${unique_suffix}.zip"
        
        upload_json="{\"last_modified\":$current_time,\"mime_type\":\"application/zip\",\"name\":\"$unique_file_name\",\"path\":\"/logsnap/v$VERSION\",\"policy_id\":\"$CLOUDREVE_POLICY_ID\",\"size\":$file_size}"
        echo "重新尝试上传请求: $upload_json"
        
        upload_response=$(curl -s -b "$COOKIE_FILE" -X PUT "$CLOUDREVE_URL/api/v3/file/upload" \
            -H "Content-Type: application/json" \
            -d "$upload_json")
        
        echo "重新尝试上传会话响应: $upload_response"
        
        # 再次检查是否成功
        if ! echo "$upload_response" | grep -q "\"code\":0"; then
            echo "重新尝试创建上传会话仍然失败"
            return 1
        fi
    fi
    
    # 解析上传文件响应，以获取到 sessionID
    local session_id=$(echo "$upload_response" | jq -r '.data.sessionID')
    echo "获取到会话ID: $session_id"
    
    # 上传文件内容
    echo "开始上传文件内容..."
    echo "文件大小: $file_size 字节"
    
    # 使用 curl 的 --data-binary 选项直接上传文件内容
    local upload_content_response=$(curl -s -b "$COOKIE_FILE" \
        -X POST "$CLOUDREVE_URL/api/v3/file/upload/$session_id/0" \
        -H "Content-Type: application/octet-stream" \
        -H "Content-Length: $file_size" \
        --data-binary @"$file_path")
    
    echo "上传文件内容响应: $upload_content_response"
    
    # 检查上传是否成功
    if echo "$upload_content_response" | grep -q "\"code\":0"; then
        echo "文件 $file_path 上传成功"
        # 将文件名添加到上传文件列表
        UPLOADED_FILES+=("$unique_file_name")
        
        # 获取文件外链
        get_file_link "$unique_file_name" "$platform"
        
        return 0
    else
        echo "文件 $file_path 上传失败"
        return 1
    fi
}

# 获取文件外链函数
get_file_link() {
    local file_name=$1
    local platform=$2  # 新增平台参数
    local encoded_name=$(echo "$file_name" | sed 's/ /%20/g')
    
    echo "正在搜索文件: $file_name"
    # 搜索文件
    local search_response=$(curl -s -b "$COOKIE_FILE" -X GET "$CLOUDREVE_URL/api/v3/file/search/keywords%2F$encoded_name")
    echo "搜索响应: $search_response"
    
    # 检查搜索是否成功
    if ! echo "$search_response" | grep -q "\"code\":0"; then
        echo "搜索文件失败"
        return 1
    fi
    
    # 解析搜索响应，获取文件 ID
    local file_id=$(echo "$search_response" | jq -r '.data.objects[0].id')
    if [ "$file_id" == "null" ] || [ -z "$file_id" ]; then
        echo "未找到文件 ID"
        return 1
    fi
    
    echo "文件 ID: $file_id"
    
    # 获取文件的外链
    local link_response=$(curl -s -b "$COOKIE_FILE" -X POST "$CLOUDREVE_URL/api/v3/file/source" \
        -H "Content-Type: application/json" \
        -d "{\"items\":[\"$file_id\"]}")
    
    echo "外链响应: $link_response"
    
    # 解析外链响应，获取外链
    local link=$(echo "$link_response" | jq -r '.data[0].url')
    if [ "$link" == "null" ] || [ -z "$link" ]; then
        echo "未找到文件外链"
        return 1
    fi
    
    echo "文件 $file_name 的外链: $link"
    
    # 根据平台保存外链
    if [ "$platform" == "linux" ]; then
        LINUX_DOWNLOAD_URL="$link"
        echo "保存 Linux 下载链接: $LINUX_DOWNLOAD_URL"
    elif [ "$platform" == "windows" ]; then
        WINDOWS_DOWNLOAD_URL="$link"
        echo "保存 Windows 下载链接: $WINDOWS_DOWNLOAD_URL"
    fi
    
    return 0
}

# 上传 Linux 版本
echo "上传 Linux 版本..."
upload_file "dist/logsnap-linux-amd64.zip" "linux"
if [ $? -ne 0 ]; then
    echo "Linux 版本上传失败"
    rm -f "$COOKIE_FILE"
    exit 1
fi

# 上传 Windows 版本
echo "上传 Windows 版本..."
upload_file "dist/logsnap-windows-amd64.zip" "windows"
if [ $? -ne 0 ]; then
    echo "Windows 版本上传失败"
    rm -f "$COOKIE_FILE"
    exit 1
fi

# 下载远程配置文件
echo "正在下载远程配置文件..."
# 获取下载链接
DOWNLOAD_LINK_RESPONSE=$(curl -s -b "$COOKIE_FILE" -X PUT "$CLOUDREVE_URL/api/v3/file/download/$REMOTE_CONFIG_FILE_ID")
echo "获取下载链接响应: $DOWNLOAD_LINK_RESPONSE"

# 检查是否成功获取下载链接
if ! echo "$DOWNLOAD_LINK_RESPONSE" | grep -q "\"code\":0"; then
    echo "获取下载链接失败"
    rm -f "$COOKIE_FILE"
    exit 1
fi

# 解析下载链接
DOWNLOAD_LINK=$(echo "$DOWNLOAD_LINK_RESPONSE" | jq -r '.data')
echo "下载链接: $DOWNLOAD_LINK"

# 使用下载链接获取文件内容
CONFIG_RESPONSE=$(curl -s "$CLOUDREVE_URL$DOWNLOAD_LINK")
# echo "配置文件内容: $CONFIG_RESPONSE"

# 更新配置文件内容
echo "正在更新配置文件内容..."
UPDATED_CONFIG=$(echo "$CONFIG_RESPONSE" | jq \
    --arg version "$VERSION" \
    --arg linux_url "$LINUX_DOWNLOAD_URL" \
    --arg windows_url "$WINDOWS_DOWNLOAD_URL" \
    '.latest_versions.stable = $version | .download_urls.stable.linux = $linux_url | .download_urls.stable.windows = $windows_url')

# echo "更新后的配置文件内容: $UPDATED_CONFIG"

# 将更新后的配置文件上传回去
echo "正在上传更新后的配置文件..."
UPDATE_RESPONSE=$(curl -s -b "$COOKIE_FILE" -X PUT "$CLOUDREVE_URL/api/v3/file/update/$REMOTE_CONFIG_FILE_ID" \
    -H "Content-Type: application/json" \
    -d "$UPDATED_CONFIG")

echo "更新配置文件响应: $UPDATE_RESPONSE"

# 检查更新是否成功
if echo "$UPDATE_RESPONSE" | grep -q "\"code\":0"; then
    echo "配置文件更新成功"
else
    echo "配置文件更新失败: $UPDATE_RESPONSE"
    rm -f "$COOKIE_FILE"
    exit 1
fi

# 清理 cookie 文件
rm -f "$COOKIE_FILE"
echo "所有文件发布完成"
