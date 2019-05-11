package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonString := `{
    "app_id": "1",
    "gateway_rsa_private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA2KaaJp7JeW91WlQCfZeS14US/ot9hIJViutv3JHojdgTx+8A\n8psStKaPl2Ac/MTJ/3mHeopCObmgjw/Au/Ne0PS1rveY0Pcazwnp+R1TDP2H9jag\nc3GJWS6cvHLB/B4uP3LOnPXN8ctwDVsF19b/howVKUKX6RAX7R2VAEyTIZJIEIQE\n0fNvRCWqbVv1RB3LU4cbQmW6nX8dP793fP8s/Lhzcj6vS6UKxLVl5CrCCGIJIBYc\n1mI8RbUYvGqwiONEnEwYvOioAoAlkMIXdFndIjngHe7JYfGW1NtPzHLG5yw8anYT\nD/3du7hJ/kSN0WM6NLa0P/vbR5+mxVdoRzY+kQIDAQABAoIBABaE2qkBADgbGbuV\n19xuENlN/7dtkFJhqbqS1kG6+M0llIjHkvWkoMEePvahCuJLIiPn4ekezdtqLAIy\nxPnERiq6BNh26+9sf+DdSvCV17gV8jfpXawiNQCME8aStw8Zo/z8VfWCpzFmz/LT\nbzwMIOs/TEPJpDiZb6M52+74BqMKfHTY14YOF8Xr4fiaUFpNTViHeOQXKzoG5PF4\nGLlhg7YNgEnjyc578izCoFp/xTjBBHQ7dtu+EnzmXD9QTlz7xUYt4P2TjUEBKy1o\nxSxDpgFL+BKYgRazilkrJ2hesbGCvbxDzcd4ivzpfmvqkN74Lq0vF9voL1JSd6D2\n3l/R9bECgYEA7nbQgsK3ResReUMumJvE4y1sl2D+rt24QHlu+jOJqVkpAo0L4HLj\nvCX0Y8tBfG/hDc5iC12YILCn+EEb9bD2giURg7V+cA+K4IJrbLTbnna/UlvA1PFK\n3kHFosdCk5cRlpAppEBLQEUjlf7mjp6k2Xxy71ozg4KlB3wf3QCrs2MCgYEA6JUk\niXd/lntdjb7V/QdwhVFdp/lzst0ClE4q04RNL8ZjwmSrYOAOGO2ktKOBG8lGT3P6\n54/BASn9TMOXks8gPE3r/pN+21RGvOq2xtHNOrnV5g1RvlqHtwtv2RUxoEoTKPjB\nm6KDeLrPNCuGZ3bzYUUNAys66v3iWM5PK2s1GnsCgYEAjtKyx95/jmzgJlTKj7Sc\nE8SdCX2ajHlXZaZVhZ1gkgFIwrJfrqqhI4tH+I1AR5tqm65EorIH72xe7h1w9ZJr\n0j8JYm1NsShd8WGrnYwlDZ/prxYtRFzQjpWuHXRit6r/acImbq3jZDcEvU3SIRF7\ngpc674iC2f1hgj4hh2hjbikCgYBWfG8ztv34xTMKrHYCOyv6R0FeXwJI9qoo39BJ\nCx9wroMWHD0mLurPFj9y9IHkBTph/SzFwszwU97fFrRcYS0Jf6hL6Cj6AiKzyUvi\nLs30EnqZq0ZEVIG27UfQH3NuuVzalXXZG9trn3vBWJYID1F9UCIAlai5DWOHxl/m\nM11x1QKBgCvi33kFLll6SVIKfkwt1Hja/DGlyq/M/xN4qn/wQwGKYzzIU+73SbQR\ng44GAiYVMJQrjISg/RVd4ClDxZ+A0cpumfpuSJdcT210L4u5FkuTQAmLZ2HOhTzK\nmW9/iR9koFHtTzTKhhYIgSWy9EWkQmcyrOKnEPYqMJjMobDJ1AuG\n-----END RSA PRIVATE KEY-----",
    "gateway_rsa_public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2KaaJp7JeW91WlQCfZeS\n14US/ot9hIJViutv3JHojdgTx+8A8psStKaPl2Ac/MTJ/3mHeopCObmgjw/Au/Ne\n0PS1rveY0Pcazwnp+R1TDP2H9jagc3GJWS6cvHLB/B4uP3LOnPXN8ctwDVsF19b/\nhowVKUKX6RAX7R2VAEyTIZJIEIQE0fNvRCWqbVv1RB3LU4cbQmW6nX8dP793fP8s\n/Lhzcj6vS6UKxLVl5CrCCGIJIBYc1mI8RbUYvGqwiONEnEwYvOioAoAlkMIXdFnd\nIjngHe7JYfGW1NtPzHLG5yw8anYTD/3du7hJ/kSN0WM6NLa0P/vbR5+mxVdoRzY+\nkQIDAQAB\n-----END PUBLIC KEY-----",
    "md5_key": "askldjlk",
    "merchant_rsa_public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2KaaJp7JeW91WlQCfZeS\n14US/ot9hIJViutv3JHojdgTx+8A8psStKaPl2Ac/MTJ/3mHeopCObmgjw/Au/Ne\n0PS1rveY0Pcazwnp+R1TDP2H9jagc3GJWS6cvHLB/B4uP3LOnPXN8ctwDVsF19b/\nhowVKUKX6RAX7R2VAEyTIZJIEIQE0fNvRCWqbVv1RB3LU4cbQmW6nX8dP793fP8s\n/Lhzcj6vS6UKxLVl5CrCCGIJIBYc1mI8RbUYvGqwiONEnEwYvOioAoAlkMIXdFnd\nIjngHe7JYfGW1NtPzHLG5yw8anYTD/3du7hJ/kSN0WM6NLa0P/vbR5+mxVdoRzY+\nkQIDAQAB\n-----END PUBLIC KEY-----"
}`
	m := make(map[string]interface{})
	e := json.Unmarshal([]byte(jsonString), &m)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(m)
		bytes, _ := json.Marshal(m)
		fmt.Println("json: ", string(bytes))
	}
}
