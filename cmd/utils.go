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
	"drip/core"
	"drip/core/proxyStats"
	"drip/engine"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Reset   = "\033[0m"
)

var psCmd = &cobra.Command{
	Use: "ps",
	Run: func(cmd *cobra.Command, args []string) {
		url := urlBase + "/" + engine.Ps
		resp, err := http.Get(url)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer resp.Body.Close()
		body := make([]string, 0)
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&body)
		if err != nil {
			cobra.CheckErr(err)
		}
		for _, v := range body {
			fmt.Println(v)
		}
	},
}

var shutDownCmd = &cobra.Command{
	Use: "shutdown",
	Run: func(cmd *cobra.Command, args []string) {
		url := urlBase + "/" + engine.SD
		resp, err := http.Get(url)
		if err != nil {
			cobra.CheckErr(err)
		}
		if resp.StatusCode != 200 {
			cobra.CheckErr(fmt.Errorf("not valid"))
		}
		done := make(chan bool)
		go spinner(done, "shutting down...")
		time.Sleep(time.Second * 5)
		done <- true
	},
}

var updateAutoScalingCmd = &cobra.Command{
	Use: "autoscaling",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		cobra.CheckErr(err)

		seconds := viper.GetInt("seconds")
		Min := viper.GetFloat64("min")
		Max := viper.GetFloat64("max")

		u, err := url.Parse(urlBase + "/" + engine.Autoscaler)
		cobra.CheckErr(err)

		minRPS := strconv.FormatFloat(Min, 'f', 2, 64)

		maxRPS := strconv.FormatFloat(Max, 'f', 2, 64)

		q := u.Query()
		q.Set("seconds", strconv.Itoa(seconds))
		q.Set("min", minRPS)
		q.Set("max", maxRPS)
		u.RawQuery = q.Encode()

		_, err = http.Get(u.String())
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var statsCmd = &cobra.Command{
	Use: "stats",
	Run: func(cmd *cobra.Command, args []string) {
		url := urlBase + "/" + engine.Stats
		resp, err := http.Get(url)
		if err != nil {
			cobra.CheckErr(err)
		}
		if resp.StatusCode != 200 {
			cobra.CheckErr(fmt.Errorf("not valid"))
		}
		data := make(map[string]*proxyStats.ProxyStats)
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			cobra.CheckErr(err)
		}

		fmt.Printf("%s%-20s  %15s  %15s  %15s	%15s%s\n", Blue, "NAME", "NUM OF REQ", "MEAN LATENCY", "REQ/S", "ACT SERVERS", Reset)
		fmt.Printf("%s%-20s  %15s  %15s  %15s	%15s%s\n", Cyan, "====", "==========", "=============", "=====", "===========", Reset)

		for modelName, stats := range data {
			fmt.Printf("%s%-20s%s  %s\n", Green, modelName, Reset, stats)
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(shutDownCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(updateAutoScalingCmd)

	updateAutoScalingCmd.Flags().IntP("seconds", "s", int(core.IntervalForMonitoring/time.Second), "Number of seconds to update")
	updateAutoScalingCmd.Flags().Float64P("min", "", core.MinRequestsPerSecond, "Number of minutes to update")
	updateAutoScalingCmd.Flags().Float64P("max", "", core.MinRequestsPerSecond, "Number of minutes to update")
}

func spinner(done chan bool, message string) {
	frames := []string{
		"⠋", "⠙", "⠹", "⠸", "⠼",
		"⠴", "⠦", "⠧", "⠇", "⠏",
	}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r\033[2K✔ %s\n", message) // Limpia línea y muestra mensaje de éxito
			return
		default:
			fmt.Printf("\r\033[36m%s\033[0m %s", frames[i%len(frames)], message)
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
