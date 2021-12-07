package os

import (
	"fmt"
	"os"
	"strconv"
)

func GetStringEnv(key string) string {
	v, b := os.LookupEnv(key)
	fmt.Printf("env key: %s value: %+v bool: %v", key, v, b)
	return os.Getenv(key)
}

func GetIntEnv(key string) int {
	val, set := os.LookupEnv(key)
	if !set {
		return 0
	}
	ival, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return ival
}

func Exit(code int) {
	os.Exit(code)
}
