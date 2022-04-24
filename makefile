build-web:
	cd services/web && \
	docker build . -t web

build-game:
	cd services/game && \
	docker build . -t game -f ./services/game