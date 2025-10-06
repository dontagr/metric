package print

import "fmt"

func PrintVersion(buildVersion string, buildDate string, buildCommit string) {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: \"%s\"\n", buildVersion)

	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: \"%s\"\n", buildDate)

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: \"%s\"\n", buildCommit)
}
