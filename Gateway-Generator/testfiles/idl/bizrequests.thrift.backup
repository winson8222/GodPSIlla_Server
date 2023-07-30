namespace go http

struct BizRequest {
    1: i64 vint64
    2: string text
    3: i32 token
    6: list<string> items
    7: i32 version
}

struct BizResponse {
    1: i32 token
    2: string text
    5: i32 http_code
}

service BizService {
    BizResponse BizMethod1(1: BizRequest req);
    
    BizResponse BizMethod2(1: BizRequest req);

    BizResponse BizMethod3(1: BizRequest req);
}

