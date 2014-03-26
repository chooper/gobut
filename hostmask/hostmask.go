
package hostmask

import (
    "strings"
    "path/filepath"
)

# A hostmask, typically the "prefix" portion of an IRC message
# originating from a client, looks like "user!ident@host". It's
# commonplace to use standard globbing for matching.
# RFC 2812 SS 2.3 Messages https://tools.ietf.org/html/rfc2812

type Hostmask string

func (h Hostmask) Nick() string {
    return strings.Split(string(h), "!")[0]
}

func (h Hostmask) Ident() string {
    prefix := strings.Split(string(h), "@")[0]
    return strings.Split(prefix, "!")[1]
}

func (h Hostmask) Host() string {
    return strings.Split(string(h), "@")[1]
}

func (h Hostmask) Matches(t Hostmask) bool {
    forward_match, _ := filepath.Match(string(h), string(t))
    backward_match, _ := filepath.Match(string(t), string(h))
    return forward_match || backward_match
}

