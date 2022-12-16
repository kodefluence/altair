package util_test

import (
	"testing"
	"time"

	"github.com/kodefluence/altair/testhelper"
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

func TestUtil_ReadFileContent(t *testing.T) {

	t.Run("Given valid file name, then it does return file content as a string", func(t *testing.T) {
		testhelper.GenerateTempTestFiles("./test_file/", "test", "valid_file.txt", 0644)
		content, err := util.ReadFileContent("./test_file/valid_file.txt")
		assert.Nil(t, err)
		assert.Equal(t, []byte("test"), content)
		testhelper.RemoveTempTestFiles("./test_file/valid_file.txt")
	})

	t.Run("Given invalid file name, then it does return error", func(t *testing.T) {
		_, err := util.ReadFileContent("./test_file/invalid_file.txt")
		assert.NotNil(t, err)
	})
}
