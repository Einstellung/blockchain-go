package p2p

import (
	utils "blockchain-go/utils"
	"bufio"
	"fmt"
	"os"
)

type BlockchainUi struct {
	inputChan chan string
	doneCh    chan struct{}

	bcommu *BlockchainCommunication
}

func NewBlockchainUi(bcommu *BlockchainCommunication) *BlockchainUi {

	// an input field for typing messages into
	inputCh := make(chan string, 32)

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

	for {
		select {
		// when the user types in a line, publish it to the blockchain network and print to the message window
		case input := <-ui.inputChan:
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
	// prompt := utils.WithColor("green", fmt.Sprintf("<%s>:", utils.ShortID(m.SenderID)))
	// fmt.Fprintln(ui.peerMsgBox, prompt, m.Message)
	fmt.Println(m.Message)
}

func (ui *BlockchainUi) displaySelfMessage(m string) {
	// fmt.Fprintln(ui.peerMsgBox, "[self msg]:", m)
	fmt.Println("[self msg]:", m)
}


// read data from terminal then send to input channel
func (ui *BlockchainUi) readDataTerminal() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		// fmt.Print("terminal input> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		if sendData == "" {
			return
		}
		ui.inputChan <- sendData
	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}