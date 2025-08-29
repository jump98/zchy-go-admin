# 启动远程容器命令

## 启动 api 服务

```sh
  # 根据端口号启动不同的服务
  docker run -d -p 8889:8888 --name zchy-api-a e1

```

## 同步 API 接口

```sh
  #对于不存在的 api 接口会插入到数据库
  docker run -d -p 8889:8888 --name zchy-api-a 88b ./go-admin server -c /app/config/settings.prod.yml -a true
```

## 同步数据库

```sh
  #！！！非初始化项目不可同步数据库，否则必出问题
  docker run -d -p 8889:8888 --name zchy-api-a 88b ./go-admin migrate -c /app/config/settings.prod.yml
```

<!-- # docker run -d -p 8898:8888 --name zchy-api-a 597 server -c /app/config/settings.dev.yml -->
<!-- # docker run -d -p 8889:8888 --name zchy-api-a 59746bd14b8c ./go-admin server -c /app/config/settings.dev.yml -->
