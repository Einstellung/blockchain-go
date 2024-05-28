package cli

import (
	p2p "blockchain-go/p2p"
	utils "blockchain-go/utils"
	"fmt"
	"io"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BlockchainUi struct {
	inputChan chan string
	doneCh    chan struct{}

	bcommu *p2p.BlockchainCommunication
	app    *tview.Application

	peerMsgBox io.Writer
	peerList   *tview.TextView
}

func NewBlockchainUi(bcommu *p2p.BlockchainCommunication) *BlockchainUi {
	app := tview.NewApplication()

	// make a text box to contain peer messages
	peerMsgBox := tview.NewTextView()
	peerMsgBox.SetDynamicColors(true)
	peerMsgBox.SetBorder(true)
	peerMsgBox.SetTitle(fmt.Sprintln("My peer ID:", bcommu.SelfID()))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	peerMsgBox.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel("broadcast >: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	// the done function is called when the user hits enter, or tabs out the field
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		// send the line onto the input chan and reset the field text
		inputCh <- line
		input.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peerList := tview.NewTextView()
	peerList.SetBorder(true)
	peerList.SetTitle("peers")
	peerList.SetChangedFunc(func() {
		app.Draw()
	})

	// peerPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	peerPanel := tview.NewFlex().
		AddItem(peerMsgBox, 0, 1, false).
		AddItem(peerList, 20, 1, false)

	// flex is a vertical box with the peerPanel on top and the input field at the bottom.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(peerPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	return &BlockchainUi{
		inputChan: inputCh,
		doneCh:    make(chan struct{}, 1),

		bcommu: bcommu,
		app:    app,

		peerMsgBox: peerMsgBox,
		peerList:   peerList,
	}
}

func (ui *BlockchainUi) Run() error{
	go ui.handleEvents()

	defer ui.end()

	return ui.app.Run()
}

func (ui *BlockchainUi) end() {
	ui.doneCh <- struct{}{}
}

func (ui *BlockchainUi) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

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

		case <-peerRefreshTicker.C:
			ui.refreshPeers()

		case <-ui.doneCh:
			return
		}
	}
}

func (ui *BlockchainUi) displayPeerMessage(m p2p.BlockMessage) {
	prompt := utils.WithColor("green", fmt.Sprintf("<%s>:", utils.ShortID(m.SenderID)))
	fmt.Fprintln(ui.peerMsgBox, prompt, m.Message)
}

func (ui *BlockchainUi) displaySelfMessage(m string) {
	fmt.Fprintln(ui.peerMsgBox, "[self msg]:", m)
}

func (ui *BlockchainUi) refreshPeers() {
	peers := ui.bcommu.ListPeers()

	// clear is thread-safe
	ui.peerList.Clear()

	for _, p := range peers {
		fmt.Fprintln(ui.peerList, utils.ShortID(p))
	}

	ui.app.Draw()
}
