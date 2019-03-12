package common

var (
	StatusPending    = "pending"
	StatusCreating   = "creating"
	StatusStarting   = "starting"
	StatusRunning    = "running"
	StatusStopping   = "stopping"
	StatusStopped    = "stopped"
	StatusMigrating  = "migrating"
	StatusResizing   = "resizing"
	StatusDestroying = "destroying"
	StatusDestroyed  = "destroyed"
)

func checkForward(now, next string) bool {
	switch next {
	case now:
		// met yet.
		return true

	case StatusDestroyed:
		return now == StatusDestroying
	case StatusDestroying:
		return now == StatusStopped || now == StatusDestroyed

	case StatusStopped:
		return now == StatusStopping || now == StatusMigrating
	case StatusStopping:
		return now == StatusRunning || now == StatusStopped

	case StatusMigrating:
		return now == StatusResizing || now == StatusStopped

	case StatusResizing:
		return now == StatusStopped

	case StatusRunning:
		return now == StatusStarting

	case StatusStarting:
		return now == StatusStopped || now == StatusCreating || now == StatusRunning

	case StatusCreating:
		return now == StatusPending

	case StatusPending:
		return now == ""

	default:
		return false
	}
}
