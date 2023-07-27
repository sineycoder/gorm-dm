package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	var arrs []*T1
	res := db.Order("id").Where("g_col = ?", true).Find(&arrs)
	fmt.Println(arrs)
	assert.Nil(t, res.Error)
}

func TestQueryJoin(t *testing.T) {
	var arrs1 []*User
	res1 := db.InnerJoins("T1").Find(&arrs1)
	assert.Nil(t, res1.Error)

	var arrs2 []*User
	res2 := db.Preload("T1").Find(&arrs2)
	assert.Nil(t, res2.Error)
}
