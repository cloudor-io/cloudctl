package api

type JobStat struct {
	Booting         int32
	Running         int32
	AssignInstances int
}

type State struct {
	Idling      int
	Busying     int
	Reserved    int
	Terminating int
	Booting     int
	Max         int
	Tag         string
	JobState    *JobStat
}

type SchedulerStatus struct {
	Waiting int
	Status  map[string]State
}
