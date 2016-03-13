package hostmask

import (
	"testing"
)

func TestNickIdentHost(t *testing.T) {
	var mask Hostmask = "user!ident@host"

	nick := mask.Nick()
	if nick != "user" {
		t.Errorf("Expected Nick to be 'user' not '%s'", nick)
	}

	ident := mask.Ident()
	if ident != "ident" {
		t.Errorf("Expected Ident to be 'ident' not '%s'", ident)
	}

	host := mask.Host()
	if host != "host" {
		t.Errorf("Expected Host to be 'host' not '%s'", host)
	}
}

func TestMatches(t *testing.T) {
	var mask Hostmask = "user!ident@host"

	var bad_patterns = []Hostmask{
		"baduser!*@*",
		"baduser!*@host",
		"baduser!ident@host",
		"*!badident@*",
		"*!badident@host",
		"user!badident@host",
		"*!*@badhost",
		"user!*@badhost",
		"user!ident@badhost",
		"baduser!badident@badhost",
	}

	var good_patterns = []Hostmask{
		"user!ident@host",
		"user!ident@*",
		"user!*@*",
		"*!*@*",
		"*!ident@*",
		"*!ident@host",
		"*!*@host",
	}

	for _, pattern := range bad_patterns {
		if mask.Matches(pattern) {
			t.Errorf("mask '%v' matched pattern '%v'", mask, pattern)
		}
		if pattern.Matches(mask) {
			t.Errorf("pattern '%v' matched mask '%v'", pattern, mask)
		}
	}

	for _, pattern := range good_patterns {
		if !mask.Matches(pattern) {
			t.Errorf("mask '%v' did not match pattern '%v'", mask, pattern)
		}
		if !pattern.Matches(mask) {
			t.Errorf("pattern '%v' did not match mask '%v'", pattern, mask)
		}
	}

	if !mask.Matches("*") {
		t.Errorf("mask should match a wildcard")
	}

	if !Hostmask("*").Matches(mask) {
		t.Errorf("wildcard should match a mask")
	}

	if !mask.Matches(mask) {
		t.Errorf("mask should match itself")
	}
}
