package limiter

type Token struct {
	key              string
	request_limit    int64
	request_interval int64
}

var permittedTokens = []Token{
	{"token-abc", 10, 1},
	{"token-vbb", 5, 10},
	{"token-bvb", 1, 10},
}
