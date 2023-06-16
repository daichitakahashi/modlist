package cmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestListPackages(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		chdir(t, "./testdata/single")

		t.Run("name", func(t *testing.T) {
			list, err := listModules()
			assert.NilError(t, err)

			list, err = expand(list, func(s string) ([]string, error) {
				return listPackages(s)
			})
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				"example.com/single",
				"example.com/single/package1",
				"example.com/single/package2",
			})
		})
	})

	t.Run("multi", func(t *testing.T) {
		chdir(t, "./testdata/multi")

		t.Run("name", func(t *testing.T) {
			list, err := listModules()
			assert.NilError(t, err)

			list, err = expand(list, func(s string) ([]string, error) {
				return listPackages(s)
			})
			assert.NilError(t, err)
			assert.DeepEqual(t, list, []string{
				"example.com/multi",
				"example.com/multi/module1",
				"example.com/multi/module1/package1",
				"example.com/multi/module1/package2",
				"example.com/multi/module2",
				"example.com/multi/module2/internal/pkg",
				"example.com/multi/package1",
				"example.com/multi/package2",
			})
		})
	})
}
