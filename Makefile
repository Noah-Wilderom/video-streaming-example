
run: build
	@./bin/app

build: templ
	@npx tailwindcss -i resources/assets/css/app.css -o ./public/assets/styles.css
	@npx esbuild resources/assets/js/app.js --bundle --outdir=public/assets
	@go build -o ./bin/app ./main.go

templ:
	@templ generate