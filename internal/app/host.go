package app

var defaultHosts = map[string]string{
	"default": "stackoverflow.com",
	// codereview
	"codereview": "codereview.stackexchange.com",
	"cr":         "codereview.stackexchange.com",
}

func Host(h string) string {
	if h == "" {
		return defaultHosts["default"]
	}

	if hostname, ok := defaultHosts[h]; ok {
		return hostname
	}

	return h
}
