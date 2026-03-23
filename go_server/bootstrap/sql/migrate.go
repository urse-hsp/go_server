// migration 执行器
// 每次启动都会执行所有 SQL
package bootstrap

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func RunMigrations() {
	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println("执行 migration:", file)

		sqlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		err = DB.Exec(string(sqlBytes)).Error
		if err != nil {
			panic("执行失败: " + file + " | " + err.Error())
		}
	}
}
