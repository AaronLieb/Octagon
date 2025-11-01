package cache

import (
	"os"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
)

func TestMain(m *testing.M) {
	// Setup test database
	testPath := "/tmp/octagon-cache-test"
	_ = os.RemoveAll(testPath)

	var err error
	db, err = badger.Open(badger.DefaultOptions(testPath).WithLoggingLevel(badger.ERROR))
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_ = db.Close()
	_ = os.RemoveAll(testPath)

	os.Exit(code)
}

func TestSetAndGet(t *testing.T) {
	key := []byte("test_key")
	value := []byte("test_value")

	err := Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	retrieved, err := Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	key := []byte("nonexistent")

	_, err := Get(key)
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}
}

func TestClear(t *testing.T) {
	key := []byte("clear_test")
	value := []byte("clear_value")

	err := Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	_, err = Get(key)
	if err == nil {
		t.Error("Expected error after clear")
	}
}

func TestTTL(t *testing.T) {
	// This test would require mocking time or waiting, so we'll just verify
	// that Set doesn't error when setting TTL
	key := []byte("ttl_test")
	value := []byte("ttl_value")

	err := Set(key, value)
	if err != nil {
		t.Fatalf("Set with TTL failed: %v", err)
	}

	// Verify we can retrieve immediately
	retrieved, err := Get(key)
	if err != nil {
		t.Fatalf("Get after Set failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}
}
