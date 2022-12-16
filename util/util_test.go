package util_test

import (
	"testing"
	"time"

	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/assert"
)

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

func TestUtil(t *testing.T) {
	intValue := 1
	assert.Equal(t, &intValue, util.ValueToPointer(intValue))
	assert.Equal(t, intValue, util.PointerToValue(&intValue))

	stringValue := "1"
	assert.Equal(t, &stringValue, util.ValueToPointer(stringValue))
	assert.Equal(t, stringValue, util.PointerToValue(&stringValue))

	timeValue := time.Now()
	assert.Equal(t, &timeValue, util.ValueToPointer(timeValue))
	assert.Equal(t, timeValue, util.PointerToValue(&timeValue))
}
