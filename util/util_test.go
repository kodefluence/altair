package util_test

import (
	"testing"
	"time"

	"github.com/codefluence-x/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {

	assert.NotPanics(t, func() {
		util.IntToPointer(1)
		util.TimeToPointer(time.Now())
		util.StringToPointer("blablabla")
		util.SHA1()
	})
}
