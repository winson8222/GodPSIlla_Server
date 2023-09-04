package constants

	 import (
		"strings"
	 )
 
	func ToConstant(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
	}


	const (
		ETCD_URL = "0.0.0.0:20000" //connects to a single etcd instance in the cluster
		LOAD_BALANCING = "ROUND_ROBIN"
	)
	const (
		FILEPATH_TO_POSTERSERVICE  = "../idl/posters.thrift"
		POSTERSERVICE_NAME         = "PosterService" //name registered with svc registry as rpcendpoint
	)
	
	const (
		FILEPATH_TO_VIEWERSERVICE  = "../idl/viewers.thrift"
		VIEWERSERVICE_NAME         = "ViewerService" //name registered with svc registry as rpcendpoint
	)
	