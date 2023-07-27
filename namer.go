package gorm_dm

import (
	"strings"

	"gorm.io/gorm/schema"
)

type Namer struct {
	schema.Namer
}

func ConvertNameToFormat(x string) string {
	return strings.ToUpper(x)
}

func (n Namer) TableName(table string) (name string) {
	return ConvertNameToFormat(n.Namer.TableName(table))
}

func (n Namer) SchemaName(table string) string {
	return ConvertNameToFormat(n.Namer.SchemaName(table))
}

func (n Namer) ColumnName(table, column string) (name string) {
	return ConvertNameToFormat(n.Namer.ColumnName(table, column))
}

func (n Namer) JoinTableName(table string) (name string) {
	return ConvertNameToFormat(n.Namer.JoinTableName(table))
}

func (n Namer) RelationshipFKName(relationship schema.Relationship) (name string) {
	return ConvertNameToFormat(n.Namer.RelationshipFKName(relationship))
}

func (n Namer) CheckerName(table, column string) (name string) {
	return ConvertNameToFormat(n.Namer.CheckerName(table, column))
}

func (n Namer) IndexName(table, column string) (name string) {
	return ConvertNameToFormat(n.Namer.IndexName(table, column))
}
