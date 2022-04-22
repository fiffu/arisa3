package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fiffu/arisa3/lib"
	"github.com/stretchr/testify/assert"
)

const testConfigFileName string = "config_test.yml"
const testConfigFileContents string = "botSecret: sample"

var testDir = lib.MustGetCallerDir()

func mustFindFile(name string) string {
	targetFile := filepath.Join(testDir, string(name))
	if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
		panic("could not find " + targetFile)
	}
	return targetFile
}

func mustWriteFile(name, contents string) {
	targetFile := filepath.Join(testDir, string(name))
	fmt.Println(targetFile)
	bytes := []byte(testConfigFileContents)
	if err := os.WriteFile(targetFile, bytes, 0600); err != nil {
		panic("could not open " + targetFile)
	}
}

func mustDeleteFile(name string) {
	targetFile := filepath.Join(testDir, string(name))
	if err := os.Remove(targetFile); err != nil {
		panic("failed to delete " + targetFile + "err: " + err.Error())
	}
}

func Test_Configure(t *testing.T) {
	mustWriteFile(testConfigFileName, testConfigFileContents)
	defer mustDeleteFile(testConfigFileName)
	path := mustFindFile(testConfigFileName)

	cfg, err := Configure(path)
	assert.NoError(t, err)
	assert.Equal(t, cfg.BotSecret, "sample")
}
