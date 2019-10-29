package cmd

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "checks the health of the Terraform private module registry",
	Long:  "checks the health of the Terraform private module registry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Fatalf("Must specify URL")
		}

		resp, err := http.Get(args[0])
		if err != nil {
			logrus.Fatalf("Healthcheck failed: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			logrus.Fatalf("Healthcheck failed: %d: %s", resp.StatusCode, string(body))
		}

		logrus.Infof("Healthcheck success")
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
}
