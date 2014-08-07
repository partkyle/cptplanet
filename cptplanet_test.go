package cptplanet

import (
	"fmt"
	"os"

	"testing"
	"time"
)

func TestBasicEnvironment(t *testing.T) {
	settings := Settings{}
	env := NewEnvironment(settings)

	os.Clearenv()

	expectedString := "this is a string value"
	os.Setenv("STRING", expectedString)
	expectedInt := 101
	os.Setenv("INT", fmt.Sprintf("%d", 101))
	expectedBool := true
	os.Setenv("BOOL", fmt.Sprintf("%v", expectedBool))
	expectedDuration := 1*time.Minute + 3*time.Second
	os.Setenv("DURATION", expectedDuration.String())

	defaultString := "nothing"
	s := env.String("STRING", defaultString, "")
	defaultInt := -1
	i := env.Int("INT", defaultInt, "")
	defaultBool := false
	b := env.Bool("BOOL", defaultBool, "")
	defaultDuration := 1 * time.Second
	d := env.Duration("DURATION", defaultDuration, "")

	if *s != defaultString {
		t.Errorf("Unexpected default; have %s got %s", defaultString, *s)
	}

	if *i != defaultInt {
		t.Errorf("Unexpected default; have %s got %s", defaultInt, *i)
	}

	if *b != defaultBool {
		t.Errorf("Unexpected default; have %s got %s", defaultBool, *b)
	}

	if *d != defaultDuration {
		t.Errorf("Unexpected default; have %s got %s", defaultDuration, *d)
	}

	env.Parse()

	if *s != expectedString {
		t.Errorf("Unexpected value; have %s got %s", expectedString, *s)
	}

	if *i != expectedInt {
		t.Errorf("Unexpected value; have %s got %s", expectedInt, *i)
	}

	if *b != expectedBool {
		t.Errorf("Unexpected value; have %s got %s", expectedBool, *b)
	}

	if *d != expectedDuration {
		t.Errorf("Unexpected value; have %s got %s", expectedDuration, *d)
	}
}
