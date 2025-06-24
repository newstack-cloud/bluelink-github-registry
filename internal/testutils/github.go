package testutils

import "fmt"

func GithubAssetURL(assetID int) *string {
	url := fmt.Sprintf(
		"https://api.github.com/repos/newstack-cloud/bluelink-provider-example/releases/assets/%d",
		assetID,
	)
	return &url
}
