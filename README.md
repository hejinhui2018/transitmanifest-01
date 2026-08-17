# TransitManifest

TransitManifest 是面向干线和转运场站的运输清单服务。它管理车辆班次、包裹装卸扫描、交接签收、异常件、路由、状态查询、审计和日报，并提供 `transitmanifest` CLI。

## 设计说明

- 核心存储是磁盘上的 `events.log` 追加日志。每条 JSON 记录包含序号、前序 SHA-256 校验和与自身校验和；打开服务时会完整验证链，发现截断、篡改或乱序会拒绝启动。
- `snapshot.json` 是已验证事件位置的原子快照。写入临时文件、同步文件、替换目标并同步目录，进程重启后加载快照并回放后续日志。
- `storage.Store` 以读写锁保护追加和快照位置；服务状态使用快照副本，查询不会拿到可修改的内部 map。
- 扫描以 `scan_id` 做幂等键，重复的包裹操作会生成异常件；关闭清单允许迟到的物理扫描，但重复扫描仍必须被拒绝。
- `audit` 从已验证日志生成时间线，`report` 从事件生成日报，二者不依赖易失内存统计。

## 运行

需要 Go 1.23，数据目录必须是绝对路径。默认数据目录为当前目录的 `data`，默认操作员为 `cli`。

```text
set TRANSITMANIFEST_DATA_DIR=C:\data\transitmanifest
set TRANSITMANIFEST_ACTOR=station-a
go run ./cmd/transitmanifest create --id m-001 --trip t-001 --plate plate-001 --origin hub-a --destination hub-b
go run ./cmd/transitmanifest scan --manifest m-001 --scan s-001 --package p-001 --station hub-a --operation load
go run ./cmd/transitmanifest close --manifest m-001 --reason departed
go run ./cmd/transitmanifest report --date 2026-08-18
go run ./cmd/transitmanifest verify
```

## 质量门禁

`Makefile` 的 `local` 目标使用本机 Go；容器目标固定使用官方 `golang:1.23.12`。生产构建强制 `CGO_ENABLED=0`，`scripts/build.sh` 同时构建 `linux/amd64` 和 `linux/arm64`。

```text
make local
make docker-test
./scripts/build.sh
```

## 提交链与专项

本候选仓库保留三段提交：`healthy` 是完整基线，`main` 在快照恢复中故意遗漏已处理的 `scan_id` 并加入跨 `storage/manifest/report` 的专项测试，`gold_model_fix` 只修复生产恢复逻辑。专项命令为：

```text
go test ./integration -run TestClosedManifestDuplicateScanSurvivesRestart -count=20
```

在 `main` 上该命令应稳定失败，在 `gold_model_fix` 上应稳定通过；没有 `test_model_fix` 分支或提交。

