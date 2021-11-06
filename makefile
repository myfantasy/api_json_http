
res:
	cd examples/service && go run ./ -s "settings/settings.json"

rec:
	cd examples/client && go run ./

recp:
	cd examples/client_conn_pull && go run ./