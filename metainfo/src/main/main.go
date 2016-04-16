package main

import (
	"database/sql"
	"log"
	"metainfo"
)

func main() {
	conn, err := sql.Open("mysql", "root@/stat?interpolateParams=true")
	if err != nil {
		log.Fatal(err)
	}

	meta := metainfo.NewMysqlMetaInfo(conn, 1024)
	tables, err := meta.GetTables()
	if err != nil {
		log.Fatal(err)
	}

	tableName := tables[0]
	info, err := meta.GetTableInfo(tableName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*info)
}
