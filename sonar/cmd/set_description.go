package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/felicianotech/sonar/sonar/docker"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var setDescriptionCmd = &cobra.Command{
	Use:   "summary <image-name> <summary-string>",
	Short: "Set the summary for an image on Docker Hub",
	Long:  "Limited to 100 characters.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// Escape file content for use in JSON
		content := []byte(strconv.Quote(args[1]))

		content = append([]byte("{\"description\": "), content[:len(content)-1]...)
		content = append(content, []byte("\"}")...)

		req, err := http.NewRequest("PATCH", "https://hub.docker.com/v2/repositories/"+args[0]+"/", bytes.NewBuffer(content))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		if viper.Get("user") == nil || viper.Get("pass") == nil || len(viper.Get("user").(string)) == 0 || len(viper.Get("pass").(string)) == 0 {
			log.Fatal("This command requires Docker Hub credentials to be set in your environment.")
		}

		resp, err := docker.SendRequest(req, viper.Get("user").(string), viper.Get("pass").(string))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			log.Fatal("There was an error updating the summary. Code " + resp.Status)
		} else {
			fmt.Printf("Successfully updated with code %d.\n", resp.StatusCode)
		}
	},
}

func init() {
	setCmd.AddCommand(setDescriptionCmd)
}
