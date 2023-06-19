package golangci

import (
	"log"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMain(m *testing.M) {
	err := os.Chdir("testdata")
	if err != nil {
		log.Panic(err)
	}

	m.Run()
}

func TestReadConfig(t *testing.T) {
	t.Parallel()

	cfg, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	assert.NilError(t, err)
	assert.DeepEqual(t, cfg.Run.SkipDirs, []string{
		"internal",
	})
	assert.Assert(t, *cfg.Run.SkipDirsUseDefault)
}
