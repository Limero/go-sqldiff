package sqldiff

import (
	"fmt"
	"os/exec"
	"strings"
)

// Returns SQL queries that would transform source database to target database
func Compare(sourceDatabasePath string, targetDatabasePath string) ([]string, error) {
	out, err := exec.Command("sqldiff", sourceDatabasePath, targetDatabasePath).Output()
	if err != nil {
		return nil, fmt.Errorf("sqldiff failed: %w", err)
	}

	return strings.Split(string(out[:]), "\n"), nil
}
