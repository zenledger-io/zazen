package log

import "go.uber.org/zap"

type Tag struct {
	Key   string
	Value any
}

func NewTag(k string, v any) Tag {
	return Tag{
		Key:   k,
		Value: v,
	}
}

func (t Tag) field() zap.Field {
	return zap.Any(t.Key, t.Value)
}
