package gonsx

import (
	"fmt"
	"log"
	"encoding/xml"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"github.com/sky-uk/gonsx/api"
	"bytes"
	"io"
	"strings"
)

func NewNSXClient(url string, user string, password string, ignoreSSL bool, debug bool) *NSXClient {
	nsxClient := new(NSXClient)
	nsxClient.URL = url
	nsxClient.User = user
	nsxClient.Password = password
	nsxClient.IgnoreSSL = ignoreSSL
	nsxClient.debug = debug
	return nsxClient
}

type NSXClient struct {
	URL		string
	User 		string
	Password	string
	IgnoreSSL	bool
	debug 		bool
}

func (nsxClient *NSXClient) Do(api api.NSXApi) error {
	requestURL := fmt.Sprintf("%s%s", nsxClient.URL, api.Endpoint())

	var requestPayload io.Reader
	if(api.RequestObject() != nil) {
		requestXmlBytes, marshallingErr := xml.Marshal(api.RequestObject())
		log.Println(string(requestXmlBytes))
		if marshallingErr != nil {
			log.Fatal(marshallingErr)
		}
		requestPayload = bytes.NewReader(requestXmlBytes)
	}
	if(nsxClient.debug) {
		log.Println("requestURL:", requestURL)
	}
	req, err := http.NewRequest(api.Method(), requestURL, requestPayload)
	if err != nil {
		log.Println("ERROR building the request: %s", err)
		return err
	}

	req.SetBasicAuth(nsxClient.User, nsxClient.Password)
	// TODO: remove this hardcoded value!
	req.Header.Set("Content-Type", "application/xml")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: nsxClient.IgnoreSSL},
	}
	httpClient := &http.Client{Transport: tr}
	res, err := httpClient.Do(req)
	if err != nil{
		log.Println("ERROR executing request: ", err)
		return err
	}
	defer res.Body.Close()
	return nsxClient.handleResponse(api, res)
}

func (nsxClient *NSXClient) handleResponse(api api.NSXApi, res *http.Response) error {
	api.SetStatusCode(res.StatusCode)
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Println("ERROR reading response: ", err)
		return err
	}

	api.SetRawResponse(bodyText)

	if(nsxClient.debug) {
		log.Println(string(bodyText))
	}

	if (isXML(res.Header.Get("Content-Type")) && api.StatusCode() == 200) {
		xmlerr := xml.Unmarshal(bodyText, api.ResponseObject())
		if xmlerr != nil {
			log.Println("ERROR unmarshalling response: ", err)
			return err
		}
	} else {
		api.SetResponseObject(string(bodyText))
	}
	return nil
}

func isXML(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "/xml")
}
