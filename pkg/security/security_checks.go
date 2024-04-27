// pkg/security/security_checks.go

package security

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// BasicSecurityCheck performs several security-related checks on the database
func BasicSecurityCheck(db *sql.DB) error {
	color.Cyan("\nRunning security checks...")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Check", "Result", "Details", "Status", "Recommendation"})

	// Perform and display results for each check
	if err := checkEmptyPasswords(db, table); err != nil {
		return err
	}
	if err := checkAdminPrivileges(db, table); err != nil {
		return err
	}
	if err := checkMySQLVersion(db, table); err != nil {
		return err
	}

	table.Render() // Send output to stdout
	return nil
}

func checkEmptyPasswords(db *sql.DB, table *tablewriter.Table) error {
	query := "SELECT user, host FROM mysql.user WHERE authentication_string = '' OR authentication_string IS NULL;"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying users with empty passwords: %w", err)
	}
	defer rows.Close()

	var usersFound bool
	for rows.Next() {
		var user, host string
		if err := rows.Scan(&user, &host); err != nil {
			return fmt.Errorf("error scanning user data: %w", err)
		}
		table.Append([]string{"Empty Passwords", "Found", fmt.Sprintf("%s@%s", user, host), "Bad", "Set strong passwords for all accounts."})
		usersFound = true
	}

	if !usersFound {
		table.Append([]string{"Empty Passwords", "None Found", "No users with empty passwords.", "Good", "No action needed."})
	}
	return nil
}

func checkAdminPrivileges(db *sql.DB, table *tablewriter.Table) error {
	query := "SELECT user, host FROM mysql.user WHERE Super_priv = 'Y';"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying users with administrative privileges: %w", err)
	}
	defer rows.Close()

	var adminsFound bool
	for rows.Next() {
		var user, host string
		if err := rows.Scan(&user, &host); err != nil {
			return fmt.Errorf("error scanning admin user data: %w", err)
		}
		table.Append([]string{"Admin Privileges", "Found", fmt.Sprintf("%s@%s", user, host), "Bad", "Minimize the number of admin accounts."})
		adminsFound = true
	}

	if !adminsFound {
		table.Append([]string{"Admin Privileges", "None Found", "No administrator accounts found.", "Good", "No action needed."})
	}
	return nil
}

func checkMySQLVersion(db *sql.DB, table *tablewriter.Table) error {
	var version string
	if err := db.QueryRow("SELECT @@version;").Scan(&version); err != nil {
		return fmt.Errorf("error retrieving MySQL version: %w", err)
	}
	// Example version compliance logic
	if version < "8.0" { // This is a basic check, a real check should parse version properly
		table.Append([]string{"MySQL Version", "Outdated", fmt.Sprintf("Version: %s", version), "Bad", "Upgrade to MySQL 8.0 or higher."})
	} else {
		table.Append([]string{"MySQL Version", "Up-to-date", fmt.Sprintf("Version: %s", version), "Good", "No action needed."})
	}
	return nil
}
