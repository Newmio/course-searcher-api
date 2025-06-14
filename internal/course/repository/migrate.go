package repository

import "github.com/Newmio/newm_helper"

func (r *psqlCourseRepo) initTables() error {
	str := `create table if not exists courses(
		id serial primary key,
		name text not null,
		description text default '',
		language text default '',
		author text default '',
		duration text default '',
		rating text default '',
		platform text not null,
		money text default '',
		link text unique not null,
		icon_link text default '',
		active boolean default true,
		date_create timestamp default now(),
		date_update timestamp default now()
	)`
	_, err := r.psql.Exec(str)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	str = `create table if not exists course_user(
		id serial primary key,
		id_user int references users(id) on delete cascade,
		id_course int references courses(id) on delete cascade,
		name text default '',
		topic text default '',
		date_start timestamp default now(),
		date_end timestamp default now(),
		credits int default 0,
		educ_name text default '',
		UNIQUE (id_user, id_course)
	)`
	_, err = r.psql.Exec(str)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}
