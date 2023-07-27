package tests

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tmps := []*T1{
		{
			ACol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			BCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			CCol: int64(1),
			DCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			ECol: time.Now(),
			FCol: "{}",
			GCol: false,
		},
		{
			ACol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			BCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			CCol: int64(2),
			DCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			ECol: time.Now(),
			FCol: "{}",
			GCol: true,
		},
	}
	//tmps := &T1{
	//	ACol: "16904291616257080568716115743334",
	//	BCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
	//	CCol: int64(1),
	//	DCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
	//	ECol: time.Now(),
	//	FCol: "{}",
	//	GCol: false,
	//}
	res := db.Model(&T1{}).Create(&tmps)
	//res := db.Clauses(clause.OnConflict{
	//	Columns:   []clause.Column{{Name: "A_COL"}},
	//	DoUpdates: clause.AssignmentColumns([]string{"B_COL", "C_COL", "D_COL", "E_COL"}),
	//}).Create(&tmps)
	assert.Nil(t, res.Error)
}

func BenchmarkSingleCreate(b *testing.B) {
	for i := 0; i < 100; i++ {
		db.Create(&T1{
			ACol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			BCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			CCol: int64(1),
			DCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			ECol: time.Now(),
			FCol: "{}",
			GCol: false,
		})
	}
}

func BenchmarkBatchCreate(b *testing.B) {
	arr := make([]*T1, 100)
	for i := 0; i < len(arr); i++ {
		arr[i] = &T1{
			ACol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			BCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			CCol: int64(1),
			DCol: strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10),
			ECol: time.Now(),
			FCol: "{}",
			GCol: false,
		}
	}
	db.Create(&arr)
}
