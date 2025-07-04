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
	"net/http"
	"net/url"
)

var addCmd = &cobra.Command{
	Use: "release",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		cobra.CheckErr(err)

		modelName := viper.GetString("model")
		releaseType := viper.GetString("release_type")
		backendType := viper.GetString("backend_type")
		deadline := viper.GetString("deadline")
		prob := viper.GetString("prob")

		u, err := url.Parse(urlBase + "/" + engine.Add)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("model", modelName)
		q.Set("release_type", releaseType)
		q.Set("backend_type", backendType)
		q.Set("deadline", deadline)
		q.Set("prob", prob)

		u.RawQuery = q.Encode()

		done := make(chan bool)

		go spinner(done, "adding model...")

		resp, err := http.Get(u.String())
		done <- true
		fmt.Print("\r")

		if err != nil {
			cobra.CheckErr(err)
		}
		if resp.StatusCode != http.StatusCreated {
			fmt.Println(resp.Status)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("model", "m", "", "model name")
	addCmd.Flags().StringP("release_type", "r", "stable", "release type")
	addCmd.Flags().StringP("backend_type", "b", "mlflow", "backend type")
	addCmd.Flags().StringP("deadline", "t", "", "deadline")
	addCmd.Flags().StringP("prob", "p", "0.5", "dev probability")
}
