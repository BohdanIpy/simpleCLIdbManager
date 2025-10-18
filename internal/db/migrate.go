package db

import "database/sql"

func Migrate(db *sql.DB) error {
	queries := []string{
		`create table if not exists logs(
    		id int GENERATED ALWAYS AS IDENTITY,
    		log_time timestamp not null default now(),
    		log_message varchar(255) not null
        );`,
		`create table if not exists users(
			id INT GENERATED ALWAYS AS IDENTITY,
			name varchar(20) not null,
			email varchar(50) not null,
			password varchar(255) not null,
			registered_at timestamp not null default now()
		);`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
