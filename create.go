package gorm_dm

import (
	"database/sql"
	"fmt"
	"reflect"

	clause_ "github.com/sineycoder/gorm-dm/clause"
	"github.com/sineycoder/gorm-dm/utils/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	schema2 "gorm.io/gorm/schema"
)

func Create(db *gorm.DB) {
	stmt, schema := db.Statement, db.Statement.Schema
	if stmt == nil || schema == nil {
		return
	}

	// for db.Scopes()
	if !stmt.Unscoped {
		for _, c := range schema.CreateClauses {
			stmt.AddClause(c)
		}
	}

	if stmt.SQL.Len() == 0 {
		boundVars := make(map[string]int)
		hasBatch := false
		values := callbacks.ConvertToCreateValues(stmt) // slice: [][]interface{}
		conflict, hasConflict := stmt.Clauses["ON CONFLICT"].Expression.(clause.OnConflict)
		if hasConflict {
			merge := clause_.Merge{}
			// 1. set up From
			merge.From = clause.From{
				Tables: []clause.Table{{Name: DualTableName()}},
			}
			// 2. set up Using
			merge.Using = values
			// 3. set up On
			for _, col := range conflict.Columns {
				merge.On = append(merge.On, clause.Eq{
					Column: clause.Column{Table: stmt.Schema.Table, Name: col.Name},
					Value:  clause.Column{Table: clause_.DefaultExcludeName(), Name: col.Name},
				})
			}
			stmt.AddClauseIfNotExists(merge)
			if !conflict.DoNothing {
				stmt.AddClauseIfNotExists(clause_.WhenMatched{Set: conflict.DoUpdates})
			}
			stmt.AddClauseIfNotExists(clause_.WhenNotMatched{Values: values})
			stmt.Build("MERGE", "WHEN MATCHED", "WHEN NOT MATCHED")
		} else {
			reflectValue := reflect.Indirect(reflect.ValueOf(stmt.Dest))
			switch reflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				// not support returning
				hasBatch = true
				stmt.SQL.Grow(180)
				stmt.AddClauseIfNotExists(clause.Insert{})
				stmt.AddClause(values)
				stmt.Build("INSERT", "VALUES")
			default:
				// just 1 data
				stmt.AddClauseIfNotExists(clause.Insert{Table: clause.Table{Name: stmt.Schema.Table}})
				stmt.AddClause(clause.Values{Columns: values.Columns, Values: [][]interface{}{values.Values[0]}})
				if len(stmt.Schema.FieldsWithDefaultDBValue) > 0 {
					stmt.AddClauseIfNotExists(clause.Returning{
						Columns: slices.Map(schema.FieldsWithDefaultDBValue, func(f *schema2.Field) clause.Column {
							return clause.Column{Name: f.DBName}
						}),
					})
				}
				stmt.Build("INSERT", "VALUES", "RETURNING")
				if len(stmt.Schema.FieldsWithDefaultDBValue) > 0 {
					_, _ = stmt.WriteString(" INTO ")
					for idx, f := range schema.FieldsWithDefaultDBValue {
						if idx > 0 {
							_ = stmt.WriteByte(',')
						}
						boundVars[f.Name] = len(stmt.Vars)
						stmt.AddVar(stmt, sql.Out{Dest: reflect.New(f.FieldType).Interface()})
					}
				}
			}
		}

		if !db.DryRun {
			result, err := stmt.ConnPool.ExecContext(stmt.Context, stmt.SQL.String(), stmt.Vars...)
			if err != nil {
				_ = db.AddError(err)
				return
			}
			if !hasBatch && !hasConflict {
				db.RowsAffected, err = result.RowsAffected()
				if err != nil {
					panic(err)
				}

				refV := stmt.ReflectValue
				if len(stmt.Schema.FieldsWithDefaultDBValue) > 0 {
					for _, f := range schema.FieldsWithDefaultDBValue {
						switch refV.Kind() {
						case reflect.Struct:
							if err = f.Set(stmt.Context, refV, stmt.Vars[boundVars[f.Name]].(sql.Out).Dest); err != nil {
								_ = db.AddError(err)
								return
							}
						case reflect.Map:
							_ = db.AddError(fmt.Errorf("not support map"))
							return
						}
					}
				}
			}
		}
	}

}
