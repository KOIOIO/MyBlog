// 服务入口：加载配置 → bootstrap 组装基础设施与应用服务 → 启动 Cron 与 HTTP。
//
// 迁移状态（Phase 1）：基础设施与应用服务均由 bootstrap 构造注入；
// 路由使用 internal/interface/http（base/user 已迁移，其余 BC 委托旧 router）；
// 中间件已迁移为注入式；global 仅作为旧代码兼容层，写入点收敛于 bootstrap。
package main

import (
	"context"
	"time"

	"server/flag"
	"server/initialize"
	"server/internal/bootstrap"
	ihttp "server/internal/interface/http"

	"github.com/fvbock/endless"
	"go.uber.org/zap"
)

func main() {
	cfg := bootstrap.InitConfig()
	zl := bootstrap.InitLogger(cfg)
	infra := bootstrap.InitInfra(cfg, zl)
	defer infra.Redis.Close()

	deps := bootstrap.BuildApp(infra)

	// 加载全部 JWT 黑名单到本地缓存
	if err := deps.Auth.LoadAll(context.Background()); err != nil {
		zl.Error("Failed to load JWT blacklist:", zap.Error(err))
	}

	flag.InitFlag()

	initialize.InitCron()
	bootstrap.InitCron(deps, zl)

	runServer(cfg.System.Addr(), deps)
}

// runServer 启动 HTTP 服务（行为与 core.RunServer 一致）。
func runServer(address string, deps *ihttp.Deps) {
	router := ihttp.NewRouter(deps)
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Minute
	s.WriteTimeout = 10 * time.Minute
	s.MaxHeaderBytes = 1 << 20

	deps.Log.Info("server run success on ", zap.String("address", address))
	deps.Log.Error(s.ListenAndServe().Error())
}
