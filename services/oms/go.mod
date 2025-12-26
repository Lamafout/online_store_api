module github.com/Lamafout/online-store-api/services/oms

go 1.24.5

require (
	github.com/Lamafout/online-store-api/libs/contracts v0.0.0

	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-playground/validator/v10 v10.27.0
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/Lamafout/online-store-api/libs/contracts => ../../libs/contracts
