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
	os.Setenv("INT", fmt.Sprintf("%d", expectedInt))
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

	err := env.Parse()

	if err != nil {
		t.Errorf("Unexpected error occurred during parsing: %s", err)
	}

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

func TestNoErrorsWhenNotUsingErrorOnExtraKeys(t *testing.T) {
	// explicit
	settings := Settings{ErrorOnExtraKeys: false}
	env := NewEnvironment(settings)

	os.Clearenv()

	expectedString := "this is a string value"
	os.Setenv("STRING", expectedString)
	expectedInt := 101
	os.Setenv("INT", fmt.Sprintf("%d", expectedInt))

	err := env.Parse()

	if err != nil {
		t.Errorf("Unexpected error; got %s", err)
	}
}

func TestErrorOnExtraKeys(t *testing.T) {
	settings := Settings{ErrorOnExtraKeys: true}
	env := NewEnvironment(settings)

	os.Clearenv()

	expectedString := "this is a string value"
	os.Setenv("STRING", expectedString)
	expectedInt := 101
	os.Setenv("INT", fmt.Sprintf("%d", expectedInt))

	err := env.Parse()

	if err == nil {
		t.Error("Expected error; got nil")
	}
}

func TestErrorOnExtraKeysAllUsed(t *testing.T) {
	settings := Settings{ErrorOnExtraKeys: true}
	env := NewEnvironment(settings)

	os.Clearenv()

	expectedString := "this is a string value"
	os.Setenv("STRING", expectedString)
	expectedInt := 101
	os.Setenv("INT", fmt.Sprintf("%d", expectedInt))

	// consume the environment
	env.String("STRING", "", "")
	env.Int("INT", 0, "")

	err := env.Parse()

	if err != nil {
		t.Error("Unexpected error; got %s", err)
	}
}
