package functions

import "github.com/streadway/amqp"

//DeclareExchange function to declare exchange to connect with queue
func DeclareExchange(cha *amqp.Channel, exchangeName string) error {
	return cha.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // args
	)
}

//DeclareQueue function to declare queue for store data
func DeclareQueue(cha *amqp.Channel, queueName string) (amqp.Queue, error) {
	return cha.QueueDeclare(
		queueName, //name
		false,     //durable
		false,     //delete when not used
		false,     //exclusive
		false,     //no-wait
		nil,       //args
	)
}

//BindQueue function to bind queue to connect queue with exchange
func BindQueue(cha *amqp.Channel, exchangeName, queueName string) error {
	return cha.QueueBind(
		queueName,    //name
		"",           //rounting key
		exchangeName, //exchange name
		false,        //no-wait
		nil,          //args
	)
}

//ConsumeData function to cunsume data in queue to use
func ConsumeData(cha *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return cha.Consume(
		queueName, //name
		"",        //consumer
		true,      //auto-ack
		false,     //exclusive
		false,     //no-local
		false,     //no-wait
		nil,       //args
	)
}

//PublishData function to publish data to exchange
func PublishData(cha *amqp.Channel, exchangeName string, data []byte) error {
	return cha.Publish(
		"keyword", //exchange name
		"",        //rounting key
		false,     //mandatory
		false,     //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
}
