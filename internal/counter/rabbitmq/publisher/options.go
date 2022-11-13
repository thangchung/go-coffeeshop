package publisher

type Option func(*Publisher)

func ExchangeName(exchangeName string) Option {
	return func(p *Publisher) {
		p.exchangeName = exchangeName
	}
}

func BindingKey(bindingKey string) Option {
	return func(p *Publisher) {
		p.bindingKey = bindingKey
	}
}

func MessageTypeName(messageTypeName string) Option {
	return func(p *Publisher) {
		p.messageTypeName = messageTypeName
	}
}
