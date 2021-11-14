migrateup:
	migrate -path src/database/migrations -database "postgresql://postgres:root@localhost:5432/other?sslmode=disable" -verbose up


migratedown:
	migrate  -path src/database/migrations -database "postgresql://postgres:root@localhost:5432/other?sslmode=disable" -verbose down

