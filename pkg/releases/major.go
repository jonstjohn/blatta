package releases

type MajorVersion struct {
	MajorVersion           string
	FirstTestingRelease    *Release
	FirstProductionRelease *Release
}

type MajorVersionSummary struct {
	MajorVersions []MajorVersion
	LatestRelease *Release
}

/*
func GetMajorVersionSummary(releases Releases) (MajorVersionSummary, error) {

}

*/
