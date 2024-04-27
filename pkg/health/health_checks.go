// pkg/health/health_checks.go

package health

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// CheckDatabaseUptime conducts several health checks on the database and outputs the results in a table format.
func CheckDatabaseUptime(db *sql.DB) error {
	color.Cyan("\nRunning health checks...")

	// Create a new table writer
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value", "Status", "Remarks"})

	// Uptime check
	uptime, err := getUptime(db)
	if err != nil {
		return fmt.Errorf("error retrieving uptime: %w", err)
	}
	table.Append([]string{"Uptime", fmt.Sprintf("%v", uptime), "OK", "Database running time"})

	// Active connections check
	activeConnections, err := getActiveConnections(db)
	if err != nil {
		return fmt.Errorf("error retrieving active connections: %w", err)
	}
	table.Append([]string{"Active Connections", fmt.Sprintf("%d", activeConnections), "OK", "Number of current connections"})

	// Slow queries check
	slowQueries, err := getSlowQueries(db)
	if err != nil {
		return fmt.Errorf("error retrieving slow queries: %w", err)
	}
	status := "OK"
	remarks := "Acceptable number of slow queries"
	if slowQueries > 100 { // Ask the user for a proper number, i have to find a better way, thinking mode...
		status = "Warning"
		remarks = "High number of slow queries"
	}
	table.Append([]string{"Slow Queries", fmt.Sprintf("%d", slowQueries), status, remarks})

	// Buffer pool usage check
	bufferUsage, err := getBufferPoolUsage(db)
	if err != nil {
		return fmt.Errorf("error retrieving buffer pool usage: %w", err)
	}

	// Declare status as new variable with an initial value
	var buff_status string = "OK"
	if bufferUsage == "N/A" {
		buff_status = "Warning"
	}

	table.Append([]string{"Buffer Pool Usage", bufferUsage, buff_status, "Percentage of buffer pool used"})

	// Open Tables check
	openTableStatus, err := getOpenTables(db)
	if err != nil {
		return fmt.Errorf("error performing open tables check: %w", err)
	}
	table.Append([]string{"Open Tables", "", openTableStatus, "Review table_open_cache if frequent opening/closing of tables."})

	// Render the table
	table.Render()

	return nil
}

func getUptime(db *sql.DB) (time.Duration, error) {
	var variableName string
	var uptimeSeconds int
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'Uptime'").Scan(&variableName, &uptimeSeconds)
	if err != nil {
		return 0, err
	}
	return time.Duration(uptimeSeconds) * time.Second, nil
}

func getActiveConnections(db *sql.DB) (int, error) {
	var variableName string
	var threadsConnected int
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'Threads_connected'").Scan(&variableName, &threadsConnected)
	if err != nil {
		return 0, err
	}
	return threadsConnected, nil
}

func getSlowQueries(db *sql.DB) (int, error) {
	var variableName string
	var slowQueries int
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'Slow_queries'").Scan(&variableName, &slowQueries)
	if err != nil {
		return 0, err
	}
	return slowQueries, nil
}

func getBufferPoolUsage(db *sql.DB) (string, error) {
	var bufferPoolTotal, bufferPoolUsed float64
	query := "SHOW GLOBAL STATUS WHERE Variable_name IN ('Innodb_buffer_pool_pages_total', 'Innodb_buffer_pool_pages_data');"
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error retrieving buffer pool status: %w", err)
	}
	defer rows.Close()

	var name string
	var value float64
	for rows.Next() {
		if err := rows.Scan(&name, &value); err != nil {
			return "", fmt.Errorf("error scanning buffer pool data: %w", err)
		}
		switch name {
		case "Innodb_buffer_pool_pages_total":
			bufferPoolTotal = value
		case "Innodb_buffer_pool_pages_data":
			bufferPoolUsed = value
		}
	}

	if bufferPoolTotal == 0 { // just to make sure, i dont fumble upon divison by zero..
		return "N/A", nil
	}

	usagePercentage := (bufferPoolUsed / bufferPoolTotal) * 100
	return fmt.Sprintf("%.2f%%", usagePercentage), nil
}
func getOpenTables(db *sql.DB) (string, error) {
	var openTables, tableOpenCache int
	var variableName string

	// Get the number of open tables
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'Open_tables';").Scan(&variableName, &openTables)
	if err != nil {
		return "", fmt.Errorf("error retrieving open tables: %w", err)
	}

	// Get the table_open_cache value
	err = db.QueryRow("SHOW VARIABLES LIKE 'table_open_cache';").Scan(&variableName, &tableOpenCache)
	if err != nil {
		return "", fmt.Errorf("error retrieving table_open_cache setting: %w", err)
	}

	var status string
	if openTables > int(0.8*float64(tableOpenCache)) { // 80% of table_open_cache as threshold
		status = fmt.Sprintf("Warning: High open tables count: %d of %d", openTables, tableOpenCache)
	} else {
		status = fmt.Sprintf("Normal: %d open tables of %d cache limit", openTables, tableOpenCache)
	}

	return status, nil
}
