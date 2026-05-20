package grammar

import (
	"github.com/vayload/vayload/internal/modules/database/query"
)

// For compile select query, insert, update delete, upsert
type QueryGrammar interface {
	Compile(ast *query.Query) (string, []any)

	Wrap(value string) string
	Placeholder(position int) string
}
