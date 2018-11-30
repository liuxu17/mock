package service

import (
	"log"
	"github.com/kaifei-bianjie/mock/sign"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/key"
)

func BatchGenSignedTxData(num int) []string {
	var (
		method       = "BatchGenSignedTx"
		signedTxData []string
	)
	resChan := make(chan types.GenSignedTxDataRes)

	senderInfos, err := key.CreateAccounts(num)
	if err != nil {
		// TODO: handle err
	}

	lens := len(senderInfos)
	log.Printf("%v: now use %v goroutine to gen signed data\n",
		method, lens)
	for i, senderInfo := range senderInfos {
		go sign.GenSignedTxData(senderInfo, conf.DefaultReceiverAddr, resChan, i)
	}

	counter := 0
	for {
		res := <-resChan
		counter ++
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

	return signedTxData
}
