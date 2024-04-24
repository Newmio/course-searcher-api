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

	return nil
}
