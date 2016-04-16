package metainfo

type FieldInfo struct {
	Field    string
	Type     string
	Nullable bool
	Default  string
}

type TableInfo struct {
	Name   string
	Fields []FieldInfo
}

type MetaInfo interface {
	GetTables() ([]string, error)
	GetTableInfo(string) (*TableInfo, error)
}
