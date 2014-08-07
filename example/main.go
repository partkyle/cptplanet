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
	env := cptplanet.NewEnvironment("EXAMPLE_")
	env.ErrorOnExtraKeys = true

	host := env.String("HOST", "127.0.0.1", "host to bind")
	port := env.Int("PORT", 9999, "port to bind")
	debug := env.Bool("DEBUG", false, "to debug or not to debug?")
	timeout := env.Duration("TIMEOUT", 5*time.Second, "timeout duration")

	var kafkas multistring
	env.Var(&kafkas, "KAFKAS", "kafka hosts to connect to")

	err := env.Parse()
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("host=%v port=%v debug=%v timeout=%v kafkas=%v", *host, *port, *debug, *timeout, kafkas)
}
