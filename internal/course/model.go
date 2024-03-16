package course

type Course struct {
	Id          int    `db:"id" json:"id" xml:"id"`
	Name        string `db:"name" json:"name" xml:"name"`
	Description string `db:"description" json:"description" xml:"description"`
	Author      string `db:"author" json:"author" xml:"author"`
	Money       int    `db:"money" json:"money" xml:"money"`
}
