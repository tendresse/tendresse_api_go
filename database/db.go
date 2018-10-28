package database

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

var db *pg.DB

func GetDB() *pg.DB {
	return db
}

func CloseDB() {
	db.Close()
}

func Init() {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// options.TLSConfig.InsecureSkipVerify = true
	db = pg.Connect(options)
	var check int
	if _, err := db.QueryOne(&check, "SELECT 1;"); err != nil {
		log.Fatal(errors.Wrap(err, "database problem, seems offline"))
	}
	EnableDebug(db)
}

func EnableDebug(db *pg.DB) {
	if os.Getenv("DEBUG") != "" {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
	}
}

func CreateDB() {
	executeQueries(create_tables)
}

func DropDB() {
	executeQueries(drop_tables)
}

func ResetDB() {
	log.Println("reset DB...")
	executeQueries(drop_tables)
	executeQueries(create_tables)
}

func executeQueries(queries []string) error {
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

var create_tables = []string{
	`CREATE TABLE tags (
	id SERIAL CONSTRAINT pk_tag PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,name TEXT NOT NULL UNIQUE
	);

CREATE TABLE achievements (
	id SERIAL CONSTRAINT pk_achievement PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,name TEXT NOT NULL UNIQUE
	,condition INT DEFAULT 10
	,icon TEXT
	,type TEXT
	,xp INT DEFAULT 10
	,tag_id INT REFERENCES tags
	);

CREATE TABLE blogs (
	id SERIAL CONSTRAINT pk_blog PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,description TEXT
	,url TEXT NOT NULL UNIQUE
	);

CREATE TABLE gifs (
	id SERIAL CONSTRAINT pk_gif PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,url TEXT NOT NULL UNIQUE
	,lame_score INT DEFAULT 0
	,blog_id INT REFERENCES blogs ON DELETE SET NULL
	);

CREATE TABLE roles (
	id SERIAL CONSTRAINT pk_role PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,name TEXT NOT NULL
	,description TEXT
	);

CREATE TABLE users (
	id SERIAL CONSTRAINT pk_user PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,email TEXT NOT NULL UNIQUE
	,username TEXT NOT NULL UNIQUE
	,passhash TEXT NOT NULL
	,premium boolean DEFAULT false
	);

CREATE TABLE tokens (
	id SERIAL CONSTRAINT pk_token PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,HASH TEXT NOT NULL UNIQUE
	,user_id INT REFERENCES users
	);

CREATE TABLE tendresses (
	id SERIAL CONSTRAINT pk_tendresse PRIMARY KEY
	,created_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,updated_at TIMESTAMP WITH TIME zone NOT NULL DEFAULT NOW()
	,sender_id INT REFERENCES users ON DELETE CASCADE
	,receiver_id INT REFERENCES users ON DELETE CASCADE
	,gif_id INT REFERENCES gifs ON DELETE CASCADE
	,viewed boolean DEFAULT false
	);

/* associations tables */

CREATE TABLE gifs_tags (
	tag_id INT REFERENCES tags ON DELETE CASCADE
	,gif_id INT REFERENCES gifs ON DELETE CASCADE
	,CONSTRAINT pk_gif_tags PRIMARY KEY (
		tag_id
		,gif_id
		)
	);

CREATE TABLE users_achievements (
	achievement_id INT REFERENCES achievements ON DELETE CASCADE
	,user_id INT REFERENCES users ON DELETE CASCADE
	,score INT DEFAULT 0
	,unlocked boolean DEFAULT false
	,CONSTRAINT pk_user_achievements PRIMARY KEY (
		achievement_id
		,user_id
		)
	);

CREATE TABLE users_friends (
	user_id INT REFERENCES users ON DELETE CASCADE
	,friend_id INT REFERENCES users ON DELETE CASCADE
	,CONSTRAINT pk_user_friends PRIMARY KEY (
		user_id
		,friend_id
		)
	);

CREATE TABLE users_roles (
	role_id INT REFERENCES roles ON DELETE CASCADE
	,user_id INT REFERENCES users ON DELETE CASCADE
	,CONSTRAINT pk_user_roles PRIMARY KEY (
		role_id
		,user_id
		)
	);`,
}

var drop_tables = []string{
	`DROP SCHEMA public CASCADE;`,
	`CREATE SCHEMA public;`,
}
