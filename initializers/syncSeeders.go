package initializers

import (
	"github.com/OltLatifi/cv-builder-back/seeders"
)

func SyncSeeders() {
	// because we are in the package initializers we can get the DB variable from the databaseConnection.go file
	var db = DB
	seeders.SeedRoles(db)
	seeders.SeedUserStatuses(db)
	seeders.SeedAdmin(db)
}