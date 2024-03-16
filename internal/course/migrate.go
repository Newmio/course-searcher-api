package course

func (r *courseRepo) initTables() {
	str := `create table if not exists courses(
		id serial primary key,
		name text not null,
		description text default '',
		author text not null,
		money int default 0,
	)`
	_, err := r.psql.Exec(str)
	if err != nil {
		panic(err)
	}
}