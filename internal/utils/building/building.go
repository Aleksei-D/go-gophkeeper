package building

import "fmt"

const (
	notAvailable = "N/A"
)

// PrintBuildVersion вывод данных сборки
func PrintBuildVersion(buildVersion, buildDate, buildCommit string) {
	if buildVersion == "" {
		buildVersion = notAvailable
	}
	if buildDate == "" {
		buildVersion = notAvailable
	}
	if buildCommit == "" {
		buildVersion = notAvailable
	}
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
}
