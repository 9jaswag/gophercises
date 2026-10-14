package data

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
)

//go:embed gopher.json
var jsonFile embed.FS

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
	data, err := fs.ReadFile(jsonFile, filepath)
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
