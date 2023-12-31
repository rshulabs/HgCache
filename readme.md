## 项目简介

- 仿 [groupcache](https://github.com/golang/groupcache) 改造的分布式缓存微服务项目
- [复盘文档](https://github.com/rshulabs/HgCache/blob/main/docs/cn/%E9%A1%B9%E7%9B%AE%E5%A4%8D%E7%9B%98.md)
- [线上体验地址](http://rshulabs.v4.idcfengye.com/api?key=Tom)(注：由于服务器是本地电脑（买不起服务器，太穷了），所以大部分时间都可能没有开放，用 k8s 集群部署的实验环境，电脑吃不消~)
  - 接口: GET /api?key=Tom (Tom 为示例)

## 版本

- **Go** go1.20.3 | > 19
- **make** GNU Make 4.3

- **docker-compose** v2.6.1
- **protobuf** libprotoc 3.18.0

## 功能实现

- 基于 groupcache 实现的⼀个分布式缓存，并在此基础上进行了服务发现及注册和网络框架扩展
- 分布式一致性模块使用一致性哈希算法对节点进行负载均衡，实现了 sigleFlight 机制应对缓存击穿问题
- 缓存模块实现过期和回调机制，并与其他模块解耦，可以使用多种淘汰算法，本项目默认是 lru 淘汰算法
- 对缓存值进行封装，实现深拷贝机制，防止被恶意修改
- 网络模块使用 gin，实现了跨域 cors 中间件
- 使用 etcd 进行服务发现和服务注册

## Usage

### _v1_

- ./run.sh

  ```
  2023/08/16 10:14:39 geecache is running at http://localhost:8001
  2023/08/16 10:14:39 geecache is running at http://localhost:8003
  2023/08/16 10:14:39 fontend server is running at http://localhost:8791
  2023/08/16 10:14:39 geecache is running at http://localhost:8002
  >>> start test
  2023/08/16 10:14:41 [Server http://localhost:8003] Pick peer http://localhost:8002
  http://localhost:8002/_hgcache/test/Tom
  2023/08/16 10:14:41 [Server http://localhost:8002] GET /_hgcache/test/Tom
  2023/08/16 10:14:41 [SlowDB] search key Tom
  6302023/08/16 10:14:41 [Server http://localhost:8003] Pick peer http://localhost:8002
  http://localhost:8002/_hgcache/test/Tom
  2023/08/16 10:14:41 [Server http://localhost:8002] GET /_hgcache/test/Tom
  2023/08/16 10:14:41 [HgCache] hit
  2023/08/16 10:14:41 [Server http://localhost:8003] Pick peer http://localhost:8002
  http://localhost:8002/_hgcache/test/Tom
  2023/08/16 10:14:41 [Server http://localhost:8002] GET /_hgcache/test/Tom
  2023/08/16 10:14:41 [HgCache] hit
  ```

### _v2_

- 前提：本地部署好 etcd

  - 参考 cd build_tools/etcd && docker-compose up -d

    ```
    Tips：
    	- ETCD_SERVERS=http://192.168.60.34:2379 // 192.168.60.34 改为本地部署 etcd IP
    	记得给data目录 777 权限
    ```

- ./dist/cache --config=configs/cache.yaml --app.is_start_http=true

  ```
  [GIN-debug] GET    /cache/api                --> github.com/rshulabs/HgCache/internal/cache/controller.Query (3 handlers)
  2023-08-16 10:16:48.395 INFO    controller/http.go:40   http 服务监听地址: 192.168.60.34:8791
  2023-08-16 10:16:48.406 INFO    discovery/registry.go:88        cache server is registered.
  2023-08-16 10:16:48.406 INFO    impl/grpc.go:62 GRPC 服务监听地址: 192.168.60.34:8553
  2023-08-16 10:16:48.407 INFO    discovery/registry.go:106       service is keepalived with 18f89f8f8e566ad
  ```
