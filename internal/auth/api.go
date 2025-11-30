package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error){
	apiString := headers.Get("Authorization")
	if apiString == ""{
		return "", fmt.Errorf("unable to get Authorization header")
	}
	apiKey, found := strings.CutPrefix(apiString, "ApiKey ")
	if !found {
		return "", fmt.Errorf("authorization format is wrong")
	}
	return apiKey, nil

}
