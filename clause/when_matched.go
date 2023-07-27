package clause

import (
	"gorm.io/gorm/clause"
)

type WhenMatched struct {
	Set clause.Set
}

func (w WhenMatched) Name() string {
	return "WHEN MATCHED"
}

func (w WhenMatched) Build(builder clause.Builder) {
	if len(w.Set) > 0 {
		_, _ = builder.WriteString("THEN UPDATE ")
		_, _ = builder.WriteString(w.Set.Name())
		_ = builder.WriteByte(' ')
		w.Set.Build(builder)
	}
}

func (w WhenMatched) MergeClause(clause *clause.Clause) {
	clause.Name = w.Name()
	clause.Expression = w
}
