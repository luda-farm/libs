package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type viesResponse struct {
	IsValid bool `json:"isValid"`
}

func IsValidVatin(vatin string) (bool, error) {
	if len(vatin) < 3 {
		return false, nil
	}

	url := fmt.Sprintf(
		"https://ec.europa.eu/taxation_customs/vies/rest-api/ms/%s/vat/%s?requesterMemberStateCode=SE&requesterNumber=556690395001",
		vatin[0:2], vatin[2:],
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("creating http request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("sending http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("non-ok response from vies: %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body: %w", err)
	}

	var parsedResponse viesResponse
	if err := json.Unmarshal(resBody, &parsedResponse); err != nil {
		return false, fmt.Errorf("parsing response: %w", err)
	}

	return parsedResponse.IsValid, nil
}
