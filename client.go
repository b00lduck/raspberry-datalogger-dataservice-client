package client
import (
	"io/ioutil"
	"net/http"
	"errors"
	"strings"
    "os"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

func request(url string, method string, body string) (string, error) {

	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	str, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("%d", resp.StatusCode) + " " + string(str))
	}

	return string(str), nil
}

func sendDataservicePut(url string, body string) error {
    dataservice := os.Getenv("DATASERVICE_HOST")
    fullUrl := "http://" + dataservice + "/" + url
    log.WithField("url", url).WithField("fullUrl", fullUrl).Info("Sending PUT")
    _, err := request(fullUrl, "PUT", body)
    return err
}

func sendDataservicePost(url string, body string) error {
    dataservice := os.Getenv("DATASERVICE_HOST")
    fullUrl := "http://" + dataservice + "/" + url
    log.WithField("url", url).WithField("fullUrl", fullUrl).Info("Sending POST")
    _, err := request(fullUrl, "POST", body)
    return err
}

func sendDataserviceGet(url string) (string, error) {
    dataservice := os.Getenv("DATASERVICE_HOST")
    fullUrl := "http://" + dataservice + "/" + url
    log.WithField("url", url).WithField("fullUrl", fullUrl).Info("Sending GET")
    return request(fullUrl, "GET", "")
}

func GetCounter(code string) (string, error) {
    log.WithField("code", code).Info("Getting counter value")
    return sendDataserviceGet("counter/" + code)
}

func SendCounterTick(code string) error {
	log.WithField("code", code).Info("Sending counter tick")
	return sendDataservicePost("counter/" + code + "/tick", "")
}

func SendThermometerReading(code string, temp float64) error {
	log.WithField("code", code).WithField("temp", fmt.Sprintf("%.2f°C", temp)).Info("Sending thermometer reading")
	svalue := fmt.Sprintf("%.2f", temp)
	return sendDataservicePut("thermometer/" + code, svalue)
}

func SendPercentage(code string, percent float64) error {
	log.WithField("code", code).WithField("percent", fmt.Sprintf("%.1f", percent)).Info("Sending percentager reading")
	svalue := fmt.Sprintf("%.1f", percent)
	return sendDataservicePut("percentage/" + code, svalue)
}

func SendFlagState(code string, state bool) error {
	log.WithField("code", code).WithField("state", state).Info("Sending flag state")
	var svalue string
    if state {
        svalue = "1"
    } else {
        svalue = "0"
    }
	return sendDataservicePut("flag/" + code, svalue)
}

func SendCounterCorrection(code string, value int32) error {
	svalue := fmt.Sprintf("%d", value)
	log.WithField("code", code).WithField("value", value).Info("Sending counter correction")
	return sendDataservicePut("counter/" + code, svalue)
}
