package model

type Stat struct {
	Type  string
	Value float64
}

var (
	RandomValue   string = "RandomValue"
	Alloc         string = "Alloc"
	TotalAlloc    string = "TotalAlloc"
	Sys           string = "Sys"
	Lookups       string = "Lookups"
	Mallocs       string = "Mallocs"
	Frees         string = "Frees"
	HeapAlloc     string = "HeapAlloc"
	HeapSys       string = "HeapSys"
	HeapIdle      string = "HeapIdle"
	HeapInuse     string = "HeapInuse"
	HeapReleased  string = "HeapReleased"
	HeapObjects   string = "HeapObjects"
	StackInuse    string = "StackInuse"
	StackSys      string = "StackSys"
	MSpanInuse    string = "MSpanInuse"
	MSpanSys      string = "MSpanSys"
	BuckHashSys   string = "BuckHashSys"
	GCSys         string = "GCSys"
	OtherSys      string = "OtherSys"
	NextGC        string = "NextGC"
	MCacheInuse   string = "MCacheInuse"
	MCacheSys     string = "MCacheSys"
	LastGC        string = "LastGC"
	PauseTotalNs  string = "PauseTotalNs"
	NumGC         string = "NumGC"
	NumForcedGC   string = "NumForcedGC"
	GCCPUFraction string = "GCCPUFraction"
	PollCount     string = "PollCount"
)

type MemStats struct {
	PollCount int
	Mapa      map[string]*Stat
}
