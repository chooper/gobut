package broadcast

import (
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	// Create testing harness
	got := []Message{}
	listeners := []MessageChannel{}
	output := make(chan Message)
	done := make(chan bool)

	// Initialize broadcaster
	b := NewBroadcaster(1024)

	// Start listener

	// Create 10 listeners, send results to shared `output` channel
	for x := 0; x < 10; x++ {
		new_listener := b.Listen(1024)
		listeners = append(listeners, new_listener)
		t.Logf("Added listener %d", x)

		go func(channel MessageChannel) {
			msg := <-channel
			t.Logf("Received data %s", msg)
			output <- msg
		}(new_listener)
	}

	// Collect all of the results from listeners
	go func() {
		for msg := range output {
			t.Logf("Received output %s", msg)
			got = append(got, msg) 
			if len(got) == 10 {
				// Since all 10 messages were received, done indicates
				// the test has passed.
				done <- true
			}
		}
	}()

	// Send that sample message
	t.Logf("Writing data: 1")
	b.Write(1)

	// Wait for results to show up
	select {
	case <- done:
		break
	case <- time.After(1 * time.Second):
		t.Fatalf("Messages were not received by listeners.")
	}
}
