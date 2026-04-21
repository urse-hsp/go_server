package bootstrap

import (
	"go-server/internal/model"
	"go-server/pkg/log"
	"os"
	"path/filepath"
	"sort"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrateServer(db *gorm.DB, log *log.Logger) *MigrateServer {
	return &MigrateServer{
		db:  db,
		log: log,
	}
}

// Start 入口
func (m *MigrateServer) Start() error {
	// 执行 SQL migration
	if err := m.RunSQLMigrations(); err != nil {
		return err
	}

	// 手动结构变更
	if err := m.migrateUser(); err != nil {
		return err
	}

	// 同步结构
	if err := m.db.AutoMigrate(model.GetModels()...); err != nil {
		return err
	}

	m.log.Info("migrate success")

	return nil
}

func (m *MigrateServer) Stop() error {
	m.log.Info("migration stop")
	return nil
}

// RunSQLMigrations SQL 迁移核心
func (m *MigrateServer) RunSQLMigrations() error {

	// 记录表 migration sql 执行历史，防止重复执行
	// CREATE TABLE schema_migrations (
	//     id INT AUTO_INCREMENT PRIMARY KEY,
	//     filename VARCHAR(255) UNIQUE,
	//     executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	// );

	// 1. 创建版本表
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return err
	}

	// 2. 已执行版本
	var applied []string
	if err := m.db.Raw("SELECT version FROM schema_migrations").Scan(&applied).Error; err != nil {
		return err
	}

	appliedMap := make(map[string]bool)
	for _, v := range applied {
		appliedMap[v] = true
	}

	// 3. 读取 migration 文件（推荐绝对路径方式）
	// dir := filepath.Join("cmd", "migration", "sql")
	dir, _ := filepath.Abs("cmd/migration/sql")
	if _, err := os.Stat(dir); err != nil {
		m.log.Fatal("migration dir not found", zap.String("dir", dir), zap.Error(err))
	}

	m.log.Info("migration dir", zap.String("dir", dir))

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// 4. 排序保证执行顺序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// 5. 执行 migration
	for _, f := range files {
		name := f.Name()

		// 只处理 sql 文件
		if f.IsDir() || filepath.Ext(name) != ".sql" {
			continue
		}

		// 已执行跳过
		if appliedMap[name] {
			m.log.Info("skip migration", zap.String("file", name))
			continue
		}

		fullPath := filepath.Join(dir, name)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		m.log.Info("running migration", zap.String("file", name))

		tx := m.db.Begin()

		// 执行 SQL
		if err := tx.Exec(string(content)).Error; err != nil {
			tx.Rollback()

			m.log.Error("migration failed",
				zap.String("file", name),
				zap.Error(err),
			)

			return err
		}

		// 写入版本记录
		if err := tx.Exec(
			"INSERT INTO schema_migrations (version) VALUES (?)",
			name,
		).Error; err != nil {
			tx.Rollback()

			m.log.Error("version insert failed",
				zap.String("file", name),
				zap.Error(err),
			)

			return err
		}

		if err := tx.Commit().Error; err != nil {
			return err
		}

		m.log.Info("migration success", zap.String("file", name))
	}

	m.log.Info("SQL migrations done")

	return nil
}

// Migrator 手动结构变更
func (m *MigrateServer) migrateUser() error {
	// if m.db.Migrator().HasColumn(&model.User{}, "name") &&
	// 	!m.db.Migrator().HasColumn(&model.User{}, "username") {

	// 	return m.db.Migrator().RenameColumn(&model.User{}, "name", "username")
	// }
	return nil
}
