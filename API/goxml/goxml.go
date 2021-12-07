package goxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type soapRQ struct {
	XMLName   xml.Name    `xml:"soapenv:Envelope"`
	XMLNsSoap string      `xml:"xmlns:soapenv,attr"`
	XMLNsXSI  string      `xml:"xmlns:ser,attr"`
	Body      interface{} `xml:"soapenv:Body"`
}

type UBResponse struct {
	Results string `xml:"Body>getResponse>return>results"`
}

func UBsoapCall(url string, soapAction string, payloadInterface interface{}) (string, error) {
	response, err := soapCall(url, soapAction, payloadInterface)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	//fmt.Printf("raw xml body: %v", body)

	var jsonResponse UBResponse

	err = xml.Unmarshal(body, &jsonResponse)
	if err != nil {
		fmt.Printf("unmarshal failed: %s", err)
		return "", err
	}

	//fmt.Printf("raw json response: %v", jsonResponse)

	return jsonResponse.Results, nil
}

func soapCall(ws string, action string, payloadInterface interface{}) (*http.Response, error) {
	v := soapRQ{
		XMLName: xml.Name{
			Space: "ser=", Local: "http://services.webservices.utilibill.com/"},
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXSI:  "http://services.webservices.utilibill.com/",
		Body:      payloadInterface,
	}

	payload, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal xml: %w", err)
	}

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		// TODO: bring logging into this package.
		fmt.Printf("print raw request: %q", dump)
		return nil, err
	}
	//fmt.Printf("print raw request: %q", dump)

	return client.Do(req)
}
