package gen

import (
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/confgen/cmd/confgen/cli/spin"
	"github.com/brezzgg/go-packages/confgen/internal/generator"
	"github.com/brezzgg/go-packages/confgen/internal/parser"
	"github.com/brezzgg/go-packages/confgen/internal/writer"
	"github.com/spf13/cobra"
)

var GenCmd = &cobra.Command{
	Use:     "generate <schema.hcl> [schema2.hcl ...]",
	Short:   "Generate config form hcl schema",
	Aliases: []string{"gen"},
	Args:    cobra.RangeArgs(1, 2^16),
	RunE: func(cmd *cobra.Command, paths []string) error {
		recursive, err := cmd.Flags().GetBool(genRecursiveFlag)
		if err != nil {
			return err
		}
		stdout, err := cmd.Flags().GetBool(genStdoutFlag)
		if err != nil {
			return err
		}
		wd, err := cmd.Flags().GetString(genWdFlag)
		if err != nil {
			return err
		}
		pattern, err := cmd.Flags().GetString(genPatternFlag)
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

			gen, err := generator.Generate(schema)
			if err != nil {
				spinCh <- fmt.Sprintf("✗ %s: generate error: %s", file, err)
				continue
			}

			var out string
			if !stdout {
				out, err = writer.Write(schema, gen)
				if err != nil {
					spinCh <- fmt.Sprintf("✗ %s: write error: %s", file, err)
					continue
				}
			}

			if out != "stdout" && !stdout {
				spinCh <- fmt.Sprintf("✓ %s -> %s", file, out)
			} else {
				spinCh <- fmt.Sprintf("✓ %s ->", file)
				fmt.Println(gen)
			}
		}

		return nil
	},
}

const (
	genRecursiveFlag  = "recursive"
	genRecursiveShort = "r"
	genRecursiveDesc  = "recursively walk directories"
	genPatternFlag    = "filename-pattern"
	genPatternShort   = "p"
	genPatternDesc    = "specify filenames regular expression"
	genStdoutFlag     = "stdout"
	genStdoutDesc     = "forcedly write result at stdout"
	genWdFlag         = "working-dir"
	genWdShort        = "d"
	genWdDesc         = "specify working directory"
)

func init() {
	GenCmd.PersistentFlags().BoolP(genRecursiveFlag, genRecursiveShort, false, genRecursiveDesc)
	GenCmd.PersistentFlags().StringP(genPatternFlag, genPatternShort, `^.+\.hcl$`, genPatternDesc)
	GenCmd.PersistentFlags().Bool(genStdoutFlag, false, genStdoutDesc)
	GenCmd.PersistentFlags().StringP(genWdFlag, genWdShort, ".", genWdDesc)
}
