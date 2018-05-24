js_source_files := $(shell find frontend/ -name '*.js' ! -path 'frontend/bundle.js' ! -path 'frontend/node_modules/*')
go_source_files := $(shell find . -name '*.go' ! -path 'frontend/*')

localshare: ${go_source_files}
	go build

internal/webui/bindata.go: frontend/bundle.js frontend/index.html frontend/style.css
	go-bindata -pkg webui -o internal/webui/bindata.go frontend/bundle.js frontend/index.html frontend/style.css

frontend/bundle.js: frontend/node_modules $(js_source_files)
	cd frontend && npx browserify index.js -o bundle.js

frontend/node_modules: frontend/package.json
	cd frontend && npm install && touch node_modules
