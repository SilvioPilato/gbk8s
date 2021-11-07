package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/silviopilato/gbk8s/pkg/tasks"
)

func TestReadTaskFromJSON(t *testing.T) {
	_, calledModule, _, _ := runtime.Caller(0)
	testfolder :=filepath.Dir(calledModule)
	mockPath := filepath.Join(testfolder, "mocks", "redis.json")
	jsonFile, _ := os.Open(mockPath)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	task := tasks.ReadTaskFromJSON(byteValue)
	if (task.Properties.Image != "docker.io/library/redis") {
		t.Fatalf("Invalid image name")
	}

	if (task.Task != "START_WORKLOAD") {
		t.Fatalf("Invalid task name")
	}

	if (task.Properties.ContainerName != "test-redis") {
		t.Fatalf("Invalid container name")
	}
}