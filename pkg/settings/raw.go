package settings

type RawSetting struct {
	ReleaseName  string
	Cpu          int
	MemoryBytes  int64
	Variable     string
	Value        string
	Type         string
	Public       bool
	Description  string
	DefaultValue string
	Origin       string
	Key          string
}

func NewRawSetting(releaseName string, cpu int, memoryBytes int64, cs ClusterSetting) *RawSetting {
	return &RawSetting{
		ReleaseName:  releaseName,
		Cpu:          cpu,
		MemoryBytes:  memoryBytes,
		Variable:     cs.Variable,
		Value:        cs.Value,
		Type:         cs.Type,
		Public:       cs.Public,
		Description:  cs.Description,
		DefaultValue: cs.DefaultValue,
		Origin:       cs.Origin,
		Key:          cs.Key,
	}
}
