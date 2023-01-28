compose-up:
	docker-compose -f ./deployments/docker-compose.yml up -d --build --force-recreate

watch:
	air

keys:
	mkdir -p ./rsa && \
	openssl genrsa -out ./rsa/access_token_private_key.pem 2048 && \
    openssl rsa -in ./rsa/access_token_private_key.pem -outform PEM -pubout -out ./rsa/access_token_public_key.pem && \
    openssl genrsa -out ./rsa/refresh_token_private_key.pem 2048 && \
    openssl rsa -in ./rsa/refresh_token_private_key.pem -outform PEM -pubout -out ./rsa/refresh_token_public_key.pem