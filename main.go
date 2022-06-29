package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	HACKERONE_UNAME = "USERNAME"
	HACKERONE_TOKEN = "H1_API_TOKEN"
	headers         = map[string][]string{"Accept": []string{"application/json"}}
)

type Data struct {
	Data []Program `json:"data"`
}

type Program struct {
	Id            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type Attributes struct {
	Handle          string `json:"handle"`
	Name            string `json:"name"`
	SubmissionState string `json:"submission_state"`
	State           string `json:"state"`
	OffersBounty    bool   `json:"offers_bounties"`
}

type Relationships struct {
	StructedScopes StructedScope `json:"structured_scopes"`
}

type StructedScope struct {
	Data []StructedScopeData `json:"data"`
}

type StructedScopeData struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	ScopeAttrs ScopeAttrs `json:"attributes"`
}

type ScopeAttrs struct {
	Asset_type              string `json:"asset_type"`
	Asset_identifier        string `json:"asset_identifier"`
	Eligible_for_bounty     bool   `json:"eligible_for_bounty"`
	Eligible_for_submission bool   `json:"eligible_for_submission"`
	Instruction             string `json:"instruction"`
	Max_severity            string `json:"max_severity"`
	Created_at              string `json:"created_at"`
	Updated_at              string `json:"updated_at"`
}

func getProgramHandleH1(p Program) Program {
	ph := p.Attributes.Handle
	req, err := http.NewRequest("GET", "https://api.hackerone.com/v1/hackers/programs/"+ph, nil)
	if err != nil {
		panic(err)
	}
	req.Header = headers
	req.SetBasicAuth(HACKERONE_UNAME, HACKERONE_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var program Program
	json.Unmarshal(body, &program)
	return program

}

func getProgramsH1() []Program {
	req, err := http.NewRequest("GET", "https://api.hackerone.com/v1/hackers/programs?page[number]=5&page[size]=100", nil)
	if err != nil {
		panic(err)
	}
	req.Header = headers
	req.SetBasicAuth(HACKERONE_UNAME, HACKERONE_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data Data

	json.Unmarshal(body, &data)

	var programs = []Program{}
	for i := 0; i < len(data.Data); i++ {
		if data.Data[i].Type == "program" && data.Data[i].Attributes.OffersBounty && data.Data[i].Attributes.SubmissionState == "open" {
			programs = append(programs, data.Data[i])
		}
	}
	return programs
}
