package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/meringu/terraform-private-registry/internal/cmd"
)

func main() {
	cmd.Execute()
}
