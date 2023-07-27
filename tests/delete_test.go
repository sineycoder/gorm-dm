package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	res := db.Where("1=1").Delete(&T1{})
	assert.Nil(t, res.Error)
}
