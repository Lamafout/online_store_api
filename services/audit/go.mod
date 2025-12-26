module github.com/Lamafout/online-store-api/services/audit

go 1.24

require (
	github.com/Lamafout/online-store-api/libs/contracts v0.0.0

    github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jmoiron/sqlx v1.4.0
	"github.com/joho/godotenv" v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/Lamafout/online-store-api/libs/contracts => ../../libs/contracts
