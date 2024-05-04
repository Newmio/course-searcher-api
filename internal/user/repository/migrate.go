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
		avatar text default '',
		date_create timestamp default now()
	)`

	_, err := db.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists user_info(
		id serial primary key,
		id_user int references users(id) on delete cascade,
		name text default '',
		middle_name text default '',
		last_name text default '',
		course_number int default 1,
		group_name text default '',
		proffession text default '',
		proffession_number text default ''
	)`

	_, err = db.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}
