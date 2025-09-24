package query

import (
	"errors"

	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

func QueryChargeStatus(outletId string) (string, int64, error) {
	client := resty.New()
	resp, err := client.R().Get("https://wemp.issks.com/charge/v1/charging/outlet/" + outletId)

	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode() != 200 {
		return "", 0, errors.New("request failed with status code: " + resp.Status())
	}
	body := resp.Bytes()

	if gjson.GetBytes(body, "code").String() != "1" {
		return "", 0, errors.New("unexpected response code: " + gjson.GetBytes(body, "code").String())
	}

	power := gjson.GetBytes(body, "data.powerFee.billingPower").String()
	usedMinutes := gjson.GetBytes(body, "data.usedmin").Int()

	return power, usedMinutes, nil
}
