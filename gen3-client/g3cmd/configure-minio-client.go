package g3cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uc-cdis/gen3-client/gen3-client/commonUtils"
	"github.com/uc-cdis/gen3-client/gen3-client/logs"
)

func init() {
	var minioCredFile string
	var apiEndpoint string
	var configureCmd = &cobra.Command{
		Use:   "configure-minio-client",
		Short: "Add or modify a configuration profile to your config file",
		Long: `Configuration file located at ~/.gen3/gen3_client_config.ini
	If a field is left empty, the existing value (if it exists) will remain unchanged`,
		Example: `./gen3-client configure --profile=<profile-name> --miniocred=<path-to-credential/miniocred.json> --apiendpoint=https://data.mycommons.org`,
		Run: func(cmd *cobra.Command, args []string) {
			// don't initialize transmission logs for non-uploading related commands
			logs.SetToBoth()

			profileConfig := conf.ReadCredentials(minioCredFile)
			profileConfig.Profile = profile
			apiEndpoint = strings.TrimSpace(apiEndpoint)
			if apiEndpoint[len(apiEndpoint)-1:] == "/" {
				apiEndpoint = apiEndpoint[:len(apiEndpoint)-1]
			}
			parsedURL, err := conf.ValidateUrl(apiEndpoint)
			if err != nil {
				log.Fatalln("Error occurred when validating apiendpoint URL: " + err.Error())
			}

			prefixEndPoint := parsedURL.Scheme + "://" + parsedURL.Host
			err = req.RequestNewAccessToken(prefixEndPoint+commonUtils.FenceAccessTokenEndpoint, &profileConfig)
			if err != nil {
				receivedErrorString := err.Error()
				errorMessageString := receivedErrorString
				if strings.Contains(receivedErrorString, "401") {
					errorMessageString = `Invalid credentials for apiendpoint '` + prefixEndPoint + `': check if your credentials are expired or incorrect`
				} else if strings.Contains(receivedErrorString, "404") || strings.Contains(receivedErrorString, "405") || strings.Contains(receivedErrorString, "no such host") {
					errorMessageString = `The provided apiendpoint '` + prefixEndPoint + `' is possibly not a valid Gen3 data commons`
				}
				log.Fatalln("Error occurred when validating profile config: " + errorMessageString)
			}
			profileConfig.APIEndpoint = apiEndpoint

			// Store user info in ~/.gen3/gen3_client_config.ini
			conf.UpdateConfigFile(profileConfig)
			log.Println(`Profile '` + profile + `' has been configured successfully.`)
			err = logs.CloseMessageLog()
			if err != nil {
				log.Println(err.Error())
			}
		},
	}

	configureCmd.Flags().StringVar(&profile, "profile", "", "Specify profile to use")
	configureCmd.MarkFlagRequired("profile") //nolint:errcheck
	configureCmd.Flags().StringVar(&minioCredFile, "miniocred", "", "Specify the file that contains the MinIO credentials")
	configureCmd.MarkFlagRequired("miniocred") //nolint:errcheck
	configureCmd.Flags().StringVar(&apiEndpoint, "apiendpoint", "", "Specify the API endpoint of the data commons")
	configureCmd.MarkFlagRequired("apiendpoint") //nolint:errcheck
	RootCmd.AddCommand(configureCmd)
}
