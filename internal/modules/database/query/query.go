package query

type StatementType uint8

const (
	StmtSelect StatementType = iota
	StmtInsert
	StmtUpdate
	StmtDelete
	StmtUpsert

	StmtCount
	StmtExists
)

type Query struct {
	Table    string
	SubQuery *Query
	Alias    string
	Columns  []string
	Distinct bool

	Joins   []Join
	Where   []Where
	OrWhere []Where

	GroupBy []string
	OrderBy []Order

	Having []Having

	Limit      *int64
	Offset     *int64
	CursorNext *string

	Args     []any
	StmtType StatementType

	InsertValues      map[string]any
	InsertMultiValues []map[string]any
	UpdateValues      map[string]any
	UpsertValues      map[string]any
	UpsertColumns     []string

	Unions       []Union
	LockMode     LockMode
	SeekColumn   string
	SeekOperator string
	SeekValue    any
}

type LockMode string

const (
	LockUpdate LockMode = "FOR UPDATE"
	LockShare  LockMode = "FOR SHARE"
	LockNone   LockMode = ""
)

type Union struct {
	Query *Query
	All   bool
}

type JoinTypes string

const (
	InnerJoinType = JoinTypes("INNER JOIN")
	LeftJoinType  = JoinTypes("LEFT JOIN")
	RightJoinType = JoinTypes("RIGHT JOIN")
	CrossJoinType = JoinTypes("CROSS JOIN")
)

type Join struct {
	Type     JoinTypes
	Table    string
	SubQuery *Query
	Alias    string
	On       string
	Args     []any
}

type WhereType string

const (
	WhereTypeBasic    WhereType = "basic"
	WhereTypeSubQuery WhereType = "subquery"
	WhereTypeIn       WhereType = "in"
	WhereTypeNotIn    WhereType = "notin"
	WhereTypeNull     WhereType = "null"
	WhereTypeNotNull  WhereType = "notnull"
	WhereTypeColumn   WhereType = "column"
	WhereTypeBetween  WhereType = "between"
)

type Where struct {
	Type     WhereType
	Column   string
	Operator string
	Value    any
	Value2   any // for between
	SubQuery *Query
	IsOr     bool
}

type Order struct {
	Column    string
	Direction string
}

type Having struct {
	Condition string
	Args      []any
}
