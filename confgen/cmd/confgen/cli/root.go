package cli

import (
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/confgen/cmd/confgen/cli/gen"
	"github.com/brezzgg/go-packages/confgen/cmd/confgen/cli/sample"
	"github.com/brezzgg/go-packages/confgen/cmd/confgen/cli/valid"
	"github.com/spf13/cobra"
)

var (
	rootCmd = cobra.Command{
		Use: "confgen",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func init() {
	addSubcommand()
}

func addSubcommand() {
	rootCmd.AddCommand(gen.GenCmd)
	rootCmd.AddCommand(valid.ValidCmd)
	rootCmd.AddCommand(sample.SampleCmd)
}
