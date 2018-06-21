js_source_files := $(shell find frontend -name '*.js' ! -path 'frontend/bundle-es5.js' ! -path 'frontend/bundle.js' ! -path 'frontend/node_modules/*')
go_source_files := $(shell find . -name '*.go' ! -path 'frontend/*')

app_version := "0.0.1"

mac_release_path := "release/mac"
linux_release_path := "release/linux"
win_release_path := "release/win"

mac_app_bundle_path := "$(mac_release_path)/LocalShare.app"
win_executable_path := "$(win_release_path)/LocalShare.exe"
linux_executable_path := "$(linux_release_path)/localshare"

localshare: ${go_source_files} internal/webui/bindata.go
	go build

run: localshare
	./localshare

bundle-mac: localshare
	mkdir -p "$(mac_app_bundle_path)/Contents/MacOS"
	mkdir -p "$(mac_app_bundle_path)/Contents/Resources"

	cp ./localshare $(mac_app_bundle_path)/Contents/MacOS/localshare
	cp ./resources/Info.plist $(mac_app_bundle_path)/Contents/Info.plist
	cp ./resources/Icons/mac/AppIcon.icns $(mac_app_bundle_path)/Contents/Resources/AppIcon.icns

bundle-linux: localshare
	mkdir -p $(linux_release_path)
	cp ./localshare $(linux_executable_path)
	# TODO: Add the icon to a .desktop file to support
	#       icons in most desktop environments on Linux.

bundle-win:
	mkdir -p $(win_release_path)
	go build -ldflags="-H windowsgui" -o $(win_executable_path)
	rcedit $(win_executable_path) --set-icon ./resources/Icons/win/AppIcon.ico

bundle-all: bundle-mac bundle-linux bundle-win

release-mac: bundle-mac
	zip -r -X release/localshare-macOS.zip $(mac_release_path)

release-linux: bundle-linux
	zip -r -X release/localshare-linux.zip $(linux_release_path)

release-win: bundle-win
	zip -r -X release/localshare-windows.zip $(win_release_path)

release-all: release-mac release-linux release-win

internal/webui/bindata.go: frontend/bundle-es5.js frontend/index.html frontend/style.css
	go-bindata -pkg webui -o internal/webui/bindata.go frontend/bundle-es5.js frontend/index.html frontend/style.css

frontend/bundle-es5.js: frontend/bundle.js
	cd frontend && npx babel bundle.js > bundle-es5.js

frontend/bundle.js: frontend/node_modules $(js_source_files)
	cd frontend && npx browserify index.js -o bundle.js

frontend/node_modules: frontend/package.json
	cd frontend && npm install && touch node_modules
