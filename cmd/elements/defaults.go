package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func defaults(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" defaults", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Gather defaults for an element set from metadata.json files in the repo")

	var regexStr, inDir string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
		// The defaults service expects a JSON document where the keys are
		// element names, and the values are the defaults dictionaries for those
		// elements. We can get these from the metadata.json files throughout
		// the Elements repo.
		defaults := make(map[string]interface{})

		for _, repo := range wf.Repositories {
			err := repo.Walk(func(f *workflow.File) error {
				dir := filepath.Dir(f.Name)
				if (!workflow.IsAnthaMetadata(f.Name)) || !regex.MatchString(dir) {
					return nil
				}

				r, err := f.Contents()
				if err != nil {
					return err
				}

				bs, err := ioutil.ReadAll(r)
				if err != nil {
					return err
				}

				var doc map[string]interface{}
				err = json.Unmarshal(bs, &doc)
				if err != nil {
					return err
				}

				name, ok := doc["name"].(string)
				if !ok {
					return fmt.Errorf("Got unexpected data in name field of %v: expected string, got %v", f.Name, reflect.TypeOf(doc["name"]))
				}
				defaults[name] = doc["defaults"]

				return nil
			})
			if err != nil {
				return err
			}
		}

		bs, err := json.Marshal(defaults)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(os.Stdout)
		_, err = w.Write(bs)
		if err != nil {
			return err
		}

		return nil
	}
}
