package eths

import (
	"context"
	"fmt"
	"go-echo/configs"
	"go-echo/utils"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// voteAsset ...
type voteAsset struct {
	tokenID int64
	Count   int64
}

// NewAcc ...
func NewAcc(pass, connstr string) (string, error) {
	cli, err := rpc.Dial(connstr)
	if err != nil {
		fmt.Println("failed to connect to geth", err)
		return "", err
	}
	defer cli.Close()
	var account string
	err = cli.Call(&account, "personal_newAccount", pass)
	if err != nil {
		fmt.Println("failed to connect to personal_newAccount", err)
		return "", err
	}
	fmt.Println("account build successfully")
	return account, err
}

// Upload ...
func Upload(from, pass, hash, name string, age int64, record bool) error {
	cli, err := ethclient.Dial(configs.Config.Eth.Connstr)
	if err != nil {
		fmt.Println("failed to ethclient.Dial", err)
		return err
	}
	instance, err := NewHotel(common.HexToAddress(configs.Config.Eth.HotelAddr), cli)
	if err != nil {
		fmt.Println("failed to eths.NewHotel", err)
		return err
	}
	// 设置签名, owner的keyStore文件
	// 需要获得文件名字
	fileName, err := utils.GetFileName(string([]rune(from)[2:]), configs.Config.Eth.Keydir)

	file, err := os.Open(configs.Config.Eth.Keydir + "/" + fileName)
	if err != nil {
		fmt.Println("failed to os.Open", err)
		return err
	}
	auth, err := bind.NewTransactor(file, pass)
	if err != nil {
		fmt.Println("failed to bind.NewTransactor", err)
		return err
	}
	// string -> [32]byte
	// Mint(opts *bind.TransactOpts, _idcord [32]byte, _name string, _age *big.Int, _bool bool)
	_, err = instance.Mint(auth, common.HexToHash(hash), name, big.NewInt(age), record)
	if err != nil {
		fmt.Println("failed to instance.Mint", err)
		return err
	}
	fmt.Printf("the account: %s Mint success...\n", from)
	return nil
}

// EventSubscribeTest ...
func EventSubscribeTest(connstr, contractAddr string) error {
	// 1.连接ws://localhost:8546
	cli, err := ethclient.Dial(connstr)
	if err != nil {
		fmt.Println("failed to ethclient.Dial", err)
		return err
	}
	// 2. 合约地址处理
	cAddress := common.HexToAddress(contractAddr)
	newAssetHash := crypto.Keccak256Hash([]byte("onNewAsset(bytes32,address,uint256)"))
	// 3. 过滤处理
	query := ethereum.FilterQuery{
		Addresses: []common.Address{cAddress},
		Topics:    [][]common.Hash{{newAssetHash}},
	}
	// 4. 通道
	pxaLogs := make(chan types.Log)
	// 5. 订阅
	sub, err := cli.SubscribeFilterLogs(context.Background(), query, pxaLogs)
	if err != nil {
		fmt.Println("failed to cli.SubscribeFilterLogs", err)
		return err
	}
	// 6. 订阅返回处理
	fmt.Println("starting operate sub...")
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("get sub err", err)
		case vLog := <-pxaLogs:
			data, err := vLog.MarshalJSON()
			fmt.Println(string(data), err)
			ParseMintEventDb([]byte(common.Bytes2Hex(vLog.Data)))
		}
	}
}
