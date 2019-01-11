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

func LoadGqlConfigFile(filename string) {
    configFile, err := ReadYamlFile(filename)
    if nil != err {
        log.Println(err)
        return
    }

    loadConfigYAML(configFile)
}

func LoadGqlConfigString(yamlconf string) {
    configFile, err := ReadYamlString(yamlconf)
    if nil != err {
        log.Println(err)
        return
    }

    loadConfigYAML(configFile)
}

func loadConfigYAML(file *YamlFile) {
    configMap, err := MapOfYaml(file.Root, "root")
    if nil != err {
        log.Println(err)
        return
    }

    for name, node := range configMap {
        configItemMap, err := MapOfYaml(node, name)
        if nil != err {
            log.Println(err)
            continue
        }

        driverName, err := StringOfYaml(configItemMap["DriverName"], name+".DriverName")
        if nil != err {
            log.Println(err)
            continue
        }
        dataSourceName, err := StringOfYaml(configItemMap["DataSourceName"], name+".DataSourceName")
        if nil != err {
            log.Println(err)
            continue
        }

        db, err := sql.Open(driverName, dataSourceName)
        if nil != err {
            log.Println(err)
            continue
        }

        maxOpenConns, err := IntOfYaml(configItemMap["MaxOpenConns"], name+".MaxOpenConns")
        if nil == err {
            db.SetMaxOpenConns(int(maxOpenConns))
        }
        maxIdleConns, err := IntOfYaml(configItemMap["MaxIdleConns"], name+".MaxIdleConns")
        if nil == err {
            db.SetMaxIdleConns(int(maxIdleConns))
        }
        connMaxLifetime, err := IntOfYaml(configItemMap["ConnMaxLifetime"], name+".ConnMaxLifetime")
        if nil == err {
            db.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))
        }

        connMap[name] = db
    }
}
