package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// Map implements a low-level get/set config that is backed by an in-memory tree of yaml
// nodes. It allows us to interact with a yaml-based config programmatically, preserving any
// comments that were present when the yaml was parsed.
type Map struct {
	Root *yaml.Node
}

// Entry represents a single entry in a Map.
type Entry struct {
	KeyNode   *yaml.Node
	ValueNode *yaml.Node
	Index     int
}

// NotFoundError is returned when a key is not found in a Map.
type NotFoundError struct {
	error
}

// Empty returns true if the Map is empty.
func (cm *Map) Empty() bool {
	return cm.Root == nil || len(cm.Root.Content) == 0
}

// GetStringValue returns the value of the given key as a string.
func (cm *Map) GetStringValue(key string) (string, error) {
	entry, err := cm.FindEntry(key)
	if err != nil {
		return "", err
	}
	return entry.ValueNode.Value, nil
}

// SetStringValue sets the value of the given key to the given string.
func (cm *Map) SetStringValue(key, value string) error {
	entry, err := cm.FindEntry(key)
	if err == nil {
		entry.ValueNode.Value = value
		return nil
	}

	var notFound *NotFoundError
	if err != nil && !errors.As(err, &notFound) {
		return err
	}

	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
	}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}

	cm.Root.Content = append(cm.Root.Content, keyNode, valueNode)
	return nil
}

// FindEntry returns the Entry for the given key.
func (cm *Map) FindEntry(key string) (*Entry, error) {
	ce := &Entry{}

	if cm.Empty() {
		return ce, &NotFoundError{errors.New("not found")}
	}

	// Content slice goes [key1, value1, key2, value2, ...].
	topLevelPairs := cm.Root.Content
	for i, v := range topLevelPairs {
		// Skip every other slice item since we only want to check against keys.
		if i%2 != 0 {
			continue
		}
		if v.Value == key {
			ce.KeyNode = v
			ce.Index = i
			if i+1 < len(topLevelPairs) {
				ce.ValueNode = topLevelPairs[i+1]
			}
			return ce, nil
		}
	}

	return ce, &NotFoundError{errors.New("not found")}
}

// RemoveEntry removes the entry with the given key.
func (cm *Map) RemoveEntry(key string) {
	if cm.Empty() {
		return
	}

	newContent := []*yaml.Node{}

	var skipNext bool
	for i, v := range cm.Root.Content {
		if skipNext {
			skipNext = false
			continue
		}
		if i%2 != 0 || v.Value != key {
			newContent = append(newContent, v)
		} else {
			// Don't append current node and skip the next which is this key's value.
			skipNext = true
		}
	}

	cm.Root.Content = newContent
}
