package tasks

import (
	"encoding/json"
)

type Task struct {
    Task string `json:"task"`
    Properties TaskProperties `json:"properties"`
}

type TaskProperties struct {
	Image string `json:"image"`
	ContainerName string `json:"containerName"`
}

func ReadTaskFromJSON(data[] byte) Task {
	var task Task

	err := json.Unmarshal(data, &task)
	if err != nil {
		panic(err)
	}
	return task
}
