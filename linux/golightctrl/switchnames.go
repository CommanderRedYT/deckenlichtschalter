// (c) Bernhard Tittelbach, 2016
package main

import "fmt"

type ActionNameHandler struct {
	handler     string
	codeon      []byte
	codeoff     []byte
	codedefault []byte
}

const (
	IRCmd2MQTT            = "IRCmd2MQTT"
	RFCode2TTY            = "RFCode2TTY"
	RFCode2BOTH           = "RFCode2BOTH"
	RFCode2MQTT           = "RFCode2MQTT"
	CeilingLightByteState = "CeilingLightByteState"
)

var RF433_chan_ chan SerialLine
var MQTT_rf_chan_ chan []byte
var MQTT_ir_chan_ chan string

var actionname_map_ map[string]ActionNameHandler = map[string]ActionNameHandler{
	"regalleinwand": ActionNameHandler{codeon: []byte{0xa2, 0xa0, 0xa8}, codeoff: []byte{0xa2, 0xa0, 0x28}, handler: RFCode2TTY}, //white remote B 1
	"bluebar":       ActionNameHandler{codeon: []byte{0xa8, 0xa0, 0xa8}, codeoff: []byte{0xa8, 0xa0, 0x28}, handler: RFCode2TTY}, //white remote C 1
	"labortisch":    ActionNameHandler{codeon: []byte{0xa2, 0xa2, 0xaa}, codeoff: []byte{0xa2, 0xa2, 0x2a}, handler: RFCode2TTY},
	"couchred":      ActionNameHandler{codeon: []byte{0x8a, 0xa0, 0x8a}, codeoff: []byte{0x8a, 0xa0, 0x2a}, handler: RFCode2TTY},  //pollin 00101 a
	"couchwhite":    ActionNameHandler{codeon: []byte{0x8a, 0xa8, 0x88}, codeoff: []byte{0x8a, 0xa8, 0x28}, handler: RFCode2TTY},  //pollin 00101 d
	"cxleds":        ActionNameHandler{codeon: []byte{0x8a, 0x88, 0x8a}, codeoff: []byte{0x8a, 0x88, 0x2a}, handler: RFCode2TTY},  //pollin 00101 b
	"mashadecke":    ActionNameHandler{codeon: []byte{0x8a, 0x28, 0x8a}, codeoff: []byte{0x8a, 0x28, 0x2a}, handler: RFCode2TTY},  //pollin 00101 c
	"boiler":        ActionNameHandler{codeon: []byte{0xa0, 0xa2, 0xa8}, codeoff: []byte{0xa0, 0xa2, 0x28}, handler: RFCode2BOTH}, //white remote A 2
	"spots":         ActionNameHandler{codeon: []byte{0x00, 0xaa, 0x88}, codeoff: []byte{0x00, 0xaa, 0x28}, handler: RFCode2TTY},  //polling 11110 d
	"olgatemp":      ActionNameHandler{codeon: []byte{0x00, 0xa2, 0x8a}, codeoff: []byte{0x00, 0xa2, 0x2a}, handler: RFCode2TTY},  // Funksteckdose an welcher olgafreezer.realraum.at hängt
	"abwasch":       ActionNameHandler{codeon: []byte{0xaa, 0xa2, 0xa8}, codeoff: []byte{0xaa, 0xa2, 0x28}, handler: RFCode2MQTT}, //alte jk16 decke vorne

	"ymhpoweroff":  ActionNameHandler{codedefault: []byte("ymhpoweroff"), handler: IRCmd2MQTT},
	"ymhpower":     ActionNameHandler{codedefault: []byte("ymhpower"), codeoff: []byte("ymhpoweroff"), handler: IRCmd2MQTT},
	"ymhpoweron":   ActionNameHandler{codedefault: []byte("ymhpoweron"), handler: IRCmd2MQTT},
	"ymhcd":        ActionNameHandler{codedefault: []byte("ymhcd"), handler: IRCmd2MQTT},
	"ymhtuner":     ActionNameHandler{codedefault: []byte("ymhtuner"), handler: IRCmd2MQTT},
	"ymhtape":      ActionNameHandler{codedefault: []byte("ymhtape"), handler: IRCmd2MQTT},
	"ymhwdtv":      ActionNameHandler{codedefault: []byte("ymhwdtv"), handler: IRCmd2MQTT},
	"ymhsattv":     ActionNameHandler{codedefault: []byte("ymhsattv"), handler: IRCmd2MQTT},
	"ymhvcr":       ActionNameHandler{codedefault: []byte("ymhvcr"), handler: IRCmd2MQTT},
	"ymh7":         ActionNameHandler{codedefault: []byte("ymh7"), handler: IRCmd2MQTT},
	"ymhaux":       ActionNameHandler{codedefault: []byte("ymhaux"), handler: IRCmd2MQTT},
	"ymhextdec":    ActionNameHandler{codedefault: []byte("ymhextdec"), handler: IRCmd2MQTT},
	"ymhtest":      ActionNameHandler{codedefault: []byte("ymhtest"), handler: IRCmd2MQTT},
	"ymhtunabcde":  ActionNameHandler{codedefault: []byte("ymhtunabcde"), handler: IRCmd2MQTT},
	"ymheffect":    ActionNameHandler{codedefault: []byte("ymheffect"), handler: IRCmd2MQTT},
	"ymhtunplus":   ActionNameHandler{codedefault: []byte("ymhtunplus"), handler: IRCmd2MQTT},
	"ymhtunminus":  ActionNameHandler{codedefault: []byte("ymhtunminus"), handler: IRCmd2MQTT},
	"ymhvolup":     ActionNameHandler{codedefault: []byte("ymhvolup"), handler: IRCmd2MQTT},
	"ymhvoldown":   ActionNameHandler{codedefault: []byte("ymhvoldown"), handler: IRCmd2MQTT},
	"ymhvolmute":   ActionNameHandler{codedefault: []byte("ymhvolmute"), handler: IRCmd2MQTT},
	"ymhmenu":      ActionNameHandler{codedefault: []byte("ymhmenu"), handler: IRCmd2MQTT},
	"ymhplus":      ActionNameHandler{codedefault: []byte("ymhplus"), handler: IRCmd2MQTT},
	"ymhminus":     ActionNameHandler{codedefault: []byte("ymhminus"), handler: IRCmd2MQTT},
	"ymhtimelevel": ActionNameHandler{codedefault: []byte("ymhtimelevel"), handler: IRCmd2MQTT},
	"ymhprgdown":   ActionNameHandler{codedefault: []byte("ymhprgdown"), handler: IRCmd2MQTT},
	"ymhprgup":     ActionNameHandler{codedefault: []byte("ymhprgup"), handler: IRCmd2MQTT},
	"ymhsleep":     ActionNameHandler{codedefault: []byte("ymhsleep"), handler: IRCmd2MQTT},
	"ymhp5":        ActionNameHandler{codedefault: []byte("ymhp5"), handler: IRCmd2MQTT},

	"ceiling1": ActionNameHandler{codeon: []byte{0, 1}, codeoff: []byte{0, 0}, handler: CeilingLightByteState},
	"ceiling2": ActionNameHandler{codeon: []byte{1, 1}, codeoff: []byte{1, 0}, handler: CeilingLightByteState},
	"ceiling3": ActionNameHandler{codeon: []byte{2, 1}, codeoff: []byte{2, 0}, handler: CeilingLightByteState},
	"ceiling4": ActionNameHandler{codeon: []byte{3, 1}, codeoff: []byte{3, 0}, handler: CeilingLightByteState},
	"ceiling5": ActionNameHandler{codeon: []byte{4, 1}, codeoff: []byte{4, 0}, handler: CeilingLightByteState},
	"ceiling6": ActionNameHandler{codeon: []byte{5, 1}, codeoff: []byte{5, 0}, handler: CeilingLightByteState},
}

func SwitchName(name string, onoff bool) error {
	nm, inmap := actionname_map_[name]
	if !inmap {
		LogRF433_.Printf("Name %s does not exist in actionname_map_", name)
		return fmt.Errorf("Name does not exist")
	}
	LogRF433_.Printf("SwitchName(%s,%s", name, onoff)
	var handler func([]byte) error
	switch nm.handler {
	case IRCmd2MQTT:
		handler = sendIRCmd2MQTT
	case RFCode2TTY:
		handler = sendRFCode2TTY
	case RFCode2BOTH:
		handler = sendRFCode2BOTH
	case RFCode2MQTT:
		handler = sendRFCode2MQTT
	case CeilingLightByteState:
		handler = setCeilingLightByteState
	default:
		return fmt.Errorf("Unknown handler %s", nm.handler)

	}
	if onoff && nm.codeon != nil {
		return handler(nm.codeon)
	} else if onoff == false && nm.codeoff != nil {
		return handler(nm.codeoff)
	} else if nm.codedefault != nil {
		return handler(nm.codedefault)
	}
	return fmt.Errorf("SwitchName could not do anything")
}

func sendRFCode2TTY(code []byte) error {
	LogRF433_.Printf("RFCode2TTY(%+v)", code)
	RF433_chan_ <- append([]byte(">"), code...)
	return nil
}

func sendRFCode2MQTT(code []byte) error {
	LogRF433_.Printf("sendRFCode2MQTT(%+v)", code)
	MQTT_rf_chan_ <- code
	return nil
}

func sendIRCmd2MQTT(code []byte) error {
	LogRF433_.Printf("IRCmd2MQTT(%s)", string(code))
	MQTT_ir_chan_ <- string(code)
	return nil
}

func sendRFCode2BOTH(code []byte) error {
	sendRFCode2TTY(code)
	sendRFCode2MQTT(code)
	return nil
}

func setCeilingLightByteState(code []byte) error {
	if len(code) != 2 {
		LogRF433_.Printf("Invalid Code %s for setCeilingLightByteState", code)
		return fmt.Errorf("Invalid Code for setCeilingLightByteState")
	}
	SetCeilingLightsState(int(code[0]), code[1] == 1)
	return nil
}