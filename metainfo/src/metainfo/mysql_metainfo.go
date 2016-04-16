package metainfo

import "database/sql"

import _ "github.com/go-sql-driver/mysql"

type MysqlMetaInfo struct {
	conn *sql.DB

	estimTablesNo int
}

func NewMysqlMetaInfo(conn *sql.DB, estim int) *MysqlMetaInfo {
	return &MysqlMetaInfo{
		conn:          conn,
		estimTablesNo: estim,
	}
}

func (m *MysqlMetaInfo) GetTables() (buf []string, err error) {
	if m.estimTablesNo > 0 {
		buf = make([]string, 0, m.estimTablesNo+5)
	}

	rows, err := m.conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		buf = append(buf, name)
	}

	return buf, nil
}

func sanitizeTableName(name string) string {
	properCode := func(b byte) bool {
		return ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == '_'
	}
	if len(name) == 0 {
		return ""
	}

	if !(properCode(name[0]) && (name[0] < '0' || name[0] > '9')) {
		return ""
	}
	for i := 1; i < len(name); i++ {
		if !properCode(name[i]) {
			return name[:i]
		}
	}
	return name
}

func (m *MysqlMetaInfo) GetTableInfo(name string) (*TableInfo, error) {
	rows, err := m.conn.Query("DESC " + sanitizeTableName(name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fieldName string
	var fieldType string
	var fieldIsNullable string
	var fieldDefault string
	var t1 string
	var t2 string
	res := &TableInfo{
		Name:   name,
		Fields: make([]FieldInfo, 0, 128),
	}
	fields := res.Fields
	for rows.Next() {
		err = rows.Scan(
			&fieldName, &fieldType, &fieldIsNullable, &t1, &fieldDefault, &t2)
		if err != nil {
			return nil, err
		}
		fields = append(fields, FieldInfo{
			Field:    fieldName,
			Type:     fieldType,
			Nullable: fieldIsNullable != "NO",
			Default:  fieldDefault,
		})
	}
	res.Fields = fields
	return res, nil
}
