package main

import (
	"github.com/antonbiluta/blockchain-golang/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("Test1")
	chain.AddBlock("Test2")
	chain.AddBlock("Test3")

	//for _, block := range chain {
	//	fmt.Printf("Prev Hash: %x\n", block.PrevHash)
	//	fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//
	//	pow := blockchain.NewProof(block)
	//	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}
}
