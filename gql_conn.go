package gokits

import (
    "database/sql"
    "log"
    "time"
)

type GqlConfig struct {
    DriverName      string
    DataSourceName  string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime int // Seconds
}

var connMap = make(map[string]*sql.DB)

func GqlConnection(name string) *sql.DB {
    return connMap[name]
}

//noinspection GoUnusedExportedFunction
func LoadGqlConfigFile(filename string) {
    configFile, err := ReadYamlFile(filename)
    if nil != err {
        log.Println(err)
        return
    }

    loadConfigYAML(configFile)
}

//noinspection GoUnusedExportedFunction
func LoadGqlConfigString(yamlconf string) {
    configFile, err := ReadYamlString(yamlconf)
    if nil != err {
        log.Println(err)
        return
    }

    loadConfigYAML(configFile)
}

func loadConfigYAML(file *YamlFile) {
    configMap, err := file.RootMap()
    if nil != err {
        log.Println(err)
        return
    }

    for name := range configMap {
        driverName, err := file.GetString(name+".DriverName")
        if nil != err {
            log.Println(err)
            continue
        }
        dataSourceName, err := file.GetString(name+".DataSourceName")
        if nil != err {
            log.Println(err)
            continue
        }

        db, err := sql.Open(driverName, dataSourceName)
        if nil != err {
            log.Println(err)
            continue
        }

        maxOpenConns, err := file.GetInt(name+".MaxOpenConns")
        if nil == err {
            db.SetMaxOpenConns(int(maxOpenConns))
        }
        maxIdleConns, err := file.GetInt(name+".MaxIdleConns")
        if nil == err {
            db.SetMaxIdleConns(int(maxIdleConns))
        }
        connMaxLifetime, err := file.GetInt(name+".ConnMaxLifetime")
        if nil == err {
            db.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))
        }

        connMap[name] = db
    }
}
