package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/partkyle/cptplanet"
)

type multistring []string

func (m *multistring) Set(s string) error {
	*m = strings.Split(s, ",")
	return nil
}

func (m *multistring) String() string {
	return fmt.Sprintf(strings.Join(*m, ","))
}

func main() {
	host := cptplanet.String("HOST", "127.0.0.1", "host to bind")
	port := cptplanet.Int("PORT", 9999, "port to bind")
	debug := cptplanet.Bool("DEBUG", false, "to debug or not to debug?")
	timeout := cptplanet.Duration("TIMEOUT", 5*time.Second, "timeout duration")

	var kafkas multistring
	cptplanet.Var(&kafkas, "KAFKAS", "kafka hosts to connect to")

	err := cptplanet.Parse()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("host=%v port=%v debug=%v timeout=%v kafkas=%v", *host, *port, *debug, *timeout, kafkas)
}
