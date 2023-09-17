package main

import (
	"flag"
	"fmt"
	"github.com/antonbiluta/blockchain-golang/blockchain"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Использование:")
	fmt.Println("  getbalance -address ADDRESS - получить баланс для адреса")
	fmt.Println("  createblockchain -address ADDRESS создать блокчеин и отправить genesis")
	fmt.Println("  print - Напечатать блоки из цепи")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Отправить кол-во монет")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")
	defer func(Database *badger.DB) {
		err := Database.Close()
		blockchain.HandleError(err)
	}(chain.Database)
	iter := chain.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	err := chain.Database.Close()
	if err != nil {
		return
	}
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer func(Database *badger.DB) {
		err := Database.Close()
		blockchain.HandleError(err)
	}(chain.Database)

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer func(Database *badger.DB) {
		err := Database.Close()
		blockchain.HandleError(err)
	}(chain.Database)

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
