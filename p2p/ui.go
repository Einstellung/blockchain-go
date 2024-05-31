package p2p

import (
	cli "blockchain-go/cmd/cli"
	utils "blockchain-go/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type BlockchainUi struct {
	inputChan chan string
	doneCh    chan struct{}

	bcommu *BlockchainCommunication
}

func NewBlockchainUi(bcommu *BlockchainCommunication) *BlockchainUi {

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	
	cli.PrintUsage()

	return &BlockchainUi{
		inputChan: inputCh,
		doneCh:    make(chan struct{}, 1),

		bcommu: bcommu,
	}
}

func (ui *BlockchainUi) end() {
	ui.doneCh <- struct{}{}
}

func (ui *BlockchainUi) handleEvents() {
	cli := cli.CLI{}

	for {
		select {
		// when the user types in a line, publish it to the blockchain network and print to the message window
		case input := <-ui.inputChan:
			cli.Run(input)
			err := ui.bcommu.Publish(input)
			if err != nil {
				utils.PrintErr("publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		// when we receive a message from the blockchain network, print it to the message window
		case m := <-ui.bcommu.Messages:
			ui.displayPeerMessage(m)

		case <-ui.doneCh:
			return
		}
	}
}

func (ui *BlockchainUi) displayPeerMessage(m BlockMessage) {
	utils.ColorPrint("green", fmt.Sprintf("<%s>: ", utils.ShortID(m.SenderID)), m.Message)
}

func (ui *BlockchainUi) displaySelfMessage(input string) {
	utils.ColorPrint("yellow", "<self>: ", input)
	fmt.Println()
}


// read data from terminal input then send to input channel
func (ui *BlockchainUi) readDataTerminal() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		// fmt.Print("terminal input> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		// Trim leading and trailing whitespace, including '\n' from the input string
		sendData = strings.TrimSpace(sendData)

		if sendData == "" {
			return
		}
		ui.inputChan <- sendData
	}
}
