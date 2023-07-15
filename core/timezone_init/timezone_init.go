package timezone_init

import (
	"log"
	"time"
)

func init() {
	// import as first package (1st on import list) on main to set default timezone on process
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatalf("Failed to read configuration because %s", err)
	}
	time.Local = loc
}
