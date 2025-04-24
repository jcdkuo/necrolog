package necrolog

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestLogWrite(t *testing.T) {
	testPath := "/var/log/uah_log/test_log.log"
	msg := "test logging output"
	Info(testPath, msg)

	file, err := os.Open(testPath)
	if err != nil {
		t.Fatalf("failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), msg) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("log message not found in file")
	}
}
