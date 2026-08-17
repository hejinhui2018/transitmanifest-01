# Bug 复现说明

## 分支

- `healthy`：无回归测试的健康项目基线。
- `main`：正式含 Bug 基线和专项回归测试。
- `gold_model_fix`：独立标准修复分支。
- `test_model_fix`：模型修复后的交付分支。

## 问题

清单事件保存在追加日志中，快照保存已应用的状态。扫描的 `scan_id` 是幂等键：相同扫描在同一进程内必须被拒绝，关闭清单重开后也必须被拒绝。含 Bug 的快照恢复路径重建了清单计数，却把已保存的扫描索引清空；重启后同一个卸货事件会被再次接受，并污染异常判断和日报计数。

## 复现命令

```sh
go test ./integration -run '^TestClosedManifestDuplicateScanSurvivesRestart$' -count=1 -v
```

含 Bug 基线退出 1，失败表现为重启后的重复扫描未返回 `ErrDuplicateScan`。正确修复必须保留快照中的扫描索引，并继续回放快照之后的日志。

## 通过条件

模型只能修改生产源代码，不能修改测试、模块文件或交付文件。修复后必须使用同一专项命令，并通过以下完整门禁：

```sh
go test ./... -count=1
go vet ./...
go build ./...
```

命令在 Linux amd64、Linux arm64、Go 1.23.12、`GOTOOLCHAIN=local`、`CGO_ENABLED=0` 环境中执行。

