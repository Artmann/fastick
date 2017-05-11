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

var workConfig fastick.WorkConfig

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		manager := fastick.NewWorkManager(workConfig)
		manager.Run()
	},
}

func init() {
	RootCmd.AddCommand(workCmd)

	workConfig = fastick.WorkConfig{}

	workCmd.Flags().IntVar(&workConfig.WorkerCount, "worker-count", 10, "Specify the number of workers")

	workCmd.Flags().StringVar(&workConfig.QueueHost, "queue-host", "localhost", "Queue host")
	workCmd.Flags().IntVar(&workConfig.QueuePort, "queue-port", 5672, "Queue port")
	workCmd.Flags().StringVar(&workConfig.QueueUsername, "queue-username", "guest", "Queue username")
	workCmd.Flags().StringVar(&workConfig.QueuePassword, "queue-password", "guest", "Queue password")

}
