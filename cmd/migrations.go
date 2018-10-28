package cmd

import (
	"github.com/go-pg/pg"
	"log"
)

type migration struct {
	description string
	migrate     func(*pg.DB)
	rollback    func(*pg.DB)
}

func getMigrations() map[string]migration {

	migrations := map[string]migration{
		"1": atai_1,
	}
	return migrations
}

var atai_1 = migration{
	description: "atai -> remove column 'numero' | 09-07-2018",
	migrate: func(db *pg.DB) {
		// to find the correct SQL, open DB Browser for SQLite, open "vue, journal sql" and copy the SQL
		sql := ``
		if _, err := db.Exec(sql); err != nil {
			log.Println(err)
		}
	},
	rollback: func(db *pg.DB) {},
}
