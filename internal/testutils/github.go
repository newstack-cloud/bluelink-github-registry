package testutils

import "fmt"

func GithubAssetURL(assetID int) *string {
	url := fmt.Sprintf(
		"https://api.github.com/repos/two-hundred/celerity-provider-example/releases/assets/%d",
		assetID,
	)
	return &url
}
