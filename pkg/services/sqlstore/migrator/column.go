package migrator


type Column struct {
	Name            string
	Type            string
	Length          int
	Length2         int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
	Default         string
}

func (col *Column) String(d Dialect) string {
	return d.ColString(col)
}

func (col *Column) StringNoPk(d Dialect) string {
	return d.ColStringNoPk(col)
}

