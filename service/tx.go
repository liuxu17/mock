package service

import (
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/key"
	"github.com/kaifei-bianjie/mock/sign"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"log"
	"time"
)

func BatchGenSignedTxData(num int, subFaucets []conf.SubFaucet) []string {
	var (
		method       = "BatchGenSignedTx"
		signedTxData []string
	)
	resChan := make(chan types.GenSignedTxDataRes, 10000)

	senderInfos, err := key.NewAccount(num, subFaucets)
	if err != nil {
		// TODO: handle err
	}

	lens := len(senderInfos)
	if lens > 0 {
		log.Printf("%v: now use %v goroutine to gen signed data\n",
			method, lens)

		for i, senderInfo := range senderInfos {
			go sign.GenSignedTxData(senderInfo, conf.DefaultReceiverAddr, resChan, i)
			if (i % 20) == 0 {
				time.Sleep(time.Second * constants.SignTxDelaySec)
			}
		}

		counter := 0
		for {
			res := <-resChan
			counter++
			if res.Res != "" {
				log.Printf("%v: successed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
				signedTxData = append(signedTxData, res.Res)
			} else {
				log.Printf("%v: failed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
			}

			if counter == lens {
				log.Printf("%v: all sign tx goroutine over\n", method)
				break
			}
		}
	} else {
		log.Printf("%v: no signed tx data\n", method)
	}

	return signedTxData
}

func SingleBatchGenSignedTxData(num int, faucet string, password string, home string, faucetAddr string) []string {
	resChan := make(chan types.GenSignedTxDataRes, 100000)
	faucetAddr, err := account.GetAccAddr(faucet, home)
	if err != nil {
		log.Printf("%v: cannot get %v address fail: %v !!!!!!!!!!!!\n",
			"New Account", faucet, err)
		return nil
	}
	senderInfos, err := key.NewAccountSingle(num, home, faucet, faucetAddr)
	if err != nil {
		// TODO: handle err
	}

	lens := len(senderInfos)
	var (
		method       = "BatchGenSignedTx"
		signedTxData [100000]string
	)
	if lens > 0 {
		log.Printf("%v: now use %v goroutine to gen signed data\n",
			method, lens)

		for i, senderInfo := range senderInfos {
			go sign.GenSignedTxDataFromSingleFaucet(faucetAddr, senderInfo, conf.DefaultReceiverAddr, resChan, i)
			if (i % 10) == 0 {
				time.Sleep(time.Second * constants.SignTxDelaySec)
			}
		}

		counter := 0
		for {
			res := <-resChan
			counter++
			if res.Res != "" {
				log.Printf("%v: successed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
				signedTxData[res.ChanNum] = res.Res
				//signedTxData = append(signedTxData, res.Res)
			} else {
				log.Printf("%v: failed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
			}

			if counter == lens {
				log.Printf("%v: all sign tx goroutine over\n", method)
				break
			}
		}
	} else {
		log.Printf("%v: no signed tx data\n", method)
	}

	var signedTxDataReturn []string
	for j := 0; j < lens; j++ {
		signedTxDataReturn = append(signedTxDataReturn, signedTxData[j])
	}
	return signedTxDataReturn
}
