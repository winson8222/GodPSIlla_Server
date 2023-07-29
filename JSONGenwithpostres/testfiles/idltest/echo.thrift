namespace go echoapi

struct Request {
	1: required string message;
}

struct Response {
	1: required string message;
}

service Echo {
    Response Echo(1: required Request req) (api.post="/Echo/Echo");
}
