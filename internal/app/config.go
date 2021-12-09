package app

import "fmt"

type Config struct {
	Host string // host to monitor
	Tags *Tags  // tags to see
	Open bool   // auto-open in browser
}

func (c Config) Feed() string {
	return fmt.Sprintf("https://%s/feeds/tag?tagnames=%s&sort=newest", c.Host, c.Tags.String())
}
