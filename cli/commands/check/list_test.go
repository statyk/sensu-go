package check

import (
	"errors"
	"testing"

	"github.com/sensu/sensu-go/cli"
	client "github.com/sensu/sensu-go/cli/client/testing"
	test "github.com/sensu/sensu-go/cli/commands/testing"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	cmd := ListCommand(cli)

	assert.NotNil(cmd, "cmd should be returned")
	assert.NotNil(cmd.RunE, "cmd should be able to be executed")
	assert.Regexp("list", cmd.Use)
	assert.Regexp("checks", cmd.Short)
}

func TestListCommandRunEClosure(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	client := cli.Client.(*client.MockClient)
	client.On("ListChecks").Return([]types.CheckConfig{
		*types.FixtureCheckConfig("name-one"),
		*types.FixtureCheckConfig("name-two"),
	}, nil)

	cmd := ListCommand(cli)
	cmd.Flags().Set("format", "json")
	out, err := test.RunCmd(cmd, []string{})

	assert.NotEmpty(out)
	assert.Contains(out, "name-one")
	assert.Contains(out, "name-two")
	assert.Nil(err)
}

func TestListCommandRunEClosureWithTable(t *testing.T) {
	assert := assert.New(t)
	cli := newCLI()

	check := types.FixtureCheckConfig("name-one")
	check.RuntimeAssets = []string{"asset-one"}

	client := cli.Client.(*client.MockClient)
	client.On("ListChecks").Return([]types.CheckConfig{*check}, nil)

	cmd := ListCommand(cli)
	cmd.Flags().Set("format", "tabular")
	out, err := test.RunCmd(cmd, []string{})

	assert.NotEmpty(out)
	assert.Contains(out, "Name")      // heading
	assert.Contains(out, "Command")   // heading
	assert.Contains(out, "Assets")    // heading
	assert.Contains(out, "name-one")  // check name
	assert.Contains(out, "asset-one") // asset name
	assert.Contains(out, "60")        // interval
	assert.Nil(err)
}

func TestListCommandRunEClosureWithErr(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	client := cli.Client.(*client.MockClient)
	client.On("ListChecks").Return([]types.CheckConfig{}, errors.New("my-err"))

	cmd := ListCommand(cli)
	out, err := test.RunCmd(cmd, []string{})

	assert.NotNil(err)
	assert.Equal("my-err", err.Error())
	assert.Empty(out)
}

func newCLI() *cli.SensuCli {
	cli := test.NewMockCLI()
	config := cli.Config.(*client.MockConfig)
	config.On("Format").Return("json")

	return cli
}
