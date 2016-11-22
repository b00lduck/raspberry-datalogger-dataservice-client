package client
import (
	"io/ioutil"
	"net/http"
	"errors"
	"strings"
	"math"
        "os"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

func request(url string, method string, body string) error {

	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		str, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("%d", resp.StatusCode) + " " + string(str))
	}

	return nil
}

func sendDataservicePut(url string, body string) error {
        dataservice := os.Getenv("DATASERVICE_HOST")
        return request("http://" + dataservice + "/" + url, "PUT", body)
}

func sendDataservicePost(url string, body string) error {
        dataservice := os.Getenv("DATASERVICE_HOST")
        return request("http://" + dataservice + "/" + url, "POST", body)
}

func SendCounterTick(code string) error {
	log.WithField("code", code).Info("Sending counter tick")
	return sendDataservicePost("counter/" + code + "/tick", "")
}

func SendThermometerReading(code string, temp float64) error {
	log.WithField("code", code).WithField("temp", fmt.Sprintf("%.2f°C", temp)).Info("Sending thermometer reading")
	svalue := fmt.Sprintf("%.0f", Round(temp * 1000))
	return sendDataservicePost("thermometer/" + code + "/reading", svalue)
}

func SendFlagState(code string, state uint8) error {
	log.WithField("code", code).WithField("state", state).Info("Sending flag state")
	svalue := fmt.Sprintf("%d", state)
	return sendDataservicePost("flag/" + code + "/state", svalue)
}

func SendCounterCorrection(code string, value int32) error {
	svalue := fmt.Sprintf("%d", value)
	log.WithField("code", code).WithField("value", value).Info("Sending counter correction")
	return sendDataservicePut("counter/" + code + "/corr", svalue)
}

func Round(f float64) float64 {
        return math.Floor(f + .5)
}
