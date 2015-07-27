package sarama

import (
	"testing"
	"time"
)

// The latest committed offset saved by one partition manager instance is
// returned by another as the initial commit.
func TestFuncOffsetManagerSimple(t *testing.T) {
	setupFunctionalTest(t)
	defer teardownFunctionalTest(t)

	// Given
	newOffset := time.Now().Unix()

	config := NewConfig()
	client, err := NewClient(kafkaBrokers, config)
	if err != nil {
		t.Error(err)
	}
	offsetMgr, err := NewOffsetManagerFromClient(client)
	if err != nil {
		t.Error(err)
	}
	pom0_1, err := offsetMgr.ManagePartition("test", "test.4", 0)
	if err != nil {
		t.Error(err)
	}

	// When: several offsets are committed.
	pom0_1.CommitOffset(newOffset, "foo")
	pom0_1.CommitOffset(newOffset+1, "bar")
	pom0_1.CommitOffset(newOffset+2, "bazz")

	// Then: last committed request is the one that becomes effective.
	pom0_1.Close()
	pom0_2, err := offsetMgr.ManagePartition("test", "test.4", 0)
	if err != nil {
		t.Error(err)
	}

	fo := <-pom0_2.InitialOffset()

	if (fo != FetchedOffset{newOffset + 2, "bazz"}) {
		t.Errorf("Unexpected offset: %v", fo)
	}

	offsetMgr.Close()
}
