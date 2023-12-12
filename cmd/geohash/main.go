package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/quantonganh/geohash"
)

const usage = `Usage of geohash:
  -d, --decode string
        Geohash for decoding
`

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("error getting FileInfo structure: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var input string
		for scanner.Scan() {
			input = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading stdin:", err)
		}

		if len(os.Args) == 1 {
			hash, err := encode(input)
			if err != nil {
				return fmt.Errorf("error encoding: %w", err)
			}
			fmt.Println(hash)
		} else {
			switch os.Args[1] {
			case "-d", "--decode":
				lat, lng, err := decode(input)
				if err != nil {
					return fmt.Errorf("error decoding: %w", err)
				}
				fmt.Printf("%.04f, %.04f\n", lat, lng)
			default:
				fmt.Print(usage)
			}
		}
	} else {
		var decodeVal string
		flag.StringVar(&decodeVal, "decode", "", "Geohash for decoding")
		flag.StringVar(&decodeVal, "d", "", "Alias for --decode")
		flag.Usage = func() { fmt.Print(usage) }
		flag.Parse()

		if len(os.Args) == 1 {
			fmt.Print(usage)
		} else {
			switch os.Args[1] {
			case "-d", "--decode":
				lat, lng, err := decode(decodeVal)
				if err != nil {
					return fmt.Errorf("error decoding: %w", err)
				}
				fmt.Printf("%.04f, %.04f\n", lat, lng)
			default:
				hash, err := encode(os.Args[1])
				if err != nil {
					return fmt.Errorf("error encoding: %w", err)
				}
				fmt.Println(hash)
			}
		}
	}

	return nil
}

func encode(coords string) (string, error) {
	lat, lng, err := geohash.ParseCoordinate(coords)
	if err != nil {
		return "", err
	}

	return geohash.Encode(lat, lng), nil
}

func decode(hash string) (lat float64, lng float64, err error) {
	if err = geohash.ParseGeohash(hash); err != nil {
		return
	}

	lat, lng = geohash.Decode(hash)
	return
}
