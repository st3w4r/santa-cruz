package dbsqlc

import (
	_ "embed"
)

//go:embed schema.sql
var CreateTableSchema string
