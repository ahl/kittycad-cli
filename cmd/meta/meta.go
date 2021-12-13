package meta

import (
	"github.com/MakeNowJust/heredoc"
	cmdInstance "github.com/kittycad/cli/cmd/meta/instance"
	"github.com/kittycad/cli/pkg/cli"
	"github.com/spf13/cobra"
)

// NewCmdMeta returns a new instance of the meta command.
func NewCmdMeta(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta <command>",
		Short: "Meta information",
		Long:  `Get information about sessions, servers, and instances. This is best used for debugging authentication sessions, etc`,
		Example: heredoc.Doc(`
			$ kittycad meta instance
		`),
	}

	cmd.AddCommand(cmdInstance.NewCmdInstance(cli))

	return cmd
}
