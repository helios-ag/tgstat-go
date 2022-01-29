package tgstat

func prepareClient() (*Client, error) {
	cfg := ClientConfig{
		token: "test",
	}

	client, err := NewClient(&cfg, WithEndpoint("http://api-tgstat///"))

	return client, err
}
