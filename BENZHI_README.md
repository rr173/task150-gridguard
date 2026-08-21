# task150-gridguard Benzhi 评测说明

GridGuard 是一个配电保护整定协调服务。工程师登记馈线、保护区和继电器整定，创建候选计划后由服务计算主保护/后备保护动作时间和选择性裕度；只有无阻断违例的计划能被激活。SQLite 持久化所有工程实体、版本、协调对、违例和故障事件，重启同一数据库会恢复活动计划。

## 本地构建、运行与测试

```bash
GOTOOLCHAIN=local go build ./...
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
GOTOOLCHAIN=local go run ./cmd/gridguard --addr :8080 --db gridguard.db
GOTOOLCHAIN=local go run ./cmd/gridguard --smoke-test --db :memory:
```

服务启动后可访问 `/healthz`、`/readyz`，通过 `/v1/feeders/create`、`/v1/zones/create`、`/v1/relays/create`、`/v1/plans/create`、`/v1/plans/validate` 和 `/v1/plans/activate` 完成真实的协调闭环。

## Docker

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh gridguard:amd64 linux/amd64
./build_benzhi_docker.sh gridguard:arm64 linux/arm64
docker run --rm gridguard:amd64 go run ./cmd/gridguard --smoke-test --db :memory:
```

构建脚本的第一个参数是镜像名，第二个参数是平台。镜像进入 bash，便于评测环境执行构建、测试和 smoke-test；项目无需外部数据库或其他服务。
