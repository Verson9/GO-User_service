hello:
	echo "HELLO!"

gen-user-service:
	oapi-codegen -old-config-style -generate "chi-server,skip-fmt" -package api ./user-service/openapi.yaml > user-service/api/server.go
	oapi-codegen -old-config-style -generate "," -package api ./user-service/openapi.yaml > user-service/api/client.go
	oapi-codegen -old-config-style -generate "spec" -package api ./user-service/openapi.yaml > user-service/api/spec.go
	oapi-codegen -old-config-style -generate "types" -package api ./user-service/openapi.yaml > user-service/api/model.go
.PHONY: gen-user-service