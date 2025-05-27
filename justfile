run:
	go run .

watch:
	bash -c 'find . \( -name "*.templ" -o -name "*.sql" \) | entr -r just compile'

tailwindcss:
	bun run tailwindcss --config tailwind.config.js -i ./src/tailwind.css -o ./public/tailwind.css

templ:
	go tool templ generate ./views/*

sqlc:
  sqlc generate

compile: tailwindcss templ sqlc run 

