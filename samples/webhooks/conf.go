package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

func init() {
	readDotEnv()
}

func getNotificationsURL(r *http.Request) string {
	notificationsURL := r.URL.Query().Get("notificationsUrl")
	if len(notificationsURL) == 0 {
		notificationsURL = os.Getenv("NOTIFICATIONS_URL")
	}

	return notificationsURL
}

func resolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), relativePath)
}

func readDotEnv() {
	envFilePath := resolveCnfgPath(".env")
	envFile, err := os.Open(envFilePath)
	if err != nil {
		return
	}
	defer func() { _ = envFile.Close() }()

	byteValue, _ := ioutil.ReadAll(envFile)
	keyVals := strings.Split(fmt.Sprintf("%s", byteValue), "\n")
	for _, keyVal := range keyVals {
		kv := strings.SplitN(keyVal, "=", 2)
		if len(kv) == 2 {
			_ = os.Setenv(kv[0], kv[1])
		}
	}
}
