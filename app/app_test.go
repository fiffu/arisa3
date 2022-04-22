package app

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testDependencyInjector struct {
	ctrl *gomock.Controller
}

func (d testDependencyInjector) NewDatabase(dsn string) (database.IDatabase, error) {
	return database.NewMockIDatabase(d.ctrl), nil
}

func (d testDependencyInjector) Bot(token string) (*discordgo.Session, error) {
	return nil, nil
}

func Test_newApp(t *testing.T) {
	mustWriteFile(testConfigFileName, testConfigFileContents)
	defer mustDeleteFile(testConfigFileName)
	configPath := mustFindFile(testConfigFileName)

	ctrl := gomock.NewController(t)
	testDI := testDependencyInjector{ctrl}
	app, sess, err := newApp(testDI, configPath)
	assert.NotNil(t, app)
	assert.Nil(t, sess)
	assert.NoError(t, err)
}
