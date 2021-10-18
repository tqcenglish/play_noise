package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/sirupsen/logrus"

	"github.com/hajimehoshi/go-mp3"
)

var files []fs.FileInfo
var otoContext *oto.Context

func run(f fs.FileInfo) error {
	file, err := os.Open(fmt.Sprintf("./assets/%s", f.Name()))
	if err != nil {
		return err
	}
	defer file.Close()

	d, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}

	p := otoContext.NewPlayer()
	defer p.Close()

	fmt.Printf("Name: %s, Length: %d[bytes]\n", file.Name(), d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func readDir() {
	var err error
	files, err = ioutil.ReadDir("./assets")
	if err != nil {
		logrus.Fatal(err)
	}
}

func Play() {
	for {
		source := rand.NewSource(time.Now().UnixNano())
		index := rand.New(source).Intn(len(files))
		file := files[index]
		if err := run(file); err != nil {
			logrus.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	var err error
	readDir()
	otoContext, err = oto.NewContext(44100, 2, 2, 8192)
	if err != nil {
		logrus.Error(err)
	}
	defer otoContext.Close()
	// 同时播放两种
	go Play()
	Play()

}
