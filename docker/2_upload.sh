#/bin/bash

REMOTE_PATH="/root/golang/images" #存储路径
REMOTE_URL="root@47.115.214.123" #远程服务器


# 进入 images 目录
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)  # 获取脚本所在目录
IMAGE_DIR="$SCRIPT_DIR/images"


shopt -s nullglob  # Bash 扩展，匹配不到时返回空数组
FILES=("$IMAGE_DIR"/*.tar)

if [ ${#FILES[@]} -eq 0 ]; then
    echo "没有找到任何 tar 文件在 $IMAGE_DIR"
    exit 1
fi

echo "找到以下 tar 文件："

# 1 列出文件并编号
i=1
for f in "${FILES[@]}"; do
    echo "$i) $(basename "$f")"
    i=$((i + 1))
done

# 提示用户选择
while true; do
    read -p "请输入你要选择的文件编号: " choice
    if [ "$choice" -ge 1 ] 2>/dev/null && [ "$choice" -lt "$i" ]; then
        SELECTED_FILE="${FILES[$((choice - 1))]}"
        echo "你选择的文件是: $(basename "$SELECTED_FILE")"
        break
    else
        echo "输入无效，请重新输入编号."
    fi
done


# 2 上传到远程服务器
echo "正在上传镜像到远程服务器..."
scp "${SELECTED_FILE}" "${REMOTE_URL}:${REMOTE_PATH}/"
if [ $? -ne 0 ]; then
    echo "上传失败！"
    exit 1
fi
echo "上传完成: ${REMOTE_URL}/" $(basename "$SELECTED_FILE")

# 3 可选：在远程服务器加载镜像
read -p "是否在远程服务器加载镜像？(y/n): " load_remote
if [[ "$load_remote" =~ ^[Yy]$ ]]; then
    echo "在远程服务器加载镜像..."
    TAR_FILE=$(basename "$SELECTED_FILE")
    ssh "${REMOTE_URL}" "docker load -i ${REMOTE_PATH}/${TAR_FILE}"
    if [ $? -ne 0 ]; then
        echo "远程加载镜像失败！"
        exit 1
    fi
    echo "远程镜像加载完成！"
fi