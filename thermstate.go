package main

import (
	"fmt"
	"log"
)

// ThermostatState represents the state handed down from the command and control
// server to the client (the thermostat itself)
type ThermostatState struct {
	Enabled            bool    `json:"Enabled,string"`
	HeatOn             bool    `json:"HeatOn,string"`
	TargetTemperature  float64 `json:"TargetTemperature,string"`
	CurrentTemperature float64 `json:"CurrentTemperature,string"`
}

// StateFromString returns a ThermostatState based on a string with the following format:
// [enabled | disabled] [temperature f]
func StateFromString(str string) ThermostatState {
	log.Printf("STR: %s", str)
	var temperature float64
	var enabledString string
	_, err := fmt.Sscanf(str, "%s %f", &enabledString, &temperature)
	if err != nil {
		log.Printf("Error parsing state: %s", err)
	}

	enabled := false
	if enabledString == "enabled" {
		enabled = true
	}

	return ThermostatState{
		Enabled:           enabled,
		TargetTemperature: temperature,
	}
}

func (state ThermostatState) String() string {
	stateString := "enabled"
	if !state.Enabled {
		stateString = "disabled"
	}

	return fmt.Sprintf("%s %f", stateString, state.TargetTemperature)
}
