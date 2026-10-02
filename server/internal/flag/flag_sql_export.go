package flag

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// SQLExport 导出 MySQL 数据。
func SQLExport(deps Deps) error {
	mysql := deps.Cfg.Mysql

	timer := time.Now().Format("20060102")
	sqlPath := fmt.Sprintf("mysql_%s.sql", timer)
	cmd := exec.Command("docker", "exec", "mysql", "mysqldump", "-u"+mysql.Username, "-p"+mysql.Password, mysql.DBName)

	outFile, err := os.Create(sqlPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	return cmd.Run()
}

// SQLImport 导入 MySQL 数据。
func SQLImport(deps Deps, sqlPath string) (errs []error) {
	byteData, err := os.ReadFile(sqlPath)
	if err != nil {
		return append(errs, err)
	}
	sqlList := splitSQL(byteData)
	for _, sql := range sqlList {
		if err := deps.DB.Exec(sql).Error; err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return nil
}
