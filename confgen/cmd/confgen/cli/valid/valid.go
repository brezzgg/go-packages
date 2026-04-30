package valid

import (
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/confgen/cmd/confgen/cli/spin"
	"github.com/brezzgg/go-packages/confgen/internal/generator"
	"github.com/brezzgg/go-packages/confgen/internal/parser"
	"github.com/spf13/cobra"
)

var (
	ValidCmd = &cobra.Command{
		Use:     "validate <schema.hcl> [schema2.hcl ...]",
		Short:   "Validate hcl schema syntax",
		Aliases: []string{"val", "valid"},
		Args:    cobra.RangeArgs(1, 2^16),
		RunE: func(cmd *cobra.Command, paths []string) error {
			recursive, err := cmd.Flags().GetBool(validRecursiveFlag)
			if err != nil {
				return err
			}
			wd, err := cmd.Flags().GetString(validWdFlag)
			if err != nil {
				return err
			}
			pattern, err := cmd.Flags().GetString(validPatternFlag)
			if err != nil {
				return err
			}

			if err := os.Chdir(wd); err != nil {
				return fmt.Errorf("failed to change working directory: %s", err)
			}

			schemas, err := parser.Parse(paths, recursive, pattern)
			if err != nil {
				fmt.Printf("parse file error: %s\n", err)
			}

			for file, schema := range schemas {
				spinCh := make(chan string)
				done := spin.Spin(file, spinCh)
				defer func() {
					<-done
				}()

				_, err := generator.Generate(schema)
				if err != nil {
					spinCh <- fmt.Sprintf("✗ %s: generate error: %s", file, err)
				} else {
					spinCh <- fmt.Sprintf("✓ %s", file)
				}
			}

			return nil
		},
	}
)

const (
	validRecursiveFlag  = "recursive"
	validRecursiveShort = "r"
	validRecursiveDesc  = "recursively walk directories"
	validPatternFlag    = "filename-pattern"
	validPatternShort   = "p"
	validPatternDesc    = "specify filenames regular expression"
	validWdFlag         = "working-dir"
	validWdShort        = "d"
	validWdDesc         = "specify working directory"
)

func init() {
	ValidCmd.PersistentFlags().BoolP(validRecursiveFlag, validRecursiveShort, false, validRecursiveDesc)
	ValidCmd.PersistentFlags().StringP(validPatternFlag, validPatternShort, `^.+\.hcl$`, validPatternDesc)
	ValidCmd.PersistentFlags().StringP(validWdFlag, validWdShort, ".", validWdDesc)
}
