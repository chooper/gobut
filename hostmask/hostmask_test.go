
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
    var bad_pattern Hostmask = "user!dontmatch@*"
    var good_pattern Hostmask = "user!*@host"

    if mask.Matches(bad_pattern) {
        t.Errorf("mask should not match bad_pattern")
    }

    if !mask.Matches(good_pattern) {
        t.Errorf("mask should match good_pattern")
    }

    if !mask.Matches("*") {
        t.Errorf("mask should match a wildcard")
    }

    if !mask.Matches(mask) {
        t.Errorf("mask should match itself")
    }
}

