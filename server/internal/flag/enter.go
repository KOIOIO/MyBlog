// Package flag 提供 CLI 管理命令（建表/导出导入/建 ES 索引/创建管理员）。
// 依赖由 bootstrap 构造注入，不读取 global。
package flag

import (
	"errors"
	"fmt"
	"os"

	"server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Deps 命令依赖。
type Deps struct {
	DB  *gorm.DB
	ES  *elasticsearch.TypedClient
	Cfg *config.Config
	Log *zap.Logger
}

// 定义 CLI 标志，用于不同操作的命令行选项
var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initializes the structure of the MySQL database table.",
	}
	sqlExportFlag = &cli.BoolFlag{
		Name:  "sql-export",
		Usage: "Exports SQL data to a specified file.",
	}
	sqlImportFlag = &cli.StringFlag{
		Name:  "sql-import",
		Usage: "Imports SQL data from a specified file.",
	}
	esFlag = &cli.BoolFlag{
		Name:  "es",
		Usage: "Initializes the Elasticsearch index.",
	}
	esExportFlag = &cli.BoolFlag{
		Name:  "es-export",
		Usage: "Exports data from Elasticsearch to a specified file.",
	}
	esImportFlag = &cli.StringFlag{
		Name:  "es-import",
		Usage: "Imports data into Elasticsearch from a specified file.",
	}
	adminFlag = &cli.BoolFlag{
		Name:  "admin",
		Usage: "Creates an administrator using the name, email and address specified in the config.yaml file.",
	}
)

// Run 执行基于命令行标志的相应操作。
func Run(c *cli.Context, deps Deps) {
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Only one command can be specified", 1)
		deps.Log.Error("Invalid command usage:", zap.Error(err))
		os.Exit(1)
	}

	switch {
	case c.Bool(sqlFlag.Name):
		if err := SQL(deps); err != nil {
			deps.Log.Error("Failed to create table structure:", zap.Error(err))
		} else {
			deps.Log.Info("Successfully created table structure")
		}
	case c.Bool(sqlExportFlag.Name):
		if err := SQLExport(deps); err != nil {
			deps.Log.Error("Failed to export SQL data:", zap.Error(err))
		} else {
			deps.Log.Info("Successfully exported SQL data")
		}
	case c.IsSet(sqlImportFlag.Name):
		if errs := SQLImport(deps, c.String(sqlImportFlag.Name)); len(errs) > 0 {
			var combinedErrors string
			for _, err := range errs {
				combinedErrors += err.Error() + "\n"
			}
			deps.Log.Error("Failed to import SQL data:", zap.Error(errors.New(combinedErrors)))
		} else {
			deps.Log.Info("Successfully imported SQL data")
		}
	case c.Bool(esFlag.Name):
		if err := Elasticsearch(deps); err != nil {
			deps.Log.Error("Failed to create ES indices:", zap.Error(err))
		} else {
			deps.Log.Info("Successfully created ES indices")
		}
	case c.Bool(esExportFlag.Name):
		if err := ElasticsearchExport(deps); err != nil {
			deps.Log.Error("Failed to export ES data:", zap.Error(err))
		} else {
			deps.Log.Info("Successfully exported ES data")
		}
	case c.IsSet(esImportFlag.Name):
		if num, err := ElasticsearchImport(deps, c.String(esImportFlag.Name)); err != nil {
			deps.Log.Error("Failed to import ES data:", zap.Error(err))
		} else {
			deps.Log.Info(fmt.Sprintf("Successfully imported ES data, totaling %d records", num))
		}
	case c.Bool(adminFlag.Name):
		if err := Admin(deps); err != nil {
			deps.Log.Error("Failed to create an administrator:", zap.Error(err))
		} else {
			deps.Log.Info("Successfully created an administrator")
		}
	default:
		err := cli.NewExitError("unknown command", 1)
		deps.Log.Error(err.Error(), zap.Error(err))
	}
}

// NewApp 创建并配置一个新的 CLI 应用程序。
func NewApp(deps Deps) *cli.App {
	app := cli.NewApp()
	app.Name = "Go Blog"
	app.Flags = []cli.Flag{
		sqlFlag,
		sqlExportFlag,
		sqlImportFlag,
		esFlag,
		esExportFlag,
		esImportFlag,
		adminFlag,
	}
	app.Action = func(c *cli.Context) error {
		Run(c, deps)
		return nil
	}
	return app
}

// InitFlag 若有命令行参数则执行 CLI 操作并退出。
func InitFlag(deps Deps) {
	if len(os.Args) > 1 {
		app := NewApp(deps)
		err := app.Run(os.Args)
		if err != nil {
			deps.Log.Error("Application execution encountered an error:", zap.Error(err))
			os.Exit(1)
		}
		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			fmt.Println("Displaying help message...")
		}
		os.Exit(0)
	}
}
