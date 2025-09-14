package optimize

import (
	"fmt"
	"database/sql"
)

func Run(sqlDB *sql.DB) {
	_, err := sqlDB.Exec("ALTER TABLE node_logs MODIFY COLUMN id bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY")
	if err != nil {
		fmt.Printf("Structure error: %v\n", err)
	}

	indexes := []struct{
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
		var exists int
		err = sqlDB.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = 'ms-uptime' AND table_name = 'node_logs' AND index_name = ?", index.name).Scan(&exists)
		if err == nil && exists > 0 {
			successCount++
			continue
		}
		_, err = sqlDB.Exec(index.sql)
		if err != nil {
			fmt.Printf("Index %s failed: %v\n", index.name, err)
		} else {
			successCount++
		}
	}

	sqlDB.Exec("ANALYZE TABLE nodes, node_logs")
	fmt.Printf("Optimization complete: %d/%d indexes\n", successCount, len(indexes))
}
