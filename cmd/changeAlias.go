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
	"github.com/MiguelAMeloM/drip/engine"
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
