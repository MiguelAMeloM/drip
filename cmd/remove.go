package cmd

import (
	"drip/engine"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
)

var removeCmd = &cobra.Command{
	Use: "deprecate",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("model", cmd.Flags().Lookup("model"))
		cobra.CheckErr(err)

		modelName := viper.GetString("model")

		u, err := url.Parse(urlBase + "/" + engine.Rem)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("model", modelName)

		u.RawQuery = q.Encode()

		done := make(chan bool)
		go spinner(done, "removing model...")

		resp, err := http.Get(u.String())
		done <- true
		if err != nil {
			cobra.CheckErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Println(resp.Status)
			return
		}
		fmt.Println("model removed")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringP("model", "m", "", "model name")

}
