package repository

import "github.com/Newmio/newm_helper"

func (r *psqlFileRepo) initTables() error {
	str := `create table if not exists report_files(
		id serial primary key,
		name text not null,
		type text not null,
		directory text not null,
		date_create timestamp default now()
	)`

	_, err := r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists education_files(
		id serial primary key,
		name text not null,
		type text not null,
		directory text not null,
		date_create timestamp default now()
	)`

	_, err = r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists report_file_user(
		id serial primary key,
		id_report_file int references report_files(id) on delete cascade,
		id_user int references users(id) on delete cascade,
		id_course int references courses(id) on delete cascade
	)`

	_, err = r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	str = `create table if not exists education_file_user(
		id serial primary key,
		id_education_file int references education_files(id) on delete cascade,
		id_user int references users(id) on delete cascade,
		id_course int references courses(id) on delete cascade
	)`

	_, err = r.db.Exec(str)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}
