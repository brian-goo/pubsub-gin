LOCAL_PORT = localhost:5000
DEPLOY_PORT = :5000
LOCAL_REDIS_ADDR = localhost:6379

AUTH0_ISS_DEV = 
AUTH0_AUD_DEV = 

TOKEN := test
TOKEN_AUTH0_DEV := 

CURL_GET_LOCAL := curl --header "Authorization: Bearer $(TOKEN)" $(LOCAL_PORT)
CURL_POST_LOCAL := curl --header "Content-Type: application/json" --header "Authorization: Bearer $(TOKEN)" --request POST $(LOCAL_PORT)
CURL_GET_LOCAL_AUTH0 := curl --header "Authorization: Bearer $(TOKEN_AUTH0_DEV)" $(LOCAL_PORT)
CURL_POST_LOCAL_AUTH0 := curl --header "Content-Type: application/json" --header "Authorization: Bearer $(TOKEN_AUTH0_DEV)" --request POST $(LOCAL_PORT)

run:
	PORT=$(LOCAL_PORT) API_KEY=$(TOKEN) ISS=$(AUTH0_ISS_DEV) AUD=$(AUTH0_AUD_DEV) APP_ENV=DEPLOYED go run .

dev: export PORT=$(LOCAL_PORT)
dev: export APP_ENV=LOCAL
dev: export REDIS_ADDR=$(LOCAL_REDIS_ADDR)
dev: export API_KEY=$(TOKEN)
dev: export ISS=$(AUTH0_ISS_DEV)
dev: export AUD=$(AUTH0_AUD_DEV)
dev:
	npx nodemon --exec go run . --signal SIGKILL


g-ping:
	$(CURL_GET_LOCAL)/ping
g-pingA:
	$(CURL_GET_LOCAL_AUTH0)/ping

p-ping:
	$(CURL_POST_LOCAL)/ping --data @_test/ping.json