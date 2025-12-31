.PHONY: dev build deploy test commit clear
SCREEN_NAME := "judge"
APP_NAME := "main.py"
buildWatch:
	# pyinstaller --onefile --distpath ./dist watcher.py
	nuitka --standalone --onefile --output-dir=target/ --output-filename=watcher watcher.py
	# sudo mv dist/watcher /usr/local/bin/

dev: stop
	@./target/watcher

stop:
	@-screen -S $(SCREEN_NAME) -X quit || true
	@echo "Stopped screen session: $(SCREEN_NAME)"
start:
	@screen -dmS $(SCREEN_NAME) sh -c 'python3 main.py' 
	@echo "Screen session $(SCREEN_NAME) started with ./$(APP_NAME)"

deploy:stop start