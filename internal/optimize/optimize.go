package optimize

import (
	"database/sql"
	"fmt"
)

func Run(sqlDB *sql.DB) {
	// تغییر ساختار ستون id
	_, err := sqlDB.Exec(`
		ALTER TABLE node_logs 
		MODIFY COLUMN id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY
	`)
	if err != nil {
		fmt.Printf("Structure error: %v\n", err)
	}

	// لیست ایندکس‌ها برای حذف و ایجاد مجدد
	indexes := []struct {
		name string
		sql  string
	}{
		{"idx_node_logs_node_id_created_at", "CREATE INDEX idx_node_logs_node_id_created_at ON node_logs(node_id, created_at DESC)"},
		{"idx_node_logs_node_id_id", "CREATE INDEX idx_node_logs_node_id_id ON node_logs(node_id, id DESC)"},
		{"idx_node_logs_status", "CREATE INDEX idx_node_logs_status ON node_logs(status)"},
		{"idx_node_logs_node_id_status", "CREATE INDEX idx_node_logs_node_id_status ON node_logs(node_id, status)"},
		{"idx_node_logs_up", "CREATE INDEX idx_node_logs_up ON node_logs(up)"},
		{"idx_node_logs_composite", "CREATE INDEX idx_node_logs_composite ON node_logs(node_id, created_at DESC, status)"},
		{"idx_node_logs_created_at", "CREATE INDEX idx_node_logs_created_at ON node_logs(created_at DESC)"},
		{"idx_node_logs_uptime", "CREATE INDEX idx_node_logs_uptime ON node_logs(node_id, up, created_at DESC)"},
	}

	successCount := 0
	for _, index := range indexes {
		// اگر ایندکس وجود دارد Drop کن
		var exists int
		err = sqlDB.QueryRow(`
			SELECT COUNT(*) 
			FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			  AND table_name = 'node_logs' 
			  AND index_name = ?`, index.name).Scan(&exists)
		if err == nil && exists > 0 {
			_, dropErr := sqlDB.Exec(fmt.Sprintf("DROP INDEX %s ON node_logs", index.name))
			if dropErr != nil {
				fmt.Printf("Failed to drop index %s: %v\n", index.name, dropErr)
			}
		}

		// ایجاد مجدد ایندکس
		_, createErr := sqlDB.Exec(index.sql)
		if createErr != nil {
			fmt.Printf("Index %s creation failed: %v\n", index.name, createErr)
		} else {
			successCount++
		}
	}

	// آنالیز جداول
	_, err = sqlDB.Exec("ANALYZE TABLE nodes")
	if err != nil {
		return
	}
	_, err = sqlDB.Exec("ANALYZE TABLE node_logs")
	if err != nil {
		return
	}

	fmt.Printf("Optimization complete: %d/%d indexes recreated\n", successCount, len(indexes))
}
