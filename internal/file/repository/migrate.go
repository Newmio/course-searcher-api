package repository

import "github.com/Newmio/newm_helper"

func (r *psqlFileRepo) initTables() error {
	str := `create table if not exists report_files(
		id serial primary key,
		name text not null,
		type text not null,
		directory text not null,
		date_create timestamp default now(),
	)`

	_, err := r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists educetion_files(
		id serial primary key,
		name text not null,
		type text not null,
		directory text not null,
		date_create timestamp default now(),
	)`

	_, err = r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists report_file_user(
		id serial primary key,
		foreign key(id_report_file) references report_files(id) on delete cascade not null,
		foreign key(id_user) references users(id) on delete cascade not null
	)`

	_, err = r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists education_file_user(
		id serial primary key,
		foreign key(id_education_file) references educetion_files(id) on delete cascade not null,
		foreign key(id_user) references users(id) on delete cascade not null
	)`

	return nil
}
