package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"math"
	"strconv"
)

const (
	ProviderCmeApiVersion = "v1.2"
	CmeApiPath            = "cme-api/" + ProviderCmeApiVersion
)


func validateCmeApiVersion(c *checkpoint.ApiClient) error {
	apiVersion, err := getCmeApiVersion(c)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if ProviderCmeApiVersion > apiVersion {
		return fmt.Errorf("invalid Provider CME API version, it must be equal or lower than CME API version in" +
			" the Management. Provider CME API version: %s, CME API version in the Management: %s", ProviderCmeApiVersion, apiVersion)
	}
	return nil
}

func getCmeApiVersion(c *checkpoint.ApiClient) (string, error) {
	url := "cme-api/api-versions"
	res, err := c.ApiCall(url, nil, c.GetSessionID(), true, false, "GET")
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	cmeAPIVersionsJson := res.GetData()
	if checkIfRequestFailed(cmeAPIVersionsJson) {
		errMessage := buildErrorMessage(cmeAPIVersionsJson)
		return "", fmt.Errorf(errMessage)
	}
	ApiVersions := cmeAPIVersionsJson["result"].(map[string]interface{})
	return ApiVersions["current_version"].(string), nil
}

func checkIfRequestFailed(resJson map[string]interface{}) bool {

	if resJson["status-code"] != nil {
		statusCode := resJson["status-code"].(float64)
		if int(math.Round(statusCode)) != 200 {
			return true
		}
	}
	return false
}

func buildErrorMessage(resJson map[string]interface{}) string {
	errMessage := ""
	if resJson["error"] != nil {
		errorResultJson := resJson["error"].(map[string]interface{})
		if v := errorResultJson["message"]; v != nil {
			errMessage = "Message: " + v.(string)
		}
		if v := errorResultJson["details"]; v != nil {
			errMessage += ". Details: " + v.(string)
		}
		if v := errorResultJson["error-code"]; v != nil {
			errMessage += " (Error code: " + strconv.Itoa(int(math.Round(v.(float64)))) + ")"
		}
	}
	if errMessage == "" {
		errMessage = "Request failed. For more details check cme_api logger on the management server"
	}
	return errMessage
}

func cmeObjectNotFound(resJson map[string]interface{}) bool {
	NotFoundErrorCode := []int{800, 802}
	if resJson["error"] != nil {
		errorResultJson := resJson["error"].(map[string]interface{})
		if v := errorResultJson["error-code"]; v != nil {
			errorCode := int(math.Round(v.(float64)))
			for i := range NotFoundErrorCode {
				if errorCode == NotFoundErrorCode[i] {
					return true
				}
			}
		}
	}
	return false
}
