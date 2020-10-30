package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"php-codesniffer-action/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ignoreWarnings bool

	rootCmd = &cobra.Command{
		Use:          "php-codesniffer-action",
		Short:        "PHP_CodeSniffer Action",
		SilenceUsage: true,
		RunE:         runRoot,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVar(
		&ignoreWarnings,
		"ignore-warnings",
		true,
		"whether annotations are added for warnings",
	)

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_", "-", "_"),
	)
}

func runRoot(cmd *cobra.Command, args []string) error {
	d, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	r := new(internal.Report)
	if err := json.Unmarshal(d, r); err != nil {
		return err
	}

	ms, err := r.CreateMessages()
	if err != nil {
		return err
	}

	for _, m := range ms {
		if _, err := fmt.Fprintln(os.Stdout, m); err != nil {
			return err
		}
	}

	return nil
}
