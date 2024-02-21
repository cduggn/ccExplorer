package awsservice

import "github.com/spf13/viper"

func Profile() string {
	awsProfile := viper.GetString("aws_profile")
	if awsProfile == "" {
		awsProfile = "not-provided"
	}

	return awsProfile
}
