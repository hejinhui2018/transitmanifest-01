# TransitManifest 转运清单服务

TransitManifest 面向干线和转运场站，管理车辆班次、包裹装卸扫描、交接签收、异常件、路由、审计和日报。核心数据写入带 SHA-256 链校验的追加日志，并通过原子快照和重启回放恢复状态。

## 工具链

- 官方镜像：`golang:1.23.12`
- `GOTOOLCHAIN=local`
- `CGO_ENABLED=0`
- 验证架构：Linux amd64、Linux arm64

## 标准命令

专项验证命令：

```sh
go test ./integration -run '^TestClosedManifestDuplicateScanSurvivesRestart$' -count=1 -v
```

完整门禁：

```sh
go test ./... -count=1
go vet ./...
go build ./...
```

## 双架构构建

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

Docker 构建阶段使用官方 Go 1.23.12 镜像并执行完整 `go test`、`go vet` 和 `go build` 门禁；运行镜像保留 CLI 入口。

