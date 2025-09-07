package main

import (
	"fmt"
	"log"
	"uptime/config"
	"uptime/database"
)

func main() {
	fmt.Println("🚀 Database Index Optimizer")
	fmt.Println("============================")
	
	// Initialize config and database
	config.Load()
	database.Connect()
	
	// Get raw SQL connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database connection:", err)
	}
	
	fmt.Println("✅ Connected to database")
	
	// Fix AUTO_INCREMENT if needed
	fmt.Println("\n🔧 Checking table structure...")
	_, err = sqlDB.Exec("ALTER TABLE node_logs MODIFY COLUMN id bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY")
	if err != nil {
		fmt.Printf("⚠️  Structure already correct or error: %v\n", err)
	} else {
		fmt.Println("✅ Table structure fixed")
	}
	
	// Create indexes
	fmt.Println("\n📚 Creating performance indexes...")
	
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
	for i, index := range indexes {
		fmt.Printf("Creating index %d/%d (%s)... ", i+1, len(indexes), index.name)
		
		// Check if index exists
		var exists int
		err = sqlDB.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = 'ms-uptime' AND table_name = 'node_logs' AND index_name = ?", index.name).Scan(&exists)
		
		if err == nil && exists > 0 {
			fmt.Printf("⏭️  Already exists\n")
			successCount++
			continue
		}
		
		_, err = sqlDB.Exec(index.sql)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
			successCount++
		}
	}
	
	// Analyze tables
	fmt.Println("\n🔍 Analyzing tables...")
	_, err = sqlDB.Exec("ANALYZE TABLE nodes, node_logs")
	if err != nil {
		fmt.Printf("⚠️  Analyze failed: %v\n", err)
	} else {
		fmt.Println("✅ Tables analyzed")
	}
	
	// Show results
	fmt.Printf("\n🎉 Optimization completed!\n")
	fmt.Printf("📊 %d/%d indexes created successfully\n", successCount, len(indexes))
	fmt.Println("🚀 Database is now optimized for high performance!")
}
