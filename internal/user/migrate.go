package user

func (r *psqlRepo) initTables() {
	str := `create table if not exists users(
		id serial primary key,
		login text not null unique,
		password text not null,
		email text default '',
		phone text default '',
		role text not null,
		first_name text not null,
		middle_name text,
		last_name text not null
	)`
	_, err := r.psql.Exec(str)
	if err != nil {
		panic(err)
	}
}
