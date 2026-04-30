package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func Parse(entries []string, recursive bool, pattern string) (map[string]*schema.Schema, error) {
	result := make(map[string]*schema.Schema)

	patternRe, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("bad pattern: %s", err)
	}

	for _, entry := range entries {
		info, err := os.Stat(entry)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", entry, err)
		}

		if info.IsDir() {
			r, err := processDirectory(entry, recursive, patternRe)
			if err != nil {
				return nil, fmt.Errorf("process directory: %s", err)
			}
			for k, v := range r {
				result[k] = v
			}
		} else {
			if !patternRe.MatchString(info.Name()) {
				continue
			}
			s, err := processFile(entry)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", entry, err)
			}
			result[entry] = s
		}
	}

	return result, nil
}

func processDirectory(entry string, recursive bool, patternRe *regexp.Regexp) (map[string]*schema.Schema, error) {
	result := make(map[string]*schema.Schema)

	if recursive {
		err := filepath.WalkDir(entry, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !patternRe.MatchString(filepath.Base(path)) {
				return nil
			}

			s, err := processFile(path)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}

			rel, err := filepath.Rel(entry, path)
			if err != nil {
				return fmt.Errorf("rel path %s: %w", path, err)
			}

			result[rel] = s
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(entry)
		if err != nil {
			return nil, fmt.Errorf("read dir %s: %w", entry, err)
		}

		for _, d := range entries {
			if d.IsDir() || filepath.Ext(d.Name()) != ".hcl" {
				continue
			}

			path := filepath.Join(entry, d.Name())
			s, err := processFile(path)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}

			result[d.Name()] = s
		}
	}

	return result, nil
}

func processFile(path string) (*schema.Schema, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse hcl file: %w", diags)
	}
	return schema.Decode(f.Body)
}
