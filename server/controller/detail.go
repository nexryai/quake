package controller

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/nexryai/quake/server/feed"
	"io"
	"log"
	"strings"

	"github.com/nexryai/quake/server/converter"
	"github.com/nexryai/quake/server/jmaseis"
)

const eventDataBaseUrl = "https://www.data.jma.go.jp/developer/xml/data/"
const force = false
const ignoreWarning = false

func GetEventDetailsJson(eventId string, isDebug bool) (*string, error) {
	baseUrl := eventDataBaseUrl
	if isDebug {
		baseUrl = "https://raw.githubusercontent.com/nexryai/quake/main/test/examples/"
	}

	respReader, err := feed.FetchEventData(baseUrl, eventId)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(respReader)

	report := &jmaseis.Report{}
	err = xml.Unmarshal(data, &report)
	if err != nil {
		return nil, err
	}

	if strings.Contains(eventId, "_VTSE41") {
		// 津波
		jmaTsunami, err := converter.Vtse2Epsp(*report)
		if err != nil {
			log.Fatalf("%s convert error: %#v", eventId, err)
		}

		// 検証
		errors := converter.ValidateJMATsunami(eventId, report, jmaTsunami)
		for _, err := range errors {
			_, ok := err.(converter.ValidationError)
			if ok && !force {
				log.Fatalf("%s has validation errors: %#v", eventId, errors)
			}
		}

		for _, err := range errors {
			_, ok := err.(converter.ValidationWarning)
			if ok && !force && !ignoreWarning {
				log.Fatalf("%s has validation warnings: %#v", eventId, errors)
			}
		}

		data, err = json.Marshal(jmaTsunami)
		if err != nil {
			log.Fatalf("%s JSON conversion error: %#v", eventId, err)
		}

		jsonString := string(data)

		return &(jsonString), nil
	} else if strings.Contains(eventId, "_VXSE43") {
		// EEW
		jmaEEW, err := converter.Vxse2EpspEEW(*report)
		if err != nil {
			return nil, err
		}

		data, err = json.Marshal(jmaEEW)
		if err != nil {
			return nil, err
		}

		jsonString := string(data)

		return &(jsonString), nil
	} else {
		// 変換
		jmaQuake, err := converter.Vxse2EpspQuake(*report)
		if err != nil {
			return nil, err
		}

		// 検証
		validationErrors := converter.ValidateJMAQuake(eventId, report, jmaQuake)
		for _, err := range validationErrors {
			var validationError converter.ValidationError
			ok := errors.As(err, &validationError)
			if ok && !force {
				return nil, fmt.Errorf("%s has validation errors: %#v", eventId, validationErrors)
			}
		}

		for _, err := range validationErrors {
			var validationWarning converter.ValidationWarning
			ok := errors.As(err, &validationWarning)
			if ok && !force && !ignoreWarning {
				return nil, fmt.Errorf("%s has validation errors: %#v", eventId, validationErrors)
			}
		}

		data, err = json.Marshal(jmaQuake)
		if err != nil {
			log.Fatalf("%s JSON conversion error: %#v", eventId, err)
		}

		jsonString := string(data)

		return &(jsonString), nil
	}

	return nil, nil
}