package cli

import (
	"fmt"
	"log"

	"blockchain-go/internal/blockchain"
)

func (cli *CLI) listAddresses() {
	wallets, err := blockchain.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
