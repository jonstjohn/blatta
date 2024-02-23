package releases

type Provider struct{}

type Release struct {
	Name          string     `json:"release_name"`
	Withdrawn     bool       `json:"withdrawn"`
	CloudOnly     bool       `json:"cloud_only"`
	ReleaseType   string     `json:"release_type"` // Production or Testing
	ReleaseDate   CustomTime `json:"release_date"`
	MajorVersion  string     `json:"major_version"` // e.g., v1.0, v23.1
	Major         int        `json:"major"`
	Minor         int        `json:"minor"`
	Patch         int        `json:"path"`
	BetaRc        string     `json:"beta_rc"`
	BetaRcVersion int        `json:"beta_rc_version"`
}

func NewFromRemote(remote RemoteRelease) *Release {
	v := remote.Version()
	return &Release{
		Name:          remote.Name,
		Withdrawn:     remote.Withdrawn,
		CloudOnly:     remote.CloudOnly,
		ReleaseType:   remote.ReleaseType,
		ReleaseDate:   remote.ReleaseDate,
		MajorVersion:  remote.MajorVersion,
		Major:         v.Major,
		Minor:         v.Minor,
		Patch:         v.Patch,
		BetaRc:        v.BetaRc,
		BetaRcVersion: v.BetaRcVersion,
	}
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
		releases[i] = *NewFromRemote(rr)
	}
	return releases, nil
}
