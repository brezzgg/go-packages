package sample

import (
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/confgen/internal/generator"
	"github.com/brezzgg/go-packages/confgen/internal/parser"
	"github.com/brezzgg/go-packages/confgen/internal/schema"
	"github.com/spf13/cobra"
)

var (
	SampleCmd = &cobra.Command{
		Use:   "sample <schema.hcl> <yaml|json>",
		Short: "Generate example config file based on hcl schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := cmd.Flags().GetString(sampleWdFlag)
			if err != nil {
				return err
			}

			if err := os.Chdir(wd); err != nil {
				return fmt.Errorf("failed to change working directory: %s", err)
			}

			schemas, err := parser.Parse([]string{args[0]}, false, `.+`)
			if err != nil {
				fmt.Printf("parse file error: %s\n", err)
			}
			if len(schemas) < 1 {
				return nil
			}
			var schema *schema.Schema
			for _, s := range schemas {
				schema = s
				break
			}

			gen, err := generator.GenerateSample(schema, args[1])
			if err != nil {
				return fmt.Errorf("generate error: %s", err)
			} else {
				fmt.Println(gen)
			}

			return nil
		},
	}
)

const (
	sampleWdFlag  = "working-dir"
	sampleWdShort = "d"
	sampleWdDesc  = "specify working directory"
)

func init() {
	SampleCmd.Flags().StringP(sampleWdFlag, sampleWdShort, ".", sampleWdDesc)
}
