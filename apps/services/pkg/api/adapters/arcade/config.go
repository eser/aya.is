package arcade

type Config struct {
	URL    string `conf:"url"    default:"https://api.arcade.dev/v1/tools/execute"`
	APIKey string `conf:"apikey"`
}
