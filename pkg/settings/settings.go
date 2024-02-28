package settings

import (
	"blatta/pkg/db"
	"blatta/pkg/host"
	"blatta/pkg/releases"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Persister struct {
	SqlExecutor *SqlExecutor
}

func ClusterSettingsFromRelease(release string) ([]ClusterSetting, error) {
	t, err := testserver.NewTestServer(
		testserver.CustomVersionOpt(release))
	if err != nil {
		return nil, err
	}
	pool, err := db.NewPoolFromUrl(t.PGURL().String())
	if err != nil {
		return nil, err
	}
	return GetLocalClusterSettings(pool)
}

// SaveClusterSettingsForVersion saves all of the cluster settings for a specific CRDB version, but only
// if the combination of release, cpu and memory has not been previously run - otherwise it bails early.
func SaveClusterSettingsForVersion(release string, url string) error {

	pool, err := db.NewPoolFromUrl(url)
	if err != nil {
		return err
	}
	p := Persister{SqlExecutor: NewSqlExecutor(pool)}

	// Get host memory and CPU
	cpu := host.GetCpu()
	memoryBytes, err := host.GetMemory()
	if err != nil {
		return err
	}

	rs := make([]string, 0)
	if strings.HasPrefix(release, "recent-") {
		rp := releases.NewPersister(pool)
		cntStr := strings.Replace(release, "recent-", "", 1)
		cnt, err := strconv.Atoi(cntStr)
		if err != nil {
			return err
		}
		rs, err = rp.SqlExecutor.GetRecentReleaseNames(cnt)
		if err != nil {
			return err
		}
	} else {
		rs = append(rs, release)
	}

	// Iterate over releases
	for _, r := range rs {

		// Check to see if save run already exists, if it does, bail early - we've already captured the settings
		exists, err := p.SqlExecutor.SaveRunExists(r, cpu, memoryBytes)
		if err != nil {
			return err
		}
		if exists {
			logrus.Info(fmt.Sprintf("Save run already exists for '%s' with cpu/memory %d/%d", r, cpu, memoryBytes))
			continue
		}

		// Get the cluster settings for this release
		settings, err := ClusterSettingsFromRelease(r)
		if err != nil {
			return err
		}
		rawSettings := make([]RawSetting, len(settings))

		// Convert the cluster settings into raw settings to be saved
		for i, s := range settings {
			rawSettings[i] = *NewRawSetting(r, cpu, memoryBytes, s)
		}

		err = p.SaveRawSettings(rawSettings)
		if err != nil {
			return err
		}

		// Save the save run so we don't have to re-run later
		err = p.SqlExecutor.UpsertSaveRun(r, cpu, memoryBytes)
		if err != nil {
			return err
		}
	}

	return nil

}

func NewPersister(pool *pgxpool.Pool) *Persister {
	return &Persister{SqlExecutor: NewSqlExecutor(pool)}
}

func (p *Persister) Init() error {
	return p.SqlExecutor.CreateRawTable()
}

func (p *Persister) SaveRawSettings(settings []RawSetting) error {
	for _, s := range settings {
		err := p.SqlExecutor.UpsertRawSetting(s)
		if err != nil {
			return err
		}
	}
	return nil
}
