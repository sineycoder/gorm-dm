package clause

import (
	"gorm.io/gorm/clause"
)

type WhenNotMatched struct {
	Values clause.Values
}

func (w WhenNotMatched) Name() string {
	return "WHEN NOT MATCHED"
}

func (w WhenNotMatched) Build(builder clause.Builder) {
	if len(w.Values.Columns) > 0 {
		_, _ = builder.WriteString("THEN")
		_, _ = builder.WriteString(" INSERT ")
		_ = builder.WriteByte('(')
		for idx, column := range w.Values.Columns {
			if idx > 0 {
				_ = builder.WriteByte(',')
			}
			builder.WriteQuoted(column)
		}
		_ = builder.WriteByte(')')
		_, _ = builder.WriteString(" VALUES ")
		_ = builder.WriteByte('(')
		for idx, column := range w.Values.Columns {
			if idx > 0 {
				_ = builder.WriteByte(',')
			}
			column.Table = DefaultExcludeName()
			builder.WriteQuoted(column)
		}
		_ = builder.WriteByte(')')
	}
}

func (w WhenNotMatched) MergeClause(clause *clause.Clause) {
	clause.Name = w.Name()
	clause.Expression = w
}
