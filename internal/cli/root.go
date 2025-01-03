package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/linecard/job/catalog"
	"github.com/linecard/job/internal/cli/job"
	"github.com/linecard/job/internal/config"
	"github.com/linecard/job/internal/output"
	"github.com/linecard/job/internal/parse"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = job.NewRootCmd()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	scripts, err := catalog.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range scripts {
		contents, err := catalog.Read(s)
		if err != nil {
			fmt.Println(output.Colorize(output.Red, err.Error()))
			continue
		}

		desc, args, err := parse.Args(contents)
		if err != nil {
			fmt.Println(output.Colorize(output.Red, err.Error()))
			continue
		}

		rootCmd.AddCommand(job.NewScriptCmd(s, desc, args))
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Version = config.Version + " (" + config.Commit + ")"
}

func initConfig() {
	config.LoadDefaults()

	viper.SetEnvPrefix("job")
	viper.AutomaticEnv()
}
