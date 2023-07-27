package tests

import (
	"time"

	gorm_dm "github.com/sineycoder/gorm-dm"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	url := gorm_dm.BuildDsn("10.174.252.190", 31406, "siney", "siney123123", nil)
	dba, err := gorm.Open(gorm_dm.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = dba
}

type T1 struct {
	ID   int64     // NUMBER
	ACol string    // NVARCHAR2
	BCol string    // NCLOB
	CCol int64     // NUMBER
	DCol string    // NVARCHAR2
	ECol time.Time `gorm:"default:SYSTIMESTAMP"` // timestamp
	FCol string
	GCol bool
}

type User struct {
	ID      int64
	Name    string
	Age     int
	OtherID int64
	T1      T1 `gorm:"foreignKey:OtherID;references:ID"`
}
