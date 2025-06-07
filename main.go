package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multibase"
)

var (
	alphabet = regexp.MustCompile("^[0-9a-z]+$")
	workers  = runtime.NumCPU()
)

const BENCHMARK_INTERVAL = 100_000

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s {suffix}\n", os.Args[0])
		os.Exit(1)
	}

	suffix := strings.ToLower(os.Args[1])
	if !alphabet.MatchString(suffix) {
		fmt.Println("Invalid characters")
		os.Exit(2)
	}

	if len(suffix) > 8 {
		fmt.Println("Suffix too long, 9+ characters would take years to find")
		os.Exit(3)
	}

	runtime.GOMAXPROCS(workers)

	privKeyChan := make(chan crypto.PrivKey)
	benchChan := make(chan int)

	for i := 0; i < workers; i++ {
		go func() {
			err := generateKey(suffix, privKeyChan, benchChan)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	go func() {
		expectedCount := math.Pow(36, float64(len(suffix)))

		start := time.Now()

		count := 0
		for range benchChan {
			count += 1

			if count%workers == 0 {
				delta := time.Now().Sub(start).Seconds()
				ops := float64(count * BENCHMARK_INTERVAL)
				opsPerSec := ops / delta
				eta := time.Duration(delta * expectedCount / ops * float64(time.Second))
				log.Printf("Benchmark: %.f keys per second - ETA: %s\n", opsPerSec, eta.Round(time.Second))
			}
		}
	}()

	<-privKeyChan
}

func generateKey(suffix string, privKeyChan chan crypto.PrivKey, benchChan chan int) error {
	for i := 1; ; i += 1 {
		privKey, pubKey, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return err
		}

		id, err := peer.IDFromPublicKey(pubKey)
		if err != nil {
			return err
		}

		keyStr, err := peer.ToCid(id).StringOfBase(multibase.Base36)
		if !strings.HasSuffix(strings.ToLower(keyStr), suffix) {
			if i%BENCHMARK_INTERVAL == 0 {
				benchChan <- 1
			}

			continue
		}

		fmt.Println(keyStr)

		privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(keyStr, privKeyBytes, 0600)
		if err != nil {
			log.Fatal(err)
		}

		privKeyChan <- privKey
		return nil
	}
}
