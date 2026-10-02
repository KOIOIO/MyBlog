// 服务入口：加载配置 → bootstrap 组装基础设施 → 启动 Cron 与 HTTP。
//
// 迁移状态（Phase 0）：基础设施实例化已收敛到 bootstrap（构造注入），
// 全局变量由 bootstrap 写入作为兼容层；路由/Cron/服务启动暂复用旧实现，
// 后续 Phase 逐步替换为 internal/interface/http 与 internal/interface/cron。
package main

import (
	"server/core"
	"server/flag"
	"server/initialize"
	"server/internal/bootstrap"
)

func main() {
	cfg := bootstrap.InitConfig()
	zl := bootstrap.InitLogger(cfg)
	infra := bootstrap.InitInfra(cfg, zl)

	defer infra.Redis.Close()

	flag.InitFlag()

	initialize.InitCron()

	core.RunServer()
}
