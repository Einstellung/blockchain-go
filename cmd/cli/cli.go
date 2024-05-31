package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type CLI struct{}

// Run parses command line arguments and processes commands
func (cli *CLI) Run(input string) {
	args := strings.Split(input, " ")

	switch args[0] {

	case "getbalance":
		cli.getBalance(args[2])

	case "createblockchain":
		cli.createBlockchain(args[2])

	case "createwallet":
		cli.createWallet()

	case "listaddresses":
		cli.listAddresses()
	case "send":
		num, err := strconv.Atoi(args[6])
		if err != nil {
			fmt.Println(err)
		} else {
			cli.send(args[2], args[4], num)
		}

	case "printchain":
		cli.printChain()

	default:
		PrintUsage()
	}
}

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
	fmt.Println()
}