package pubsubevent

import "github.com/k-yomo/pm/middleware/pm_effectively_once"

func NewDeduplicateKey(topic TopicName, key string) string {
	return string(topic) + ":" + key
}

func SetDeduplicateKey(attributes map[string]string, deduplicateKey string) map[string]string {
	if attributes == nil {
		attributes = make(map[string]string, 1)
	}
	attributes[pm_effectively_once.DefaultDeduplicateKey] = deduplicateKey

	return attributes
}
