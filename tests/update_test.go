package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	t1 := T1{}
	res := db.Model(&t1).Where("1=1").Update("B_COL", "test1")
	assert.Nil(t, res.Error)
}
