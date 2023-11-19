package app

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testDependencyInjector struct {
	ctrl *gomock.Controller
}

func (d testDependencyInjector) NewDatabase(ctx context.Context, dsn string) (database.IDatabase, error) {
	return database.NewMockIDatabase(d.ctrl), nil
}

func (d testDependencyInjector) NewInstrumentationClient(ctx context.Context) (instrumentation.Client, error) {
	return instrumentation.NewInstrumentationClient(ctx)
}

func (d testDependencyInjector) Bot(token string, debugMode bool) (*discordgo.Session, error) {
	return nil, nil
}

func Test_newApp(t *testing.T) {
	mustWriteFile(testConfigFileName, testConfigFileContents)
	defer deleteFile(testConfigFileName)
	configPath := mustFindFile(testConfigFileName)

	ctrl := gomock.NewController(t)
	testDI := testDependencyInjector{ctrl}
	app, err := newApp(context.Background(), testDI, configPath)
	assert.NotNil(t, app)
	assert.Nil(t, app.BotSession())
	assert.NoError(t, err)
}
