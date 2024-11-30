package service_test

import (
	_ "embed"
	"testing"

	"maxim.tbank/xml2pg/service"
)

//go:embed data.xml
var dataXml []byte

func BenchmarkCsvMarshal(b *testing.B) {

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		service.ParseXml(dataXml)
	}
}
