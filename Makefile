hello:
	echo "HELLO!"

oapi-gen:
	oapi-codegen -old-config-style -generate "chi-server" -package api openapi.yaml > user-service/api/server.go
	oapi-codegen -old-config-style -generate "spec" -package api openapi.yaml > user-service/api/spec.go
	oapi-codegen -old-config-style -generate "types" -package api openapi.yaml > user-service/api/model.go
.PHONY: oapi-gen