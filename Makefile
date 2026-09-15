run-service:
	docker compose up -d # -d для работы в фоне и продолжения выполнения мейка
	docker exec -i postgres psql -U postgres -d documents < migration.sql
	go run -mod=vendor -v document-service/cmd