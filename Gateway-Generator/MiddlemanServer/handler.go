package main

import (
	"context"
	"encoding/json"
	"fmt"
	"server/constants"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"log"
	"server/model"
	"server/mychacha20"
)

type server struct {
	model.UnimplementedMiddlemanServer
}

// PSI implements PSI on MiddleMan
func (s *server) PSI(ctx context.Context, in *model.Request) (*model.Response, error) {
	secretKey, err := mychacha20.GenerateChaCha20Key()
	if err != nil {
		fmt.Println("mychacha20.GenerateChaCha20Key failed:", err)
		return &model.Response{}, status.Error(codes.Internal, "Internal Server Error")
	}
	secretNonce, err := mychacha20.GenerateChaCha20Nonce()
	if err != nil {
		fmt.Println("mychacha20.GenerateChaCha20Nonce failed:", err)
		return &model.Response{}, status.Error(codes.Internal, "Internal Server Error")
	}

	//Asynchronously make HTTP requests to API Gateway asking for intersection of microsvc data
	//If invalid params, immediately terminate entire operation
	log.Printf("Received SvcInfo: %v", in.GetSvcInfo())
	serviceInfoList, err := extractRequestSvcInfo(in)
	if err != nil {
		fmt.Println("extractRequestSvcInfo failed:", err.Error())
		return &model.Response{}, status.Error(codes.InvalidArgument, "Invalid Service Info")
	}

	APIGatewayURL := constants.APIGATEWAY_URL
	response, err := makePSIReqToAPIGateway(ctx, "http://"+APIGatewayURL, serviceInfoList)
	if response == nil {
		fmt.Println("makePSIReqToAPIGateway failed: " + "request to API gateway returned nil")
		return nil, status.Error(codes.Internal, "Internal Server Error")
	} else if response.StatusCode() == consts.StatusInternalServerError || response == nil {
		fmt.Println("makePSIReqToAPIGateway failed: " + string(response.Body()))
		return nil, status.Error(codes.Internal, "Internal Server Error")
	} else if response.StatusCode() == consts.StatusBadRequest {
		fmt.Println("makePSIReqToAPIGateway failed: " + string(response.Body()))
		return nil, status.Error(codes.InvalidArgument, "Invalid Service Info")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	microsvcIntersection, err := extractResponseAsListBytes(response)
	if err != nil {
		fmt.Println("extractResponseAsListBytes failed:", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	//Encrypt data from client
	err = mychacha20.Encrypt(secretKey, secretNonce, in.EncryptedElems)
	if err != nil {
		fmt.Println("mychacha20.Encrypt failed:", err)
		return &model.Response{}, status.Error(codes.Internal, "Internal Server Error")
	}

	err = mychacha20.Encrypt(secretKey, secretNonce, microsvcIntersection)
	if err != nil {
		fmt.Println("mychacha20.Encrypt failed:", err)
		return &model.Response{}, status.Error(codes.Internal, "Internal Server Error")
	}

	return &model.Response{DoubleEncryptedElems: in.EncryptedElems, EncryptedServerElems: microsvcIntersection}, nil
}

// HELPER METHODS
// Used to write code for overarching steps of handlers.
func extractRequestSvcInfo(in *model.Request) ([][]string, error) {
	if in.GetSvcInfo() == nil || len(in.SvcInfo) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid Service Info")
	}
	serviceInfoList := make([][]string, len(in.SvcInfo))
	for i, info := range in.SvcInfo {
		serviceInfoList[i] = []string{info.ServiceName, info.MethodName}
	}

	return serviceInfoList, nil
}

// Response body should be a json list containing intersection
func makePSIReqToAPIGateway(ctx context.Context, serverURL string, serviceInfoList [][]string) (*protocol.Response, error) {
	// Convert serviceInfoList to JSON
	jsonData, err := json.Marshal(serviceInfoList)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	// Create an Hertz HTTP client to communicate with API Gateway
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	req := &protocol.Request{}
	req.SetBody(jsonData)
	res := &protocol.Response{}
	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI(serverURL + "/PSI")
	err = c.Do(ctx, req, res)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Service Info")
	}

	return res, nil
}

func extractResponseAsListBytes(resp *protocol.Response) ([][]byte, error) {
	// Unmarshal the JSON response into a []string
	var stringList []string
	fmt.Println(resp)
	if err := json.Unmarshal(resp.BodyBytes(), &stringList); err != nil {
		return nil, status.Error(codes.Internal, "Unable to marshal JSON to string list")
	}

	// Convert each string to []byte and append it to byteList
	var byteList [][]byte
	for _, str := range stringList {
		byteList = append(byteList, []byte(str))
	}

	return byteList, nil
}
