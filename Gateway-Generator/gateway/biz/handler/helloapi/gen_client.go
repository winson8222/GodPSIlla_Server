
package helloapi

import (
	"context"
	"gateway/constants"
	"log"
	"strings"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	

	"github.com/cloudwego/kitex/pkg/klog"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}


// Creates generic client "[ServiceName]GenericClient"
func HelloServiceGenericClient() genericclient.Client {
	r, err := etcd.NewEtcdResolver([]string{constants.ETCD_URL})
	if err != nil {
		log.Fatal(err)
	}

	path := constants.FILEPATH_TO_HELLOSERVICE
	p, err := generic.NewThriftFileProvider(path)
	if err != nil {
		klog.Fatalf("new thrift file provider failed: %v", err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		klog.Fatalf("new JSON thrift generic failed: %v", err)
	}

	if (constants.LOAD_BALANCING == "ROUND_ROBIN") {
		cli, err := genericclient.NewClient(constants.HELLOSERVICE_NAME, g, client.WithResolver(r),
			client.WithLoadBalancer(loadbalance.NewWeightedBalancer()))
		if err != nil {
			klog.Fatalf("new JSON generic client failed: %v", err)
		}
		return cli
	} else {
		cli, err := genericclient.NewClient(constants.HELLOSERVICE_NAME, g, client.WithResolver(r))
		if err != nil {
			klog.Fatalf("new JSON generic client failed: %v", err)
		}
		return cli
	}
}


func DoHelloMethod(ctx context.Context, cli genericclient.Client, req string) (interface{}, error) {
	fmt.Print(req)
	resp, err := cli.GenericCall(ctx, "HelloMethod", req)

	if err != nil {
		return nil, err
	}
	//OWN CODE ABOVE
	return resp, nil
}

