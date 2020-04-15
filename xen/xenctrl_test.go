package xen

import (
	"testing"
	"fmt"
	"strings"
	"os/exec"
	"bytes"
)

func isPaused() bool {
	cmd := exec.Command("xl", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	output := out.String()
	rows := strings.Split(output, "\n")
	for _, row := range rows {
		if strings.Contains(row, "stardust") {
			values := strings.Fields(row)
			state := values[4]
			return strings.Contains(state, "p")
		}
	}
	return false
}


func TestPause(t *testing.T) {
	domainid := testSetup()
	cntrl := Xenctrl{DomainID: uint32(domainid)}
	cntrl.Init()

	err := cntrl.Pause()
	if err != nil {
		t.Error(err)
		return 
	}
	if !isPaused() {
		t.Error("Error: did not pause VM")
	}

	testTeardown()
}