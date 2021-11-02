package config

import (
	"fmt"
)

const host = "localhost"
const port = "5432"
const user = "postgres"
const dbname = "hospital"

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbname)
}
