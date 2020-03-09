package notification

// Option for building rabbit client.
type Option func(*client)

// Exchange set options for connect rabbit.
func Exchange(exchange string) Option {
	return func(client *client) {
		client.exchange = exchange
	}
}

// RoutingKey set options for connect rabbit.
func RoutingKey(key string) Option {
	return func(client *client) {
		client.key = key
	}
}

// Mandatory set options for connect rabbit.
func Mandatory() Option {
	return func(client *client) {
		client.mandatory = true
	}
}

// Immediate set options for connect rabbit.
func Immediate() Option {
	return func(client *client) {
		client.immediate = true
	}
}

// AppID set options for connect rabbit.
func AppID(appID string) Option {
	return func(client *client) {
		client.appID = appID
	}
}

// GeneratorID set options for connect rabbit.
func GeneratorID(generator func() (string, error)) Option {
	return func(client *client) {
		client.generatorID = generator
	}
}
