let ENABLE_BUTTON = "enable_button";
let INC_BUTTON = "increase_temp";
let DEC_BUTTON = "decrease_temp";
let CURRENT_TEMP = "current_temp";
let TARGET_TEMP = "target_temp";
function elem(id) { return document.getElementById(id); }

var is_enabled = false;
var target_temperature = 0.0;
var current_temperature = 0.0;
var heat_on = false;

function refreshState() {
    let request = new XMLHttpRequest();
    request.open("GET", "/refreshState");
    request.onload = () => {
        if (request.status == 200) {
            let state = JSON.parse(request.responseText);
            is_enabled = (state["Enabled"] == "true");
            heat_on = (state["HeatOn"] == "true");
            target_temperature = Math.round(parseFloat(state["TargetTemperature"]))
            current_temperature = Math.round(parseFloat(state["CurrentTemperature"]))

            setStatus("Connected.");
            updateUI();
        } else {
            console.log("Error getting response");
            setStatus("Error connecting to server");
        }
    }

    request.onerror = (err) => {
        console.log("err");
        setStatus("Error connecting to CHIP: " + err);
    };
    
    request.send();
}

function updateUI() {
    // Enable button
    let enable_button = elem(ENABLE_BUTTON);
    if (is_enabled) {
        enable_button.classList.add("enabled");
    } else {
        enable_button.classList.remove("enabled");
    }
    enable_button.innerHTML = (is_enabled ? "DISABLE" : "ENABLE");

    // Target temp
    let target_temp = elem(TARGET_TEMP);
    target_temp.style.opacity = (is_enabled ? 1.0 : 0.5);
    target_temp.innerHTML = target_temperature;

    let PULSE_CLASS = "pulse_animation";
    if (heat_on) {
        elem(TARGET_TEMP).classList.add(PULSE_CLASS);
    } else {
        elem(TARGET_TEMP).classList.remove(PULSE_CLASS);
    }

    // Current temp
    elem(CURRENT_TEMP).innerHTML = current_temperature;
}

function sendNewState()
{
    let request = new XMLHttpRequest();
    request.open("POST", "/setState");

    let body = (is_enabled ? "enabled" : "disabled") + " " + target_temperature.toString();
    console.log("Sending: " + body);
    request.send(body);
}

function setEnabled(enabled) {
    is_enabled = enabled;
    sendNewState();
    updateUI();
}

function setStatus(status) {
    var status_label = document.getElementById("status");
    status_label.innerHTML = status;
}

function setHeatOn(on) {
    heat_on = on;
    updateUI();
}

function setTargetTemperature(target) {
    target_temperature = target;
    sendNewState();
    updateUI();
}

window.onload = () => {
    setStatus("Connecting...");

    elem(ENABLE_BUTTON).onclick = () => {
        setEnabled(!is_enabled);
    }

    elem(DEC_BUTTON).onclick = () => {
        setTargetTemperature(target_temperature - 1.0);
    };

    elem(INC_BUTTON).onclick = () => {
        setTargetTemperature(target_temperature + 1.0);
    }

    refreshState();
    setInterval(refreshState, 2500);
}
