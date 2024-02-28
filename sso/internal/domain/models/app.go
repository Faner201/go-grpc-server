package models

// App сама по себе нужна, чтобы знать, в какое преложение мы авторизируемся, так же подписывает ему уникальный jwt

type App struct {
	ID     int
	Name   string
	Secret string
}
