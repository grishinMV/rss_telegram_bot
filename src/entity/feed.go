package entity

type Feed struct {
	Id        int    `db:"id"`
	Link      string `db:"link"`
	LastNew   int64  `db:"last_new"`
	NextParse int64  `db:"next_parse"`
}
