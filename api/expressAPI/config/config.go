package config

import (
	"apiProject/api/expressAPI/types"
	"fmt"
	"os"
)

var EnvConfig = initConfig()

func initConfig() types.MysqlConfig {
	return types.MysqlConfig{
		ServerPort: getEnv("SERVER_PORT", "8088"),
		DbUser:     getEnv("DB_USER", "root"),
		DbPass:     getEnv("DB_PASS", "root244112311"),
		DbAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		DbName:     getEnv("DB_NAME", "eladmin"),
		JWTSecret:  getEnv("JWT_SECRET", "testjwtsecretadmin"),
		PublicKey:  "-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3JnmWnHSYY+D703RzGS6\nPK/YvslCz/I8K+vr7PZq5ND0OhoLP+pexhBnDtHhDw4vFhFsosNl70YZkOImQiX+\nglF0PJqZld1TAtzSl4UtTGo/TI5fHuxy4q31Q+b/kun8lvfvU7ZJe/uWUbKx6dC7\nw6Pum7k5C9oqaOKXHllp78e1Wu3+qnvU2KnNI8gO2PyGUnLxJMJCIfvwqZPwRh07\nKi0yqsUFWuoIpvyLZePR/R35MGf6FQyavlZx7rNFUpw9prHbzC1nu9iSVi4A7B6l\nikjoEThi+uXIQHptQseIxosVhs2z04SCEobRv1H86vs/ciJKE7Y8nUIcwC0H8a/h\nTQIDAQAB\n-----END RSA PUBLIC KEY-----",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA3JnmWnHSYY+D703RzGS6PK/YvslCz/I8K+vr7PZq5ND0OhoL\nP+pexhBnDtHhDw4vFhFsosNl70YZkOImQiX+glF0PJqZld1TAtzSl4UtTGo/TI5f\nHuxy4q31Q+b/kun8lvfvU7ZJe/uWUbKx6dC7w6Pum7k5C9oqaOKXHllp78e1Wu3+\nqnvU2KnNI8gO2PyGUnLxJMJCIfvwqZPwRh07Ki0yqsUFWuoIpvyLZePR/R35MGf6\nFQyavlZx7rNFUpw9prHbzC1nu9iSVi4A7B6likjoEThi+uXIQHptQseIxosVhs2z\n04SCEobRv1H86vs/ciJKE7Y8nUIcwC0H8a/hTQIDAQABAoIBAF29xEZAwd6VRsJE\n9lb9oqoxK1B/Y7XLwMgFO775Q5kyNeYOtSMW6+kMhU6l3xYvt9CP3PMZR1KzHiAU\nCZ/oV0t3Y4ZxR7yITUMVJSQgAozLRVS51y/j2Dn9JBETsxzx81UPzJJtDrLxyQG0\nhqfN/Ev5eGaSAezIa2cgiojqA/tQvpP42qNvp+cDJGgRyLVmYfXiRco8pGrnW9k/\n0AFe5N9GgvjFvKBtYhxMLdnZLbgF/a0piYcas6ja3Wnb6K1HQXzfE4fy66WCh2Xn\nWoU3FfP9PTqWVT7JSCrSR0ot6ctY/2mKNn/AfoK556PZHJ7BuW4ylwXRPiVAhRrZ\nYPgmYyECgYEA/gmpjM8Ge4YQK33nOZ3qOOiuGQ6Pvzm5i6+VZmq4+6M14AFEeQH/\noQqxHcABlHym9+eWKHOv0nAy61p8qIkg2OQVH0nOdQSvd51N09c0ZNGsOX5XDJlU\n/6ls5VEdI0FPgKpEk53RD+eUCQ6L+xOGM6NFe8eQzqtPYCN4jl8dnWUCgYEA3k4e\nnL8Y16XpJGF5wp+ixsmdwsh9XlWygfILw5onYgGnaSv5FbxjCDSJedGJ0CM6jJXL\niLp7tTwpqe0oBOYW53KkYXn0posEU4lfGe27OQkbvzRXEQE3AGMcjOmz/mNepx8z\n3ZxmG43+NO8Eh6ElKYTKK7dkULMOgPo5bTs/yckCgYBtyHsvUOB6TUt7oCNm8Omh\nwlxKk9JnT2jyBuVHp2NdzACiV6nhqY1xaQ91zd5g7yWxCLIJtUUMalR3BVnN88Tw\nNlEyflDsnSO/S4mwvNX1o+8LwZ+Y4EKtYeifiVhQPg8/iVWtfYw1lVySNWklDiD2\n+94xSeM4jSv2Xh3hWRWRSQKBgD8wy4jY1THvak86GgdVo0qIYvzMSr6282/2optu\nRUWZnMHLixk/nJLnhDCJfHgam3j814c9Iw8IU/uGezqxQM93ifxfU0jH+WnZgZv4\nNKDo0udN9HXT95N3mNUBVXW5P12YBAE5hNjOSvU2//2hs9OSeHlmvvAlhbjp58sB\n7YbpAoGBALJrHmV9VXJ+nZH3GqSvLtscs6fnnOevToeEfukr56lBGylhfq8iUqea\nrzcJSDpUq+Ead38uBba5iapvFLeTSkFBYFNY4yHnPW4Uv/VuuIGIH9NVI0PTOK7Y\n5yeZ8oGYJBbjhTbrDF0YpMdcZOsvhjwp7D44TpcV9DDmsdkcLccY\n-----END RSA PRIVATE KEY-----",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
