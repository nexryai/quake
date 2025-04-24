package geo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/koron/go-dproxy"
	"github.com/nexryai/archer"
)

var (
	ErrRequestFailedWithStatus = fmt.Errorf("request failed with status")

	jmaAreaData *interface{}

	/*
		気象庁のAPIで使用される市レベルのコードの一覧
		[総務省のWebサイト](https://www.soumu.go.jp/denshijiti/code.html) からExcelファイルをダウンロードしCSVに変換し、
		cut -c1-5 input.txt > cities.txt で作れる。
	*/
	//go:embed cities.txt
	citiesData string

	cityCodes = strings.Split(strings.TrimSpace(citiesData), "\n")
)


type RevGeo struct {
	Results struct {
		MuniCd string `json:"muniCd"`
		Lv01Nm string `json:"lv01Nm"`
	} `json:"results"`
}

type City struct {
	ForecastCode string `json:"forecastcode"`
	PrefName     string `json:"prefname"`
	Cityname     string `json:"cityname"`
	Centers      string `json:"centers"`
	Offices      string `json:"offices"`
	Class10s     string `json:"class10s"`
	Class15s     string `json:"class15s"`
	Class20s     string `json:"class20s"`
}

func init() {
	const url = "https://www.jma.go.jp/bosai/common/const/area.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	requester := archer.SecureRequest{
		Request: req,
		TimeoutSecs: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, err := requester.Send()
	if err != nil {
		log.Printf("Failed to fetch JMA area data: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}
	defer resp.Body.Close()
	
	err = json.Unmarshal(data, &jmaAreaData)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
	}

	log.Printf("JMA area data loaded successfully: %d bytes", len(data))
}

func latLonToLocalCode(lat, lon float64) (string, error) {
	url := fmt.Sprintf("https://mreversegeocoder.gsi.go.jp/reverse-geocoder/LonLatToAddress?lat=%v&lon=%v", lat, lon)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	requester := archer.SecureRequest{
		Request: req,
		TimeoutSecs: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, err := requester.Send()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s", resp.Status)
		return "", ErrRequestFailedWithStatus
	}


	var rg RevGeo
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	json.Unmarshal(body, &rg)

	return rg.Results.MuniCd, nil
}

/*
	区レベルのコードを市レベルのコードに変換する。
	渡されたコードが既に市レベルのコードの場合はそのまま返す。
	区レベルのコードが渡された場合は、最も近い市レベルのコードを返す。
*/
func localCodeToCityCode(code string) string {
	// 二分探索
	left, right := 0, len(cityCodes)-1
	var mid int
	for right-left > 1 {
		mid = (left + right) / 2
		if code == cityCodes[mid] {
			break
		} else if code > cityCodes[mid] {
			left = mid
		} else {
			right = mid
		}
	}

	if code < cityCodes[mid] {
		mid--
	}

	return cityCodes[mid]
}

func getParentCode(code, class string) string {
	// category := []string{"class20s", "class15s", "class10s", "offices"}
	areainfo := dproxy.New(*jmaAreaData)

	pcode, err := areainfo.M(class).M(code).M("parent").String()
	if err == nil {
		return pcode
	}

	return ""
}

func toCityName(citycode string) string {
	areainfo := dproxy.New(*jmaAreaData)
	name, _ := areainfo.M("class20s").M(citycode).M("name").String()
	return name
}

func toPrefName(officecode string) string {
	areainfo := dproxy.New(*jmaAreaData)
	name, _ := areainfo.M("offices").M(officecode).M("name").String()
	return name
}

func LatLonToCity(lat, lon float64) (*City, error) {
	code, err := latLonToLocalCode(lat, lon)
	if err != nil {
		return nil, err
	}

	citycode := localCodeToCityCode(code)
	log.Default().Printf("City code: %s", citycode)

	class20s := citycode + "00"
	class15s := getParentCode(class20s, "class20s")
	class10s := getParentCode(class15s, "class15s")
	offices := getParentCode(class10s, "class10s")
	centers := getParentCode(offices, "offices")

	var forecastCode string
	switch offices {
		case "014030":
			forecastCode = "014100"
		case "460040":
			forecastCode = "460100"
		default:
			forecastCode = offices
	}

	return &City{
		ForecastCode: forecastCode,
		PrefName:     toPrefName(offices),
		Cityname:     toCityName(class20s),
		Centers:      centers,
		Offices:      offices,
		Class10s:     class10s,
		Class15s:     class15s,
		Class20s:     class20s,
	}, nil
}
