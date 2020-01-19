// Package generated contains generated code based on go-swagger.
package generated

//go:generate rm -rf models restapi client
//go:generate swagger generate server -f ../swagger.yml --exclude-main --principal app.AuthUser --include-buildapi --strict
//go:generate swagger generate client -f ../swagger.yml
