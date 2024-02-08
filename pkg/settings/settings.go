package settings

import (
	"bytes"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// releaseDataURL is the location of the YAML file maintained by the
// docs team where release information is encoded. This data is used
// to render the public CockroachDB releases page. We leverage the
// data in structured format to generate release information used
// for testing purposes.
const releaseDataURL = "https://raw.githubusercontent.com/cockroachdb/docs/main/src/current/_data/releases.yml"

// var namePattern = regexp.MustCompile(`^v(\d+).(\d+).(\d+)-?(beta|rc)?\.?(\d+)$`)
var namePattern = regexp.MustCompile(`^v(\d+).(\d+).(\d+)-?(beta|rc)?\.?(\d+)?$`)

type CustomTime struct {
	time.Time
}

// Release contains the information we extract from the YAML file in
// `releaseDataURL`.
type Release struct {
	Name        string     `yaml:"release_name"`
	Withdrawn   bool       `yaml:"withdrawn"`
	CloudOnly   bool       `yaml:"cloud_only"`
	ReleaseType string     `yaml:"release_type"` // Production or Testing
	ReleaseDate CustomTime `yaml:"release_date"`
}

type Version struct {
	Major         int
	Minor         int
	Patch         int
	BetaRc        string
	BetaRcVersion int
}

type ReleaseFilter struct {
	Widthdrawn  bool
	ReleaseType string
	From        time.Time
	To          time.Time
}

func (ct *CustomTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	layout := "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func ServerForVersion(v string) (testserver.TestServer, error) {
	return testserver.NewTestServer(
		testserver.CustomVersionOpt(v))
}

func GetReleases() ([]Release, error) {
	resp, err := http.Get(releaseDataURL)
	if err != nil {
		return nil, fmt.Errorf("could not download release data: %w", err)
	}
	defer resp.Body.Close()

	var blob bytes.Buffer
	if _, err := io.Copy(&blob, resp.Body); err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var data []Release
	if err := yaml.Unmarshal(blob.Bytes(), &data); err != nil { //nolint:yaml
		return nil, fmt.Errorf("failed to YAML parse release data: %w", err)
	}

	return data, nil
}

func (r *Release) Version() Version {
	matches := namePattern.FindStringSubmatch(r.Name)
	v := Version{}
	if matches != nil {
		major, err := strconv.Atoi(matches[1])
		minor, err := strconv.Atoi(matches[2])
		patch, err := strconv.Atoi(matches[3])
		if err != nil {
			return v
		}
		v.Major = major
		v.Minor = minor
		v.Patch = patch
		v.BetaRc = matches[4]
		if len(matches[5]) > 0 {
			betaRcVersion, err := strconv.Atoi(matches[5])
			if err != nil {
				return v
			}
			v.BetaRcVersion = betaRcVersion
		}
	}

	return v
}

/*
func GetLatestReleases(num int) ([]Release, error) {
	releases, err := GetReleases()
	if err != nil {
		return err
	}
	filtered := make([]Release, 0)
	lastMajor := 0
	lastMinor := 0
	for _, r := range releases {
		major, minor := majorMinorFromReleaseName(r.Name)
		filtered = append(filtered, r)
	}
}

func majorMinorFromReleaseName(n string) (int, int) {

}

*/
