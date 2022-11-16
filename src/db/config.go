package db

type Config struct {
	host     string
	port     string
	user     string
	password string
	charset  string
	dbName   string
}

func (d Config) Host() string {
	return d.host
}

func (d Config) Port() string {
	return d.port
}

func (d Config) User() string {
	return d.user
}

func (d Config) Password() string {
	return d.password
}

func (d Config) Charset() string {
	return d.charset
}

func (d Config) DbName() string {
	return d.dbName
}
