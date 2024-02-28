package releases

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"slices"
	"time"
)

// The releases package provides functionality around CockroachDB releases, including listing and saving to
// a remote cluster

type Provider struct{}
type Persister struct {
	SqlExecutor *SqlExecutor
}

type Releases []Release

type SortBy int

const (
	SortByVersion     SortBy = iota
	SortByReleaseDate SortBy = iota
)

type Release struct {
	Name          string    `json:"release_name"`
	Withdrawn     bool      `json:"withdrawn"`
	CloudOnly     bool      `json:"cloud_only"`
	ReleaseType   string    `json:"release_type"` // Production or Testing
	ReleaseDate   time.Time `json:"release_date"`
	MajorVersion  string    `json:"major_version"` // e.g., v1.0, v23.1
	Major         int       `json:"major"`
	Minor         int       `json:"minor"`
	Patch         int       `json:"path"`
	BetaRc        string    `json:"beta_rc"`
	BetaRcVersion int       `json:"beta_rc_version"`
}

// NewReleaseFromRemote
func NewReleaseFromRemote(remote RemoteRelease) *Release {
	v := remote.Version()
	return &Release{
		Name:          remote.Name,
		Withdrawn:     remote.Withdrawn,
		CloudOnly:     remote.CloudOnly,
		ReleaseType:   remote.ReleaseType,
		ReleaseDate:   remote.ReleaseDate.Time,
		MajorVersion:  remote.MajorVersion,
		Major:         v.Major,
		Minor:         v.Minor,
		Patch:         v.Patch,
		BetaRc:        v.BetaRc,
		BetaRcVersion: v.BetaRcVersion,
	}
}

func (r *Release) CompareDates(r2 *Release) int {
	if r.ReleaseDate == r2.ReleaseDate {
		return 0
	} else if r.ReleaseDate.Before(r2.ReleaseDate) {
		return -1
	}
	return 1
}

func (r *Release) CompareVersion(r2 *Release) int {
	// Check to see if r is < r2 and default to false
	if r.Major < r2.Major { // Major is less
		return -1
	} else if r.Major == r2.Major { // Majors are equal
		if r.Minor < r2.Minor { // Minor is less
			return -1
		} else if r.Minor == r2.Minor { // Minors are equal
			if r.Patch < r2.Patch { // Patch is less
				return -1
			} else if r.Patch == r2.Patch { // Patches are equal
				if r.BetaRc == "alpha" && (r2.BetaRc == "beta" || r2.BetaRc == "rc" || r2.BetaRc == "") { // Alpha v
					return -1
				} else if r.BetaRc == "beta" && (r2.BetaRc == "rc" || r2.BetaRc == "") { // Beta vs RC or prod
					return -1
				} else if r.BetaRc == "rc" && r2.BetaRc == "" { // RC vs prod
					return -1
				} else if r.BetaRc == r2.BetaRc { // Same alpha, beta or rc, look at version
					if r.BetaRcVersion < r2.BetaRcVersion {
						return -1
					} else if r.BetaRcVersion == r2.BetaRcVersion {
						return 0
					} else {
						return 1
					}
				} else {
					return 1
				}
			} else {
				return 1
			}
		} else {
			return 1
		}
	} else {
		return 1
	}
	return 1
}

func UpdateReleases(pool *pgxpool.Pool) []error {
	p := NewProvider()
	errors := make([]error, 0)
	releases, err := p.GetReleases()
	if err != nil {
		errors = append(errors, err)
		return errors
	}
	persister := NewPersister(pool)
	err = persister.Init()
	if err != nil {
		errors = append(errors, err)
		return errors
	}
	return persister.SaveReleases(releases)
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) GetReleases() ([]Release, error) {
	rp := NewRemoteProvider()
	rrs, err := rp.GetReleases()
	if err != nil {
		return nil, err
	}
	releases := make([]Release, len(rrs))
	for i, rr := range rrs {
		releases[i] = *NewReleaseFromRemote(rr)
	}
	return releases, nil
}

func NewPersister(pool *pgxpool.Pool) *Persister {
	return &Persister{SqlExecutor: NewSqlExecutor(pool)}
}

func (p *Persister) Init() error {
	return p.SqlExecutor.CreateTable()
}

func (p *Persister) SaveReleases(releases []Release) []error {
	errors := make([]error, 0)
	for _, r := range releases {
		err := p.SqlExecutor.UpsertRelease(r)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetRelevantReleases gets the following:
//   - Production releases from the past 2 years that have not been withdrawn
//   - If the most recent major release has no production releases, then it gets all testing releases
//   - Gets most release for all major releases older than 2 years
func (p *Provider) GetRelevantReleases() ([]Release, error) {
	rp := NewRemoteProvider()
	rrs, err := rp.GetReleases()
	if err != nil {
		return nil, err
	}
	releases := make([]Release, len(rrs))
	for i, rr := range rrs {
		releases[i] = *NewReleaseFromRemote(rr)
	}
	return releases, nil
}

func (rs *Releases) SortBy(sort SortBy) {
	switch sort {
	case SortByVersion:
		slices.SortFunc(*rs, func(a, b Release) int {
			return a.CompareVersion(&b)
		})
	case SortByReleaseDate:
		slices.SortFunc(*rs, func(a, b Release) int {
			return a.CompareDates(&b)
		})
	}
}
