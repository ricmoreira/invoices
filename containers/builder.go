// Inspiration to create dependcy injection came from this post: https://blog.drewolson.org/dependency-injection-in-go/

package containers

import (
	"invoices/config"
	controllers "invoices/controllers/v1"
	"invoices/handlers"
	"invoices/repositories"
	"invoices/server"
	"invoices/services"

	"go.uber.org/dig"
)

// BuildContainer returns a container with all app dependencies built in
func BuildContainer() *dig.Container {
	container := dig.New()

	// config
	err := container.Provide(config.NewConfig)
	if err != nil {
		panic(err)
	}

	// persistance layer
	err = container.Provide(repositories.NewDBCollections)
	if err != nil {
		panic(err)
	}
	err = container.Provide(repositories.NewInvoiceRepository)
	if err != nil {
		panic(err)
	}

	// services
	err = container.Provide(services.NewInvoiceService)
	if err != nil {
		panic(err)
	}
	err = container.Provide(services.NewKafkaConsumer)
	if err != nil {
		panic(err)
	}
	err = container.Provide(services.NewKafkaProducer)
	if err != nil {
		panic(err)
	}

	// controllers
	err = container.Provide(controllers.NewInvoiceController)
	if err != nil {
		panic(err)
	}

	// generic http layer
	err = container.Provide(handlers.NewHttpHandlers)
	if err != nil {
		panic(err)
	}

	// server
	err = container.Provide(server.NewServer)
	if err != nil {
		panic(err)
	}

	return container
}
