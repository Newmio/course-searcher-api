package user

func (r *psqlRepo) initTables() {
	str := `create table if not exists users(
		id serial primary key,
		login text not null unique,
		password text not null,
		email text default '',
		phone text default '',
		role text not null,
	)`
	_, err := r.db.Exec(str)
	if err != nil {
		panic(err)
	}

	str = `create table if not exists person(
		id serial primary key,
		first_name text not null,
		middle_name text,
		last_name text not null,
		id_user int not null,
	)`
	_, err = r.db.Exec(str)
	if err != nil {
		panic(err)
	}
}
