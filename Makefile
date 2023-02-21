dev_install:
	@go build -o main && ./main install ffmpeg

dev_uninstall:
	@go build -o main && ./main uninstall ffmpeg -d

.PHONY: dev_install dev_uninstall
