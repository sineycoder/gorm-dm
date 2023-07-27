package clause

import (
	"gorm.io/gorm/clause"
)

type Merge struct {
	From  clause.From
	Using clause.Values
	On    []clause.Eq // resolve conflict, e.g. a=1 or b=2
}

func (merge Merge) Name() string {
	return "MERGE"
}

func DefaultExcludeName() string {
	return "excluded"
}

// Build for expression
func (merge Merge) Build(builder clause.Builder) {
	cols, vals := merge.Using.Columns, merge.Using.Values
	clause.Insert{}.Build(builder)
	sel := clause.Select{}
	_, _ = builder.WriteString(" USING (")
	for idx, v := range vals {
		if idx > 0 {
			_, _ = builder.WriteString(" UNION ")
		}
		_, _ = builder.WriteString(sel.Name())
		_ = builder.WriteByte(' ')
		for i, col := range cols {
			if i > 0 {
				_ = builder.WriteByte(',')
			}
			builder.AddVar(builder, v[i])
			col.Alias = col.Name
			col.Name = ""
			builder.WriteQuoted(col)
		}
		_ = builder.WriteByte(' ')
		_, _ = builder.WriteString(merge.From.Name())
		_ = builder.WriteByte(' ')
		merge.From.Build(builder)
	}
	_, _ = builder.WriteString(") ")
	_, _ = builder.WriteString(DefaultExcludeName())
	_, _ = builder.WriteString(" ON (")
	for idx, on := range merge.On {
		if idx > 0 {
			_, _ = builder.WriteString(clause.OrWithSpace)
		}
		on.Build(builder)
	}
	_ = builder.WriteByte(')')
}

func (merge Merge) MergeClause(clause *clause.Clause) {
	clause.Name = merge.Name()
	clause.Expression = merge
}
