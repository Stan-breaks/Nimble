run:
	go run .

compile: tailwindcss templ run 

watch:
	find . -name "*.templ" | entr -r make compile

tailwindcss:
	bun run tailwindcss --config tailwind.config.js -i ./src/tailwind.css -o ./public/tailwind.css

templ:
	templ generate ./views/*

.PHONY: run compile watch tailwindcss templ

