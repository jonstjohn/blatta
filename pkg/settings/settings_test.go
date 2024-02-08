package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerForVersion(t *testing.T) {
	ts, err := ServerForVersion("v23.1.14")
	assert.Nil(t, err)
	pgurl := ts.PGURL()
	assert.NotNil(t, pgurl)

}

func TestGetReleases(t *testing.T) {
	rs, err := GetReleases()
	assert.Nil(t, err)
	for _, r := range rs {
		assert.NotNil(t, r.ReleaseDate)
	}
}

func TestReleaseVersion(t *testing.T) {

	type NamePatternResult struct {
		Major         int
		Minor         int
		Patch         int
		BetaRc        string
		BetaRcVersion int
	}
	type NamePatternTest struct {
		Release Release
		Result  NamePatternResult
	}
	tests := []NamePatternTest{
		NamePatternTest{
			Release: Release{Name: "v23.1.14"},
			Result:  NamePatternResult{Major: 23, Minor: 1, Patch: 14},
		},
		NamePatternTest{
			Release: Release{Name: "v23.2.0-beta.3"},
			Result:  NamePatternResult{Major: 23, Minor: 2, Patch: 0, BetaRc: "beta", BetaRcVersion: 3},
		},
	}
	for _, n := range tests {
		v := n.Release.Version()
		assert.Equal(t, n.Result.Major, v.Major)
		assert.Equal(t, n.Result.Minor, v.Minor)
		assert.Equal(t, n.Result.Patch, v.Patch)

		if len(n.Result.BetaRc) > 0 {
			assert.Equal(t, n.Result.BetaRc, v.BetaRc)
		}

		if n.Result.BetaRcVersion != 0 {
			assert.Equal(t, n.Result.BetaRcVersion, v.BetaRcVersion)
		}
	}
}
