package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shreddedbacon/fronius-client/fronius"
)

// FakePowerwall holds the value for the inverter host that will get passed to the fronius client
type FakePowerwall struct {
	Inverter     string
	InverterType string
	CloudBased   bool
	AuthToken    string
}

// MetersAggregates is the data structure for `/api/meters/aggretates` response from the PW api.
type MetersAggregates struct {
	Site    *AggregateData `json:"site,omitempty"`
	Battery *AggregateData `json:"battery,omitempty"`
	Load    *AggregateData `json:"load,omitempty"`
	Solar   *AggregateData `json:"solar,omitempty"`
}

// AggregateData is the data structure for each of the meter aggregates.
type AggregateData struct {
	LastCommunicationTime time.Time `json:"last_communication_time"`
	InstantPower          float64   `json:"instant_power"`
	InstantReactivePower  float64   `json:"instant_reactive_power"`
	InstantApparentPower  float64   `json:"instant_apparent_power"`
	Frequency             float64   `json:"frequency"`
	EnergyExported        float64   `json:"energy_exported"`
	EnergyImported        float64   `json:"energy_imported"`
	InstantAverageVoltage float64   `json:"instant_average_voltage"`
	InstantTotalCurrent   int       `json:"instant_total_current"`
	IACurrent             int       `json:"i_a_current"`
	IBCurrent             int       `json:"i_b_current"`
	ICCurrent             int       `json:"i_c_current"`
	Timeout               int       `json:"timeout"`
}

// GetMetersAggregates returns a response that is similar to `/api/meters/aggregates` that a local PW would return
func (f *FakePowerwall) GetMetersAggregates(w http.ResponseWriter, r *http.Request) {
	pwd := MetersAggregates{}

	switch f.InverterType {
	case "fronius":
		d, _ := fronius.New(f.Inverter)
		// Get the realtime powerflow data
		p, _ := d.GetPowerFlowRealtimeData()
		// Craft our response payload, we only need the instant power value as this is
		// what is needed to determine how many watts are being fed into the grid
		pwd = MetersAggregates{
			Site: &AggregateData{
				LastCommunicationTime: time.Now().UTC(),
				InstantPower:          p.Body.Data.Site.PGrid, // This is what we need from the inverter
			},
			Load: &AggregateData{
				LastCommunicationTime: time.Now().UTC(),
				InstantPower:          p.Body.Data.Site.PLoad * -1, // powerwall inverts this value
			},
			Solar: &AggregateData{
				LastCommunicationTime: time.Now().UTC(),
				InstantPower:          p.Body.Data.Site.PPV,
			},
			Battery: &AggregateData{
				LastCommunicationTime: time.Now().UTC(),
				InstantPower:          0,
			},
		}
	}
	// Marshal the data into json bytes
	pwb, _ := json.Marshal(pwd)
	// Return the bytes as string
	fmt.Fprintln(w, string(pwb))
}

// Request performs a request against a powerwall api endpoint
func (f *FakePowerwall) Request(endpoint string) ([]byte, error) {
	netClientTimeout := 10
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var netClient = &http.Client{
		Timeout:   time.Second * time.Duration(netClientTimeout),
		Transport: tr,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", f.Inverter, endpoint), bytes.NewBuffer([]byte{}))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	rBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return []byte{}, errors.New("error performing check or connecting to powerwall")
	}
	return rBody, nil
}
