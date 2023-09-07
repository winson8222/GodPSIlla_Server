package constants

import (
	"strings"
)

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}

const (
	APIGATEWAY_URL = "localhost:80"
)

const (
	GRPC_PORT = 7999 //To be customised by the user
)
