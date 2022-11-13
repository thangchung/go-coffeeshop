package consumer

type Option func(*Consumer)

func ExchangeName(exchangeName string) Option {
	return func(p *Consumer) {
		p.exchangeName = exchangeName
	}
}

func QueueName(queueName string) Option {
	return func(p *Consumer) {
		p.queueName = queueName
	}
}

func BindingKey(bindingKey string) Option {
	return func(p *Consumer) {
		p.bindingKey = bindingKey
	}
}

func ConsumerTag(consumerTag string) Option {
	return func(p *Consumer) {
		p.consumerTag = consumerTag
	}
}

func MessageTypeName(messageTypeName string) Option {
	return func(p *Consumer) {
		p.messageTypeName = messageTypeName
	}
}

func WorkerPoolSize(workerPoolSize int) Option {
	return func(p *Consumer) {
		p.workerPoolSize = workerPoolSize
	}
}
