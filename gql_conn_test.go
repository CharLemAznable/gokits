package gokits

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "testing"
)

func TestConnection(t *testing.T) {
    dsn := os.Getenv("MYSQL_DATA_SOURCE_NAME")
    LoadGqlConfigString(`
Default:
  DriverName:       mysql
  DataSourceName:   ` + dsn + `
  MaxOpenConns:     50
  MaxIdleConns:     10
    `)

    connection := GqlConnection("Default")
    fmt.Println(connection)
}
