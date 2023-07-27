package gorm_dm

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "code.byted.org/videoarch-onpremise/dm"
	"github.com/sineycoder/gorm-dm/utils"
	"github.com/sineycoder/gorm-dm/utils/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	gorm_utils "gorm.io/gorm/utils"
)

type Config struct {
	DriverName        string
	DSN               string
	Conn              gorm.ConnPool // *sql.DB
	DefaultStringSize uint          // varchar2 default size if not specify
	DBVersion         string        // save db version
	FormatTimestamp   bool          // using <TIMESTAMP 'yyyy-MM-dd HH:mm:ss.999'> when gorm time.Time to db timestamp
}

type Dialector struct {
	*Config
}

func Open(dsn string) gorm.Dialector {
	return &Dialector{Config: &Config{DSN: dsn}}
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config}
}

func BuildDsn(server string, port int, user, password string, options map[string]string) string {
	ret := fmt.Sprintf("dm://%s:%s@%s?autoCommit=true", url.PathEscape(user), url.PathEscape(password),
		net.JoinHostPort(server, strconv.Itoa(port)))
	if options != nil {
		ret += "?"
		for key, val := range options {
			val = strings.TrimSpace(val)
			for _, temp := range strings.Split(val, ",") {
				temp = strings.TrimSpace(temp)
				if strings.ToUpper(key) == "SERVER" {
					ret += fmt.Sprintf("%s=%s&", key, temp)
				} else {
					ret += fmt.Sprintf("%s=%s&", key, url.QueryEscape(temp))
				}
			}
		}
		ret = strings.TrimRight(ret, "&")
	}
	return ret
}

func DualTableName() string {
	return "DUAL"
}

func (d Dialector) Name() string {
	return "dm"
}

func (d Dialector) Initialize(db *gorm.DB) (err error) {
	db.NamingStrategy = Namer{Namer: utils.IfThen[schema.Namer](db.NamingStrategy == nil, schema.NamingStrategy{}, db.NamingStrategy)}
	d.DefaultStringSize = utils.IfThen(d.DefaultStringSize == 0, 1024, d.DefaultStringSize)

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		CreateClauses: []string{"INSERT", "VALUES", "ON CONFLICT", "RETURNING"},
		UpdateClauses: []string{"UPDATE", "SET", "WHERE", "RETURNING"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE", "RETURNING"},
	})

	if d.Conn != nil {
		db.ConnPool = d.Conn
	} else {
		d.DriverName = utils.IfThen(len(d.DriverName) == 0, "dm", d.DriverName)
		db.ConnPool, err = sql.Open(d.DriverName, d.DSN)
		if err != nil {
			return err
		}
	}

	if err = db.ConnPool.QueryRowContext(context.Background(), "SELECT BANNER FROM v$VERSION LIMIT 1").Scan(&d.DBVersion); err != nil {
		return err
	}

	if err = db.Callback().Create().Replace("gorm:create", Create); err != nil {
		return err
	}

	return nil
}

func (d Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	panic("implement me!")
}

// DataTypeOf for migrator
func (d Dialector) DataTypeOf(field *schema.Field) string {
	delete(field.TagSettings, "RESTRICT")

	var sqlType string

	switch field.DataType {
	case schema.Bool:
		sqlType = "INTEGER(1,0)"
	case schema.Int, schema.Uint:
		if field.Size <= 8 {
			sqlType = "INTEGER(3,0)"
		} else if field.Size <= 16 {
			sqlType = "INTEGER(5,0)"
		} else if field.Size <= 32 {
			sqlType = "INTEGER(10,0)"
		} else {
			sqlType = "INTEGER"
		}
	case schema.Float:
		sqlType = "FLOAT"
	case schema.String, "VARCHAR2":
		size := field.Size
		defaultSize := d.DefaultStringSize
		if size == 0 {
			if defaultSize > 0 {
				size = int(defaultSize)
			} else {
				hasIndex := field.TagSettings["INDEX"] != "" || field.TagSettings["UNIQUE"] != ""
				if field.PrimaryKey || field.HasDefaultValue || hasIndex {
					size = 191
				}
			}
		}

		if size >= 2000 {
			sqlType = "CLOB"
		} else {
			sqlType = fmt.Sprintf("VARCHAR2(%d)", size)
		}

	case schema.Time:
		sqlType = "TIMESTAMP"

	case schema.Bytes:
		sqlType = "BLOB"
	default:
		sqlType = string(field.DataType)

		if strings.EqualFold(sqlType, "text") {
			sqlType = "CLOB"
		}

		if sqlType == "" {
			panic(fmt.Sprintf("invalid sql type %s (%s) for dameng", field.FieldType.Name(), field.FieldType.String()))
		}

	}
	if val, ok := field.TagSettings["AUTOINCREMENT"]; ok && gorm_utils.CheckTruth(val) {
		sqlType += " IDENTITY"
	} else {
		notNull, _ := field.TagSettings["NOT NULL"]
		unique, _ := field.TagSettings["UNIQUE"]
		additionalType := fmt.Sprintf("%s %s", notNull, unique)
		if value, ok := field.TagSettings["DEFAULT"]; ok {
			additionalType = fmt.Sprintf("%s %s %s%s", "DEFAULT", value, additionalType, func() string {
				if value, ok := field.TagSettings["COMMENT"]; ok {
					return " COMMENT " + value
				}
				return ""
			}())
		}
		sqlType = fmt.Sprintf("%v %v", sqlType, additionalType)
	}
	return sqlType
}

func (d Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}

func (d Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	_, _ = writer.WriteString(":")
	_, _ = writer.WriteString(strconv.Itoa(len(stmt.Vars)))
}

func (d Dialector) QuoteTo(writer clause.Writer, str string) {
	_, _ = writer.WriteString(str)
}

var (
	numericPlaceholder = regexp.MustCompile(`:(\d+)`)
	tmFmtWithMS        = "2006-01-02 15:04:05.999"
	tmFmtZero          = "0000-00-00 00:00:00"
	nullStr            = "NULL"
)

func (d Dialector) Explain(sql string, vars ...interface{}) string {
	sql = logger.ExplainSQL(sql, numericPlaceholder, `'`, slices.Map(vars, func(v interface{}) interface{} {
		switch v := v.(type) {
		case bool:
			if v {
				return 1
			}
			return 0
		case time.Time:
			if v.IsZero() {
				return "TIMESTAMP " + tmFmtZero
			} else {
				return "TIMESTAMP " + v.Format(tmFmtWithMS)
			}
		case *time.Time:
			if v != nil {
				if v.IsZero() {
					return "TIMESTAMP " + tmFmtZero
				} else {
					return "TIMESTAMP " + v.Format(tmFmtWithMS)
				}
			} else {
				return nullStr
			}
		default:
			return v
		}
	})...)

	return strings.ReplaceAll(sql, `'TIMESTAMP `, `TIMESTAMP '`)
}

func (d Dialector) SavePoint(tx *gorm.DB, name string) error {
	tx.Exec("SAVEPOINT " + name)
	return tx.Error
}

func (d Dialector) RollbackTo(tx *gorm.DB, name string) error {
	tx.Exec("ROLLBACK TO SAVEPOINT " + name)
	return tx.Error
}
