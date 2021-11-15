package tasks

import (
	"encoding/json"
)

type Task struct {
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

func SerializeTask(task *Task) []byte {
	data, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	return data
}

func GetStartWorkloadTask(properties TaskProperties) Task {
	return Task{Properties: properties}
}

func GetRemoveWorkloadTask(properties TaskProperties) Task {
	return Task{Properties: properties}
}