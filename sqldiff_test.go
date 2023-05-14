package sqldiff

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFakeDatabase(t *testing.T, path string, queries ...string) string {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	for _, query := range queries {
		_, err := db.Exec(query)
		require.NoError(t, err)
	}

	return path
}

func TestCompare(t *testing.T) {
	schema := `
		CREATE TABLE tablename (
			column1 VARCHAR(1024) NOT NULL,
			column2 INTEGER(1) NOT NULL
		)
	`

	db1 := createFakeDatabase(
		t,
		"/tmp/source.db",
		schema,
		"INSERT INTO tablename (column1, column2) VALUES ('abc', 1)",
		"INSERT INTO tablename (column1, column2) VALUES ('def', 0)",
	)
	defer os.Remove(db1)
	db2 := createFakeDatabase(
		t,
		"/tmp/target.db",
		schema,
		"INSERT INTO tablename (column1, column2) VALUES ('abc', 0)",
		"INSERT INTO tablename (column1, column2) VALUES ('def', 1)",
	)
	defer os.Remove(db2)

	out, err := Compare(db1, db2)
	require.NoError(t, err)

	assert.Equal(t, []string{
		"UPDATE tablename SET column2=0 WHERE rowid=1;",
		"UPDATE tablename SET column2=1 WHERE rowid=2;",
		"",
	}, out)
}
