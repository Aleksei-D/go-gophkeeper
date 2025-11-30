package building

import "fmt"

const (
	notAvaible = "N/A"
)

func PrintBuildVersion(buildVersion, buildDate, buildCommit string) {
	if buildVersion == "" {
		buildVersion = notAvaible
	}
	if buildDate == "" {
		buildVersion = notAvaible
	}
	if buildCommit == "" {
		buildVersion = notAvaible
	}
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
}
