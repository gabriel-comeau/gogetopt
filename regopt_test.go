package gogetopt

import (
	"testing"
)

// Test all valid register use cases.  Clears the options between registrations
func TestRegisterValid(t *testing.T) {
	ClearAll()
	var err error

	// both a short and long, boolean, not required
	err = RegisterOpt("valid", "valid", "v", true, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a short, boolean, not required
	err = RegisterOpt("valid", "", "v", true, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a short, boolean, not required
	err = RegisterOpt("valid", "valid", "", true, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// both a short and long, non-boolean, non-required
	err = RegisterOpt("valid", "valid", "v", false, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a short, non-boolean, non-required
	err = RegisterOpt("valid", "", "v", false, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a long, non-boolean, non-required
	err = RegisterOpt("valid", "valid", "v", false, false, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// both a short and a long, non-boolean, required
	err = RegisterOpt("valid", "valid", "v", false, true, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a short, non-boolean, required
	err = RegisterOpt("valid", "", "v", false, true, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}

	ClearAll()

	// only a long, non-boolean, required
	err = RegisterOpt("valid", "valid", "", false, true, "test usage string")
	if err != nil {
		t.Errorf("RegisterOpt() failed: %v", err)
	}
}

// Tests invalid register user cases
func TestRegisterInvalid(t *testing.T) {

	var err error

	// Fail because no long or short key provided
	err = RegisterOpt("invalid", "", "", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered without either a short or long key")
	}

	ClearAll()
	err = nil

	// Fail because too long short key provided
	err = RegisterOpt("invalid", "", "toolong", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered with a short key longer than 1 character")
	}

	ClearAll()
	err = nil

	// Fail because too short long key provided
	err = RegisterOpt("invalid", "l", "", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered with a long key shorter than 2 characters")
	}

	ClearAll()
	err = nil

	// Fail because option set to both boolean and required
	err = RegisterOpt("invalid", "l", "", true, true, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered as both boolean and required")
	}

	ClearAll()
	err = nil

	// Now we're going to register a valid option, not clear it and attempt
	// to register options with the same keys
	err = RegisterOpt("sovalid", "soval", "s", true, false, "test usage")
	if err != nil {
		t.Error("RegisterOpt(): Oops this one should have passed")
	}

	// Should fail on same key
	err = RegisterOpt("sovalid", "meha", "m", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered with the same key as another option")
	}

	err = nil

	// Should fail on same long val
	err = RegisterOpt("newsovalid", "soval", "x", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered with the same long key as a previous one")
	}

	err = nil

	// Should fail on same short val
	err = RegisterOpt("newestsovalid", "wow", "s", true, false, "test usage")
	if err == nil {
		t.Error("RegisterOpt(): An option was registered with the same short key as a previous one")
	}

	err = nil

	ClearAll()
}
