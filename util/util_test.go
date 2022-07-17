package util_test

import (
	"testing"
	"time"

	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {

	assert.NotPanics(t, func() {
		util.IntToPointer(1)
		util.TimeToPointer(time.Now())
		util.StringToPointer("blablabla")
		util.SHA1()
	})

	assert.Equal(t, nil, util.PointerToInt(nil))
	assert.Equal(t, nil, util.PointerToTime(nil))
	assert.Equal(t, nil, util.PointerToString(nil))
	assert.Equal(t, 1, util.PointerToInt(util.IntToPointer(1)))
	now := time.Now()
	assert.Equal(t, now, util.PointerToTime(util.TimeToPointer(now)))
	assert.Equal(t, "xx", util.PointerToString(util.StringToPointer("xx")))
}
