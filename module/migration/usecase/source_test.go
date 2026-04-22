package usecase

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestMaxSourceVersion(t *testing.T) {
	tests := []struct {
		name string
		fs   fstest.MapFS
		path string
		want uint
	}{
		{
			name: "finds the highest numeric prefix in a standard migration set",
			fs: fstest.MapFS{
				"migrations/mysql/1_create_users.up.sql":   {Data: []byte("")},
				"migrations/mysql/1_create_users.down.sql": {Data: []byte("")},
				"migrations/mysql/7_add_index.up.sql":      {Data: []byte("")},
				"migrations/mysql/7_add_index.down.sql":    {Data: []byte("")},
				"migrations/mysql/3_alter_table.up.sql":    {Data: []byte("")},
			},
			path: "migrations/mysql",
			want: 7,
		},
		{
			name: "empty directory yields zero",
			fs: fstest.MapFS{
				"migrations/mysql/.keep": {Data: []byte("")},
			},
			path: "migrations/mysql",
			want: 0,
		},
		{
			name: "ignores files without a leading integer",
			fs: fstest.MapFS{
				"migrations/mysql/notes.txt": {Data: []byte("")},
				"migrations/mysql/README.md": {Data: []byte("")},
			},
			path: "migrations/mysql",
			want: 0,
		},
		{
			name: "missing directory yields zero, not a panic",
			fs:   fstest.MapFS{},
			path: "migrations/mysql",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maxSourceVersion(tt.fs, tt.path))
		})
	}
}
