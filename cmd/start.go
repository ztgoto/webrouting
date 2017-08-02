// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ztgoto/webrouting/config"
	"github.com/ztgoto/webrouting/http"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Long:  `start http routing server`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.PrepareSetting()
		if err != nil {
			panic(err)
		}
		http.StartServer()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&config.ConfPath, "config", "f", config.DefaultConfPath, "http server config file path")

	// startCmd.Flags().IntVarP(&config.SysConf.MaxProcs, "proces", "c", runtime.NumCPU(), "system max proces")

	// startCmd.Flags().StringVarP(&config.HTTPConf.Addr, "listen", "l", config.DefaultAddr, "http server listen addr")
	// startCmd.Flags().BoolVarP(&config.HTTPConf.Vhost, "vhost", "v", config.DefaultVhost, "http server vhost")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
