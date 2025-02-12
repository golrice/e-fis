# 分布式缓存系统

## 概述

分布式缓存系统是一个简单的缓存解决方案，支持缓存数据的分布式存储和访问。此项目使用Go语言编写，包含了服务端、客户端和一些辅助组件。

## 项目结构

```
e-fis/
├── bin/
├── cmd/
│   ├── client/
│   │   └── main.go
│   └── server/
│       ├── info.go
│       ├── main.go
│       └── pool.go
├── internal/
│   ├── cache/
│   │   ├── byteview.go
│   │   ├── cache.go
│   │   ├── cache_test.go
│   │   └── flowcontrol/
│   │       └── controler.go
│   ├── consistenthash/
│   │   └── hash.go
│   ├── peer/
│   │   ├── getter.go
│   │   └── peer.go
│   └── protocal/
│       ├── cachepb.pb.go
│       └── cachepb.proto
├── Makefile
└── run.sh
```

## 功能特性

- **缓存服务**：可以根据指定的端口启动缓存服务。
- **客户端**：提供简单的客户端命令行工具，用于与缓存服务进行交互。
- **分布式缓存**：通过一致性哈希算法实现缓存数据的分布式存储。
- **缓存淘汰策略**：使用LRU算法进行缓存数据的淘汰。
- **缓存一致性控制**：通过控制缓存的并发访问，确保缓存的一致性。

## 安装与运行

### 安装依赖

确保已安装`go`和`protobuf`编译器。如果尚未安装，请参考[Go安装指南](https://golang.org/doc/install)和[protobuf安装指南](https://developers.google.com/protocol-buffers/docs/gotutorial)。

然后，安装项目依赖：

```bash
go mod tidy
```

### 编译项目

使用`Makefile`进行编译：

```bash
make build
```

### 启动服务

启动多个缓存服务实例：

```bash
make run-server
```

启动API服务实例：

```bash
make run-client
```

### 运行测试

项目包含单元测试，可以使用以下命令运行测试：

```bash
make test
```

### 代码格式化

使用`go fmt`格式化代码：

```bash
make fmt
```

### 代码静态检查

使用`golangci-lint`对代码进行静态检查：

```bash
make lint
```

### 清理构建文件

使用以下命令清理构建生成的二进制文件：

```bash
make clean
```

## 贡献

欢迎贡献代码和文档。在提交PR之前，请确保代码通过所有测试并符合项目编码规范。

## 许可证

本项目使用MIT许可证。更多细节请参见[LICENSE](LICENSE)文件。
