package simtest

import (
	"fmt"
	"testing"

	"github.com/FactomProject/factomd/common/globals"
	"github.com/FactomProject/factomd/engine"
	. "github.com/FactomProject/factomd/testHelper"
)

/*
Test brainswapping F <-> L with no auditors

This test is useful for catching a failure scenario where the timing between
identity swap is off leading to a stall
*/
func TestLeaderBrainSwap(t *testing.T) {
	ResetFactomHome(t)       // clear out old test home
	for i := 0; i < 6; i++ { // build config files for the test
		WriteConfigFile(i, i, "", t) // just write the minimal config
	}

	params := map[string]string{"--factomhome": globals.Params.FactomHome}

	// start the 6 fnodes running "012345"
	state0 := SetupSim("LLLFFF", params, 30, 0, 0, t)
	state3 := engine.GetFnodes()[3].State // Get node 2

	WaitForAllNodes(state0)
	WaitForBlock(state0, 6)

	batchCount := 3 // use odd number to match LFFFLL as end condition

	for batch := 0; batch < batchCount; batch++ {

		target := batch + 7

		change := fmt.Sprintf("ChangeAcksHeight = %v\n", target)

		if batch%2 == 0 {
			WriteConfigFile(1, 5, change, t) // Setup A brain swap between L1 and F5
			WriteConfigFile(5, 1, change, t)

			WriteConfigFile(2, 4, change, t) // Setup A brain swap between L2 and F4
			WriteConfigFile(4, 2, change, t)

		} else {
			WriteConfigFile(5, 5, change, t) // Un-Swap
			WriteConfigFile(1, 1, change, t)

			WriteConfigFile(4, 4, change, t)
			WriteConfigFile(2, 2, change, t)
		}

		WaitForBlock(state3, target)
		WaitMinutes(state3, 1)
	}

	WaitBlocks(state0, 1)
	AssertAuthoritySet(t, "LFFFLL")
	WaitForAllNodes(state0)
	ShutDownEverything(t)
}