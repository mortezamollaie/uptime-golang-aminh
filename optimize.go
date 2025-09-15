package main

import (
	"database/sql"
	"fmt"
	"log"
	"uptime/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load config
	config.Load()

	// Connect to MySQL
	db, err := sql.Open("mysql", config.AppConfig.Database.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// List of SQL statements
	sqlStatements := []string{
		"DROP INDEX idx_node_logs_node_id_created_at ON node_logs;",
		"CREATE INDEX idx_node_logs_node_id_created_at ON node_logs(node_id, created_at DESC);",

		"DROP INDEX idx_node_logs_node_id_id ON node_logs;",
		"CREATE INDEX idx_node_logs_node_id_id ON node_logs(node_id, id DESC);",

		"DROP INDEX idx_node_logs_status ON node_logs;",
		"CREATE INDEX idx_node_logs_status ON node_logs(status);",

		"DROP INDEX idx_node_logs_node_id_status ON node_logs;",
		"CREATE INDEX idx_node_logs_node_id_status ON node_logs(node_id, status);",

		"DROP INDEX idx_node_logs_up ON node_logs;",
		"CREATE INDEX idx_node_logs_up ON node_logs(up);",

		"DROP INDEX idx_node_logs_composite ON node_logs;",
		"CREATE INDEX idx_node_logs_composite ON node_logs(node_id, created_at DESC, status);",

		"DROP INDEX idx_node_logs_created_at ON node_logs;",
		"CREATE INDEX idx_node_logs_created_at ON node_logs(created_at DESC);",

		"DROP INDEX idx_node_logs_uptime ON node_logs;",
		"CREATE INDEX idx_node_logs_uptime ON node_logs(node_id, up, created_at DESC);",

		"ANALYZE TABLE nodes;",
		"ANALYZE TABLE node_logs;",
	}

	// Execute each statement
	for _, stmt := range sqlStatements {
		if _, err := db.Exec(stmt); err != nil {
			fmt.Printf("Error executing statement: %s\nError: %v\n", stmt, err)
		} else {
			fmt.Printf("Executed: %s\n", stmt)
		}
	}

	fmt.Println("Optimization SQL executed successfully.")
}
