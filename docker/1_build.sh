#/bin/bash

#copy
cp ../go.mod ./go.mod
cp ../go.sum ./go.sum

#设置镜像版本
export IMAGE_NAME="zchy-go-admin:v1.0.6"

# 设置镜像名称和版本

echo "开始构建 Docker 镜像..."
docker-compose build
echo "Docker 镜像构建完成！"

#清理
rm  ./go.mod
rm  ./go.sum


# 将dockers images 打包成tar文件
IMAGE_NAME2="${IMAGE_NAME//\//-}"
IMAGE_NAME2="${IMAGE_NAME2//./-}"
IMAGE_NAME2="${IMAGE_NAME2//:/-}"
# 生成镜像文件名
TAR_FILE="${IMAGE_NAME2}.tar"
echo "选择镜像: "$TAR_FILE

docker save -o "./images/$TAR_FILE" "$IMAGE_NAME"
if [ $? -ne 0 ]; then
    echo "Docker save 失败！"
    exit 1
fi
echo "镜像已打包为 $TAR_FILE"


