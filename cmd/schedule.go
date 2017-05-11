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

var scheduleConfig fastick.SchedulerConfig

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		scheduler := fastick.NewScheduler(scheduleConfig)
		scheduler.Run()
	},
}

func init() {
	RootCmd.AddCommand(scheduleCmd)

	scheduleConfig = fastick.SchedulerConfig{}
	scheduleCmd.Flags().IntVar(&scheduleConfig.Interval, "interval", 15, "Specify the interval in seconds in which tasks are scheduled")

	scheduleCmd.Flags().StringVar(&scheduleConfig.DatabaseHost, "database-host", "localhost", "Database hostname")
	scheduleCmd.Flags().StringVar(&scheduleConfig.DatabaseName, "database-name", "fastick", "Database name")
	scheduleCmd.Flags().StringVar(&scheduleConfig.DatabaseUsername, "database-username", "root", "Database username")
	scheduleCmd.Flags().StringVar(&scheduleConfig.DatabasePassword, "database-password", "", "Database password")

	scheduleCmd.Flags().StringVar(&scheduleConfig.QueueHost, "queue-host", "localhost", "Queue host")
	scheduleCmd.Flags().IntVar(&scheduleConfig.QueuePort, "queue-port", 5672, "Queue port")
	scheduleCmd.Flags().StringVar(&scheduleConfig.QueueUsername, "queue-username", "guest", "Queue username")
	scheduleCmd.Flags().StringVar(&scheduleConfig.QueuePassword, "queue-password", "guest", "Queue password")

}
