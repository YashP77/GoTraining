package internal

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func collectN(t *testing.T, sub Subscriber, want int) []string {
	t.Helper()

	got := make([]string, 0, want)

	timeout := time.After(2 * time.Second)

	for len(got) < want {
		select {
		case msg := <-sub:
			got = append(got, msg)
		case <-timeout:
			t.Fatalf("timeout: got %d messages, want %d", len(got), want)
		}
	}

	return got
}

func TestConcurrentPublishers(t *testing.T) {
	t.Parallel()

	StartActor()

	const (
		numPublishers = 5
		numMessages   = 10
		numSubs       = 3
	)

	// create subscribers
	subs := make([]Subscriber, 0, numSubs)
	for i := 0; i < numSubs; i++ {
		subs = append(subs, Subscribe())
		defer Unsubscribe(subs[i])
	}

	// Start publishers concurrently
	var wg sync.WaitGroup
	wg.Add(numPublishers)

	for p := 0; p < numPublishers; p++ {
		go func(pid int) {
			defer wg.Done()
			for m := 0; m < numMessages; m++ {
				Publish("pub" + strconv.Itoa(pid) + "-msg" + strconv.Itoa(m))
			}
		}(p)
	}

	// Wait for publisher to finish
	wg.Wait()

	// Each subscriber should get numPublishers * numMessages
	expected := numPublishers * numMessages

	for i, sub := range subs {
		got := collectN(t, sub, expected)
		if len(got) != expected {
			t.Fatalf("subscriber %d: received %d messages; want %d", i, len(got), expected)
		}
	}
}
