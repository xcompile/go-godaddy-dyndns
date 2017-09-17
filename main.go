/**
* Dynamic DNS using GoDaddy API 
* https://developer.godaddy.com
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"
	"errors"
	"io/ioutil"
	"flag"
)

type DomainWithRecord struct {
	
	Domain	string `json:"domain"`
	Record	string `json:"record"`
	Type    string `json:"type"`
}

type DomainRecord struct {
	Type	string `json:"type"`
	Name	string `json:"name"`
	Data	string `json:"data"`
	TTL	int    `json:"ttl"`
}

type DNSRecordCreateTypeName struct {
	Data	    string `json:"data"`
	TTL         int    `json:"ttl"`
}

type IPFy struct{
	IP	string `json:"ip"`
}

var (
  GoDaddyHttp = initHttpClient()
  ApiKey, ApiSecret, Domain = initParams()
)

func initParams() (string,string,DomainWithRecord) {
  var ak string
  var as string
  var dwrStr string

  flag.StringVar(&ak,"key","", "GoDaddy API-Key")
  flag.StringVar(&as,"secret","", "GoDaddy API-Secret")
  flag.StringVar(&dwrStr,"dwr","", "GoDaddy Domain (JSON: {\"domain\": \"yourdomain.com\",\"record\": \"yourdnsrecord\",\"type\": \"A\"}'")
  flag.Parse()

  byt := []byte(dwrStr)

  var dwr DomainWithRecord
  if err := json.Unmarshal(byt, &dwr); err != nil {
    panic(err)
  }

  return ak, as, dwr

}

func initHttpClient() *http.Client {
	return &http.Client{}
}


func main() {

	ip := getExternalIP()
	currentIp, err := getRecord()
	if (nil != err) {
		panic(err)
	}

	if (currentIp == ip) {
		log.Printf("skip update, current assigned ip matches external ip => %s", ip)
		return
	}

	err = updateRecord(Domain.Domain,Domain.Type, Domain.Record, ip)
	if (nil != err) {
		panic(err)
	}





}
func getRecord() (string, error) {
	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",Domain.Domain,Domain.Type,Domain.Record)
	auth := fmt.Sprintf("sso-key %s:%s",ApiKey,ApiSecret)


	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth)


	resp, err := GoDaddyHttp.Do(req)
	if (err != nil) {
		log.Fatal("Do: ",err)
		return "",err
	}
	defer resp.Body.Close()

	// list of records for domain
	records := make([]DomainRecord,0)

	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		log.Println(err)
		return "",err
	}
	return records[0].Data, nil
}

func updateRecord(domain string,recordType string, recordName string,ip string) error {
        url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",domain,recordType,recordName)
	auth := fmt.Sprintf("sso-key %s:%s",ApiKey,ApiSecret)

	domainRecords  := make([]DNSRecordCreateTypeName,1)
	domainRecord := DNSRecordCreateTypeName{
		Data: ip,
		TTL: 600,
	}
	domainRecords[0] = domainRecord

	GoDaddyHttp := &http.Client{}

	json, err := json.Marshal(domainRecords)
	if (err != nil) {
		return err
	}
	log.Print(string(json))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")

	resp, err := GoDaddyHttp.Do(req)
	if (err != nil) {
		log.Fatal("Do: ",err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) {
		raw, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("status code %d, expected 200 OK: %s", resp.StatusCode, string(raw)))
	}

	log.Print(fmt.Sprintf("update %s record %s to %s was successfull", domain, recordName, ip)) 
	return nil
}

//function to get the public ip address
func getExternalIP() string {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if (err != nil) {
		log.Fatal("Do: ", err)
		return ""
	}
	defer resp.Body.Close()

	var record IPFy

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Fatal("Decode ", err)
		return ""
	}
	return record.IP

}
