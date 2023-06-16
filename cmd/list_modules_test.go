package cmd

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func chdir(t *testing.T, dir string) {
	t.Helper()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.Chdir(pwd)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestListModules(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		chdir(t, "./testdata/single")

		t.Run("path", func(t *testing.T) {
			list, err := listModules(false)
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				".",
			})
		})

		t.Run("name", func(t *testing.T) {
			list, err := listModules(true)
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				"example.com/single",
			})
		})
	})

	t.Run("multi", func(t *testing.T) {
		chdir(t, "./testdata/multi")

		t.Run("path", func(t *testing.T) {
			list, err := listModules(false)
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				".",
				"./module1",
				"./module2",
			})
		})

		t.Run("name", func(t *testing.T) {
			list, err := listModules(true)
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				"example.com/multi",
				"example.com/multi/module1",
				"example.com/multi/module2",
			})
		})
	})
}
