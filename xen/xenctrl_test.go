package xen

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
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

// Test that the Unpuase method works correctly
func TestUnpause(t *testing.T) {
	domainid := testSetup()
	cntrl := Xenctrl{DomainID: uint32(domainid)}
	cntrl.Init()

	err := cntrl.Pause()
	if err != nil {
		t.Error(err)
		return
	}

	err = cntrl.Unpause()
	if err != nil {
		t.Error(err)
		return
	}

	if isPaused() {
		t.Error("Error: did not unpause VM")
		return
	}

	testTeardown()
}

// Tests the Pause will pause the domain
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

// Tests the IsPaused method will correctly determine 
// whether the domain is paused
func TestIsPaused(t *testing.T) {
	domainid := testSetup()
	cntrl := Xenctrl{DomainID: uint32(domainid)}
	cntrl.Init()

	err := cntrl.Pause()
	if !isPaused() {
		t.Error("Error: did not pause VM")
		return
	}

	if !cntrl.IsPaused() {
		t.Error("Error: the the IsPaused did not return true")
	}

	err = cntrl.Unpause()
	if err != nil {
		t.Error(err)
		return
	}

	if cntrl.IsPaused() {
		t.Error("Error: VM is unpaused should be false")
	}
	testTeardown()
}
