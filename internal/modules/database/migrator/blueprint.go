package migrator

type Column struct {
	Name       string
	Type       string
	Length     int
	Precision  int
	Scale      int
	IsNullable bool
	IsUnique   bool
	IsPrimary  bool
	IsAutoInc  bool
	IsUnsigned bool
	DefaultVal any
}

type ForeignColumn struct {
	*Column
	ForeignReference *ForeignReference
	OnDeleteAction   string
}

type ForeignReference struct {
	Table  string
	Column string
}

type Index struct {
	Name     string
	Columns  []string
	IsUnique bool
}

type CommandType int

const (
	AddColumnCommand CommandType = iota
	DropColumnCommand
	RenameColumnCommand
	AddIndexCommand
	DropIndexCommand
)

type Command struct {
	Type     CommandType
	Name     string
	Column   *Column
	From     string
	To       string
	Columns  []string
	IsUnique bool
}

type Blueprint struct {
	TableName      string
	Columns        []*Column
	ForeignColumns []*ForeignColumn
	Indexes        []*Index
	Commands       []*Command
	IsAltering     bool
}

func NewBlueprint(tableName string) *Blueprint {
	return &Blueprint{
		TableName:      tableName,
		Columns:        make([]*Column, 0),
		ForeignColumns: make([]*ForeignColumn, 0),
		Indexes:        make([]*Index, 0),
		Commands:       make([]*Command, 0),
		IsAltering:     false,
	}
}

func (bp *Blueprint) ID() *Column {
	col := &Column{
		Name:      "id",
		Type:      "INTEGER",
		IsPrimary: true,
	}

	bp.addColumn(col)
	return col
}

func (bp *Blueprint) String(name string, length int) *Column {
	col := &Column{
		Name:   name,
		Type:   "VARCHAR",
		Length: length,
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Text(name string) *Column {
	col := &Column{
		Name: name,
		Type: "TEXT",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Integer(name string) *Column {
	col := &Column{
		Name: name,
		Type: "INTEGER",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) BigInteger(name string) *Column {
	col := &Column{
		Name: name,
		Type: "BIGINT",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Boolean(name string) *Column {
	col := &Column{
		Name: name,
		Type: "BOOLEAN",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Timestamp(name string) *Column {
	col := &Column{
		Name: name,
		Type: "TIMESTAMP",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Decimal(name string, precision, scale int) *Column {
	col := &Column{
		Name:      name,
		Type:      "DECIMAL",
		Precision: precision,
		Scale:     scale,
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) Float(name string) *Column {
	col := &Column{
		Name: name,
		Type: "FLOAT",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) JSON(name string) *Column {
	col := &Column{
		Name: name,
		Type: "JSON",
	}
	bp.addColumn(col)
	return col
}

func (bp *Blueprint) addColumn(col *Column) {
	if bp.IsAltering {
		bp.Commands = append(bp.Commands, &Command{
			Type:   AddColumnCommand,
			Column: col,
		})
	} else {
		bp.Columns = append(bp.Columns, col)
	}
}

func (bp *Blueprint) Timestamps() {
	bp.Timestamp("created_at").Default("now()")
	bp.Timestamp("updated_at").Default("now()")
}

func (bp *Blueprint) SoftDeletes() {
	bp.Timestamp("deleted_at").Nullable()
}

func (bp *Blueprint) Foreign(name string) *ForeignColumn {
	col := &Column{
		Name: name,
		Type: "INTEGER",
	}

	fc := &ForeignColumn{
		Column:           col,
		ForeignReference: &ForeignReference{},
	}

	bp.ForeignColumns = append(bp.ForeignColumns, fc)
	return fc
}

func (bp *Blueprint) Index(columns []string, name string) {
	if name == "" {
		name = "idx_" + bp.TableName + "_" + columns[0]
	}

	idx := &Index{
		Name:     name,
		Columns:  columns,
		IsUnique: false,
	}

	if bp.IsAltering {
		bp.Commands = append(bp.Commands, &Command{
			Type:     AddIndexCommand,
			Name:     name,
			Columns:  columns,
			IsUnique: false,
		})
	} else {
		bp.Indexes = append(bp.Indexes, idx)
	}
}

func (bp *Blueprint) UniqueIndex(columns []string, name string) {
	if name == "" {
		name = "uniq_" + bp.TableName + "_" + columns[0]
	}

	idx := &Index{
		Name:     name,
		Columns:  columns,
		IsUnique: true,
	}

	if bp.IsAltering {
		bp.Commands = append(bp.Commands, &Command{
			Type:     AddIndexCommand,
			Name:     name,
			Columns:  columns,
			IsUnique: true,
		})
	} else {
		bp.Indexes = append(bp.Indexes, idx)
	}
}

func (bp *Blueprint) DropIndex(name string) {
	bp.Commands = append(bp.Commands, &Command{
		Type: DropIndexCommand,
		Name: name,
	})
}

func (bp *Blueprint) DropColumn(name string) {
	bp.Commands = append(bp.Commands, &Command{
		Type: DropColumnCommand,
		Name: name,
	})
}

func (bp *Blueprint) RenameColumn(from, to string) {
	bp.Commands = append(bp.Commands, &Command{
		Type: RenameColumnCommand,
		From: from,
		To:   to,
	})
}

func (c *Column) Nullable() *Column {
	c.IsNullable = true
	return c
}

func (c *Column) Default(val interface{}) *Column {
	c.DefaultVal = val
	return c
}

func (c *Column) Unique() *Column {
	c.IsUnique = true
	return c
}

func (c *Column) Unsigned() *Column {
	c.IsUnsigned = true
	return c
}

func (c *Column) AutoIncrement() *Column {
	c.IsAutoInc = true
	return c
}

func (c *Column) Primary() *Column {
	c.IsPrimary = true
	return c
}

func (fc *ForeignColumn) References(column string) *ForeignColumn {
	fc.ForeignReference.Column = column
	return fc
}

func (fc *ForeignColumn) On(table string) *ForeignColumn {
	fc.ForeignReference.Table = table
	return fc
}

func (fc *ForeignColumn) OnDelete(action string) *ForeignColumn {
	fc.OnDeleteAction = action
	return fc
}

func (fc *ForeignColumn) Nullable() *ForeignColumn {
	fc.Column.IsNullable = true
	return fc
}
