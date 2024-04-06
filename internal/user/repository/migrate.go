package repository

import "github.com/Newmio/newm_helper"

func (db *psqlUserRepo) initTables() error {
	str := `create table if not exists users(
		id serial primary key,
		login text unique not null,
		password text not null,
		email text unique not null,
		phone text default '',
		role text not null,
		date_create timestamp default now()
	)`

	_, err := db.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists persons(
		id serial primary key,
		first_name text default '',
		middle_name text default '',
		last_name text default '',
		city text default '',
		street text default '',
		university_role text not null,
		id_user int not null,
		date_create timestamp default now()
	)`

	_, err = db.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}
