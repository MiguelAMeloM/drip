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
	"crypto/rand"
	"github.com/MiguelAMeloM/drip/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		cobra.CheckErr(err)
		token, err := cmd.Flags().GetString("token")
		cobra.CheckErr(err)
		engine.RaiseMainServer(token)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 6000, "port")
	serveCmd.Flags().String("token", rand.Text(), "Authentication token")
}
