package cmd

import (
	"drip/engine"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

var changeAliasCmd = &cobra.Command{
	Use:  "ch-al [stage] [alias]",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		stage := args[0]
		alias := args[1]
		u, err := url.Parse(urlBase + "/" + engine.AliasCh)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("stage", stage)
		q.Set("alias", alias)

		u.RawQuery = q.Encode()

		_, err = http.Get(u.String())
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(changeAliasCmd)
}
