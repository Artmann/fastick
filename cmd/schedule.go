// Copyright © 2017 Christoffer Artmann <christoffer@artmann.co>
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
	"github.com/artmann/fastick/fastick"
)

var config fastick.SchedulerConfig

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		scheduler := fastick.NewScheduler(config)
		scheduler.Run()
	},
}

func init() {
	RootCmd.AddCommand(scheduleCmd)

	config = fastick.SchedulerConfig{}
	scheduleCmd.Flags().IntVar(&config.Interval, "interval", 15, "Specify the interval in seconds in which tasks are scheduled")

	scheduleCmd.Flags().StringVar(&config.DatabaseHost, "database-host", "localhost", "Database hostname")
	scheduleCmd.Flags().StringVar(&config.DatabaseName, "database-name", "fastick", "Database name")
	scheduleCmd.Flags().StringVar(&config.DatabaseUsername, "database-username", "root", "Database username")
	scheduleCmd.Flags().StringVar(&config.DatabasePassword, "database-password", "", "Database password")

	scheduleCmd.Flags().StringVar(&config.QueueHost, "queque-host", "localhost", "Queue host")
	scheduleCmd.Flags().IntVar(&config.QueuePort, "queque-port", 5672, "Queue port")
	scheduleCmd.Flags().StringVar(&config.QueueUsername, "queque-username", "guest", "Queue username")
	scheduleCmd.Flags().StringVar(&config.QueuePassword, "queque-password", "guest", "Queue password")

}
