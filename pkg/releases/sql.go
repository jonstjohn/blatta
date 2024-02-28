package releases

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type ReleasesRow struct {
	Name          string
	Withdrawn     bool
	CloudOnly     bool
	ReleaseType   string
	ReleaseDate   time.Time
	MajorVersion  string
	Major         int
	Minor         int
	Patch         int
	BetaRc        string
	BetaRcVersion int
}

type SqlExecutor struct {
	Pool *pgxpool.Pool
}

const CREATE_TABLE = `
CREATE TABLE IF NOT EXISTS releases (
	name STRING NOT NULL PRIMARY KEY,
	withdrawn BOOL NOT NULL,
	cloud_only BOOL NOT NULL,
	release_type STRING NOT NULL,
	release_date TIMESTAMP NOT NULL,
	major_version STRING NOT NULL,
	major INT NOT NULL,
	minor INT NOT NULL,
	patch INT NOT NULL,
	beta_rc STRING,
	beta_rc_version INT,
	INDEX (release_date),
	INDEX (major, minor, patch)
)	
`

const SelectAllReleasesSql = `
SELECT
	name,
	withdrawn,
	cloud_only,
	release_type,
	release_date,
	major_version,
	major,
	minor,
	patch,
	beta_rc,
	beta_rc_version
FROM
	releases
ORDER BY major ASC, minor ASC, patch ASC, beta_rc ASC, beta_rc_version ASC
`

const UPSERT = `
UPSERT INTO releases (
	name, withdrawn, cloud_only,
	release_type, release_date, major_version, 
	major, minor, patch, 
	beta_rc, beta_rc_version)
VALUES (
	$1, $2, $3,
	$4, $5, $6,
	$7, $8, $9,
	$10, $11)
`

const MostRecentSql = `
SELECT name, withdrawn, cloud_only,
	release_type, release_date, major_version,
	major, minor, patch,
	beta_rc, beta_rc_version
FROM releases
ORDER BY major DESC, minor DESC, patch DESC
LIMIT 1
`

func NewSqlExecutor(pool *pgxpool.Pool) *SqlExecutor {
	return &SqlExecutor{
		Pool: pool,
	}
}

func (s *SqlExecutor) CreateTable() error {
	_, err := s.Pool.Exec(context.Background(), CREATE_TABLE)
	return err
}

func (s *SqlExecutor) UpsertRelease(r Release) error {
	_, err := s.Pool.Exec(context.Background(), UPSERT,
		r.Name, r.Withdrawn, r.CloudOnly,
		r.ReleaseType, r.ReleaseDate, r.MajorVersion,
		r.Major, r.Minor, r.Patch,
		r.BetaRc, r.BetaRcVersion,
	)
	return err
}

func (s *SqlExecutor) GetRecentReleaseNames(cnt int) ([]string, error) {
	rows, err := s.Pool.Query(context.Background(),
		"SELECT name FROM releases WHERE withdrawn = false AND cloud_only = false ORDER BY release_date DESC LIMIT $1 ", cnt)

	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	return names, nil
}

func (s *SqlExecutor) GetAllReleasesRows() ([]ReleasesRow, error) {
	rows, err := s.Pool.Query(context.Background(), SelectAllReleasesSql)
	if err != nil {
		return nil, err
	}
	rrs := make([]ReleasesRow, 0)
	for rows.Next() {
		var name string
		var withdrawn bool
		var cloudOnly bool
		var releaseType string
		var releaseDate time.Time
		var majorVersion string
		var major int
		var minor int
		var patch int
		var betaRc string
		var betaRcVersion int
		rows.Scan(&name, &withdrawn, &cloudOnly, &releaseType, &releaseDate, &majorVersion, &major,
			&minor, &patch, &betaRc, &betaRcVersion)
		rrs = append(rrs, ReleasesRow{Name: name, Withdrawn: withdrawn, CloudOnly: cloudOnly,
			ReleaseType: releaseType, ReleaseDate: releaseDate, MajorVersion: majorVersion,
			Major: major, Minor: minor, Patch: patch, BetaRc: betaRc, BetaRcVersion: betaRcVersion,
		})
	}
	return rrs, nil
}
