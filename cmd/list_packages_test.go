package cmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestListPackageNames(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		chdir(t, "./testdata/single")

		list, err := listModules(false)
		assert.NilError(t, err)

		list, err = expand(list, func(s string) ([]string, error) {
			return listPackageNames(s)
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, list, []string{
			"example.com/single",
			"example.com/single/package1",
			"example.com/single/package2",
		})
	})

	t.Run("multi", func(t *testing.T) {
		chdir(t, "./testdata/multi")

		list, err := listModules(false)
		assert.NilError(t, err)

		list, err = expand(list, func(s string) ([]string, error) {
			return listPackageNames(s)
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, list, []string{
			"example.com/multi",
			"example.com/multi/package1",
			"example.com/multi/package2",
			"example.com/multi/module1",
			"example.com/multi/module1/package1",
			"example.com/multi/module1/package2",
			"example.com/multi/module2",
			"example.com/multi/module2/internal/pkg",
		})
	})
}

func TestListPackagePaths(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		chdir(t, "./testdata/single")

		list, err := listModules(true)
		assert.NilError(t, err)

		list, err = expand(list, func(s string) ([]string, error) {
			return listPackagePaths(s)
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, list, []string{
			".",
			"./package1",
			"./package2",
		})
	})

	t.Run("multi", func(t *testing.T) {
		chdir(t, "./testdata/multi")

		list, err := listModules(true)
		assert.NilError(t, err)

		list, err = expand(list, func(s string) ([]string, error) {
			return listPackagePaths(s)
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, list, []string{
			".",
			"./package1",
			"./package2",
			"./module1",
			"./module1/package1",
			"./module1/package2",
			"./module2",
			"./module2/internal/pkg",
		})
	})
}
