package gokits

import (
    "database/sql"
    "log"
)

type Gql struct {
    conn *sql.DB
}

func NewGql(connName string) (*Gql, error) {
    connection := GqlConnection(connName)
    if nil == connection {
        return nil, &UnknownGqlConnectionName{Name: connName}
    }
    gql := new(Gql)
    gql.conn = connection
    return gql, nil
}

func DefaultGql() (*Gql, error) {
    return NewGql("Default")
}

type gqlInst struct {
    gql    *Gql
    sql    string
    params []interface{}
}

func (gql *Gql) New() *gqlInst {
    instance := new(gqlInst)
    instance.gql = gql
    return instance
}

func (instance *gqlInst) Sql(sql string) *gqlInst {
    instance.sql = sql
    return instance
}

func (instance *gqlInst) Params(params ... interface{}) *gqlInst {
    instance.params = params
    return instance
}

func (instance *gqlInst) Query() ([]map[string]string, error) {
    err := instance.gql.conn.Ping()
    if err != nil {
        log.Println(err)
        return nil, err
    }

    stmt, err := instance.gql.conn.Prepare(instance.sql)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    rows, err := stmt.Query(instance.params...)
    defer rows.Close()
    if err != nil {
        log.Println(err)
        return nil, err
    }

    columns, err := rows.Columns()
    if err != nil {
        log.Println(err)
        return nil, err
    }

    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    list := make([]map[string]string, 0)
    for rows.Next() {
        record := make(map[string]string)
        // 将行数据保存到record字典
        err = rows.Scan(scanArgs...)
        if err != nil {
            log.Println(err)
            return nil, err
        }
        for i, col := range values {
            if col != nil {
                record[columns[i]] = string(col)
            }
        }
        list = append(list, record)
    }
    return list, nil
}

func (instance *gqlInst) Execute() (int64, error) {
    err := instance.gql.conn.Ping()
    if err != nil {
        log.Println(err)
        return 0, err
    }

    stmt, err := instance.gql.conn.Prepare(instance.sql)
    if err != nil {
        log.Println(err)
        return 0, err
    }

    result, err := stmt.Exec(instance.params...)
    if err != nil {
        log.Println(err)
        return 0, err
    }

    return result.RowsAffected()
}
