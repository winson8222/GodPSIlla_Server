namespace go hello

struct HelloReq {
    1: string Name; // Add api annotations for easier parameter binding
}

struct HelloResp {
    1: string RespBody;
}

service HelloService {
    HelloResp HelloMethod(1: HelloReq request) (api.post="/HelloService/HelloMethod");
}
