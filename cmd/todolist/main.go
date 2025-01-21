package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "todolist",
	Short:            description,
	Long:             ``,
	PersistentPreRun: doPersistentPreRun,
	SilenceErrors:    true, // allows us to log errors uniformly using the logger, without a duplicate from cobra
}

var (
	debug bool
)

const (
	description = "todolist server"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
}

func doPersistentPreRun(cmd *cobra.Command, args []string) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg(description + " failed")
		os.Exit(1)
	}
	
}
