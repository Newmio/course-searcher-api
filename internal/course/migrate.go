package course

func (r *courseRepo) initTables() {
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
		link text not null
	)`
	_, err := r.psql.Exec(str)
	if err != nil {
		panic(err)
	}
}