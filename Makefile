.PHONY: all frontend clean

all: frontend
	GOOS=darwin GOARCH=amd64 go build -o conjugreater-macos .
	GOOS=windows GOARCH=amd64 go build -o conjugreater-windows.exe .

frontend:
	cd web && npm run build

clean:
	rm -f conjugreater-macos conjugreater-windows.exe
	rm -rf web/build
