package migrate

type Config struct {
	Postgres struct {
		URL    string
		Schema string
	}
}
