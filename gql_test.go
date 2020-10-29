package gokits

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func loadConfig() {
    dsn := os.Getenv("MYSQL_DATA_SOURCE_NAME")
    LoadGqlConfigString(`
Default:
  DriverName:       mysql
  DataSourceName:   ` + dsn)
}

func TestNewGql(t *testing.T) {
    a := assert.New(t)
    loadConfig()

    gql, _ := DefaultGql()

    result, _ := gql.New().Sql(`select 'x' as "X"`).Query()
    a.Equal("x", result[0]["X"])

    result, _ = gql.New().Sql("select 'xx' where 'x' = ?").Params("x").Query()
    a.Equal("xx", result[0]["xx"])
}

func TestDml(t *testing.T) {
    a := assert.New(t)
    loadConfig()

    gql, _ := DefaultGql()

    _, _ = gql.New().Sql(`
    create table app_config (
        config_name varchar(100) not null,
        config_value text not null,
        primary key (config_name)
    );
    `).Execute()

    count, _ := gql.New().Sql(`
    insert into app_config
          (config_name  ,config_value)
    values(?            ,?)`).Params("TEST", 123).Execute()
    a.EqualValues(1, count)

    result, _ := gql.New().Sql(`
select config_name
      ,config_value
  from app_config
 where config_name = ?`).Params("TEST").Query()
    a.Equal("123", result[0]["config_value"])

    count, _ = gql.New().Sql(`
    update app_config
       set config_value = ?
     where config_name = ?`).Params(345, "TEST").Execute()
    a.EqualValues(1, count)

    result, _ = gql.New().Sql(`
select config_name
      ,config_value
  from app_config
 where config_name = ?`).Params("TEST").Query()
    a.Equal("345", result[0]["config_value"])

    count, _ = gql.New().Sql(`
    delete from app_config
    where config_name = ?`).Params("TEST").Execute()
    a.EqualValues(1, count)
}
