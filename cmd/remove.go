/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package cmd

import (
	"fmt"
	"github.com/MiguelAMeloM/drip/engine"
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
