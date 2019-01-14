package gokits

import "fmt"

type UnknownGqlConnectionName struct {
    Name string
}

func (e *UnknownGqlConnectionName) Error() string {
    return fmt.Sprintf("gql: Unknown connection named: %s", e.Name)
}
