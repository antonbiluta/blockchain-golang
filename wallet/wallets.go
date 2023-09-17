package wallet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/antonbiluta/blockchain-golang/utils"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()
	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(utils.WalletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(utils.WalletFile)
	utils.HandleError(err)

	//gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	utils.HandleError(err)

	ws.Wallets = wallets.Wallets
	return nil
}

func (ws *Wallets) SaveToFile() {
	var content bytes.Buffer

	//gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	utils.HandleError(err)

	err = os.WriteFile(utils.WalletFile, content.Bytes(), 0644)
	utils.HandleError(err)
}
