package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var authKeysModTime time.Time

func main() {
	filePath := flag.String("f", "authorized_keys", "path to the authorized_keys file")
	keyPath := flag.String("k", "test-key.pub", "public key to add to authorized_keys file (-f)")
	key := readPublicKey(*keyPath)
	flag.Parse()

	log.Println("File path:", *filePath)
	ticker := time.NewTicker(2 * time.Second)
	for {
		<-ticker.C
		log.Println("Ticked!")
		stat, err := os.Stat(*filePath)
		if err != nil {
			log.Panic("Could not get stat of file: ", err)
		}
		if stat.ModTime().After(authKeysModTime) {
			authKeysModTime = stat.ModTime()
			log.Println("FILE IS MODIFIED AFTER LAST CHECK")
			content, err := ioutil.ReadFile(*filePath)
			if err != nil {
				log.Panic("Could not read file content", err)
			}
			if !bytes.Contains(content, key) {
				file, err := os.Create(*filePath)
				if err != nil {
					log.Panic("Could not open file for writing: ", err)
				}
				content = bytes.TrimSpace(content)
				if len(content) > 0 {
					content = append(content, '\n', '\n')
				}
				content = append(content, key...)
				_, err = file.Write(content)
				if err != nil {
					log.Panic("Could not write key to file: ", err)
				}
				authKeysModTime = time.Now()
			}
		}
	}
}

func readPublicKey(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("Can't read public key file:", err)
	}
	return content
}
