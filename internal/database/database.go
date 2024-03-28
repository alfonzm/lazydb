package database

type Database struct {
  string
}

func NewConnection() *Database {
  return &Database{"mysql"}
}

func (d *Database) Query() string {
  return "query results"
}
