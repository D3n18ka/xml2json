package service_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"maxim.tbank/xml2pg/service"
)

//go:embed data.xml
var data []byte

func TestStreamXml(t *testing.T) {

	err := service.ParseXml(data)
	if err != nil {
		t.Error("error parsing xml")
	}
}

func TestParseXml(t *testing.T) {

	msg, err := service.ParseXmlStruct(data)
	if err != nil {
		t.Error("error parsing xml")
	}

	marshal, err := json.Marshal(msg)
	if err != nil {
		t.Error("error marshaling json")
	}

	fmt.Println(string(marshal))
}
