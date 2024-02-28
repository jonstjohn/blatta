package settings

// Summary - create a summary of all the raw settings. Here is what we're interested in:
// - variable name
// - default values over releases
// -

// Iterate over every setting
// Track the first version that it appeared in, if there is a difference between value with different CPU,
// and when the value changes across versions, and if it doesn't exist in the current version
/*
select settings_raw.release_name, settings_raw.cpu,
	settings_raw.memory_bytes, settings_raw.variable, settings_raw.value
from
	settings_raw INNER JOIN releases ON settings_raw.release_name = releases.name
order
	by settings_raw.variable, releases.major, releases.minor, releases.patch,
	releases.beta_rc, releases.beta_rc_version, settings_raw.cpu,
	settings_raw.memory_bytes
limit 100;
*/

type Summarizer struct {
	SqlExecutor SqlExecutor
}

type Change struct {
	Release string
	From    string
	To      string
}

type Summary struct {
	Variable           string
	Value              string
	Type               string
	Public             bool
	Description        string
	DefaultValue       string
	Origin             string
	Key                string
	FirstReleases      []string
	LastReleases       []string
	HostDependent      bool
	ValueChanges       []Change
	DescriptionChanges []Change
}

func NewSummarizer(sqlExecutor SqlExecutor) *Summarizer {
	return &Summarizer{SqlExecutor: sqlExecutor}
}

func NewSummaryFromRaw(rs RawSetting) *Summary {
	return &Summary{
		Variable:           rs.Variable,
		Value:              rs.Value,
		Type:               rs.Type,
		Public:             rs.Public,
		Description:        rs.Description,
		DefaultValue:       rs.DefaultValue,
		Origin:             rs.Origin,
		Key:                rs.Key,
		FirstReleases:      []string{rs.ReleaseName},
		LastReleases:       make([]string, 0),
		ValueChanges:       make([]Change, 0),
		DescriptionChanges: make([]Change, 0),
	}
}

func (sum *Summarizer) Summarize() error {
	rawSettings, err := sum.SqlExecutor.GetSettingsRawOrderedByVersion()
	if err != nil {
		return err
	}
	summary := Summary{}
	for _, rs := range rawSettings {
		// New summary if variable has changed
		if summary.Variable != rs.Variable {
			// Persist summary and create a new one
			sum.SqlExecutor.UpsertSummary(summary)
			summary = *NewSummaryFromRaw(rs)
		} else { // otherwise process
			// Value change
			if summary.Value != rs.Value {
				summary.ValueChanges = append(summary.ValueChanges, Change{Release: rs.ReleaseName, From: summary.Value, To: rs.Value})
				summary.Value = rs.Value
			}

			// Description change
			if summary.Description != rs.Description {
				summary.DescriptionChanges = append(summary.DescriptionChanges, Change{Release: rs.ReleaseName, From: summary.Description, To: rs.Description})
				summary.Description = rs.Description

			}
		}
	}

	// Persist last summary
	sum.SqlExecutor.UpsertSummary(summary)
}
