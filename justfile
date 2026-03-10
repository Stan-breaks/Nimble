run:
	go run .

sqlc:
	sqlc generate

watch:
	bash -c 'find . \( -name "*.go" -o -name "*.sql" \) -not -path "./database/*" | entr -r just compile'

compile: sqlc run
