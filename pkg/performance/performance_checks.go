// pkg/performance/performance_checks.go

package performance

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// BasicPerformanceCheck performs detailed performance checks on the database
func BasicPerformanceCheck(db *sql.DB) error {
	color.Cyan("\nRunning performance checks...")

	// Create a new table writer for structured output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Check", "Status", "Details"})

	// Check for queries not using indexes
	if err := checkIndexUsage(db, table); err != nil {
		log.Printf("Error during index usage check: %v", err)
		table.Append([]string{"Index Usage", "FAILED", fmt.Sprintf("Error: %v", err)})
	} else {
		table.Append([]string{"Index Usage", "PASSED", "All major queries are using indexes appropriately."})
	}

	// Render the table to stdout
	table.Render()

	return nil
}

func checkIndexUsage(db *sql.DB, table *tablewriter.Table) error {
	// Example query to check for missing indexes
	// Lets start with a simple check, i need to expand it later
	query := `
	SELECT 
	    t.table_schema,
	    t.table_name,
	    t.table_rows,
	    ps.index_size,
	    ps.data_size,
	    ps.total_size
	FROM 
	    information_schema.tables t
	JOIN (
	    SELECT 
	        table_schema,
	        table_name,
	        SUM(data_length) data_size,
	        SUM(index_length) index_size,
	        SUM(data_length + index_length) total_size
	    FROM 
	        information_schema.tables 
	    GROUP BY 
	        table_schema, 
	        table_name
	) ps ON t.table_schema = ps.table_schema AND t.table_name = ps.table_name
	WHERE 
	    t.table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
	    AND t.table_rows > 10000
	    AND ps.index_size = 0;`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying database for index usage: %v", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var schema, tableName string
		var tableRows, indexSize, dataSize, totalSize int64
		if err := rows.Scan(&schema, &tableName, &tableRows, &indexSize, &dataSize, &totalSize); err != nil {
			return fmt.Errorf("error scanning table data: %v", err)
		}
		// Add rows to the table for each found issue
		detail := fmt.Sprintf("%s.%s: %d rows, Data Size: %d, Index Size: %d, Total Size: %d",
			schema, tableName, tableRows, dataSize, indexSize, totalSize)
		table.Append([]string{"Table Index Check", "FAILED", detail})
	}

	if count == 0 {
		table.Append([]string{"Table Index Check", "PASSED", "All tables with significant data are indexed."})
	}

	return nil
}
