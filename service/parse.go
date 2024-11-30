package service

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"maxim.tbank/xml2pg/db"
)

func ParseXmlStruct(data []byte) (*db.Message, error) {
	msg := &db.Message{}
	err := xml.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func ParseXml(data []byte) error {

	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			fmt.Println("start", se.Name)

			for _, attr := range se.Attr {
				fmt.Println("attr", attr.Name, attr.Value)
			}
		case xml.CharData:
			fmt.Println(se.Copy())
		case xml.EndElement:
			break
		}
	}

	return nil
}
