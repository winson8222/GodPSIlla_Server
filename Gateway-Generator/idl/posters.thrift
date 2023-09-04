namespace go postersapi

struct Request {}

struct Response {
    1: required list<string> posterslist;
}

service PosterService {
    Response getuniqueusernames(1: Request req) (api.post="/PosterService/getuniqueusernames");
}

