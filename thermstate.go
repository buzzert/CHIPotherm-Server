package main

import (
	"fmt"
	"log"
)

// ThermostatState represents the state handed down from the command and control
// server to the client (the thermostat itself)
type ThermostatState struct {
	enabled     bool
	temperature float32
}

// StateFromString returns a ThermostatState based on a string with the following format:
// [enabled | disabled] [temperature f]
func StateFromString(str string) ThermostatState {
	var temperature float32
	var enabledString string
	_, err := fmt.Sscanf(str, "%s %f", &enabledString, &temperature)
	if err != nil {
		log.Printf("Error parsing state")
	}

	enabled := false
	if enabledString == "enabled" {
		enabled = true
	}

	return ThermostatState{
		enabled:     enabled,
		temperature: temperature,
	}
}

func (state ThermostatState) String() string {
	stateString := "enabled"
	if !state.enabled {
		stateString = "disabled"
	}

	return fmt.Sprintf("%s %f", stateString, state.temperature)
}
