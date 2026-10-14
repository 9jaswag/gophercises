package data

import (
	"encoding/json"
	"errors"
	"os"
)

type (
	Arcs         map[string]storyContent
	storyContent struct {
		Title      string        `json:"title,omitempty"`
		Paragraphs []string      `json:"story,omitempty"`
		Options    []storyOption `json:"options,omitempty"`
	}
	storyOption struct {
		Arc  string `json:"arc,omitempty"`
		Text string `json:"text,omitempty"`
	}
)

func ReadAndParseJson(filepath string) (Arcs, error) {
	data, err := readJsonFile(filepath)
	if err != nil {
		return Arcs{}, err
	}

	stories, err := parseJsonBytes(data)
	if err != nil {
		return Arcs{}, err
	}

	return stories, nil
}

func readJsonFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return []byte{}, errors.New("Error reading JSON file")
	}

	return data, nil
}

func parseJsonBytes(data []byte) (Arcs, error) {
	var stories Arcs

	err := json.Unmarshal(data, &stories)
	if err != nil {
		return Arcs{}, err
	}

	return stories, nil
}
