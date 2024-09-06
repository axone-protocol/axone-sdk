package template

import (
	"time"

	"github.com/axone-protocol/axone-sdk/credential"
)

type Option[T credential.Descriptor] func(descriptor T)

type HasID interface {
	setID(string)
}

func WithID[T interface {
	HasID
	credential.Descriptor
}](id string) Option[T] {
	return func(descriptor T) {
		descriptor.setID(id)
	}
}

type HasDatasetDID interface {
	setDatasetDID(string)
}

func WithDatasetDID[T interface {
	HasDatasetDID
	credential.Descriptor
}](did string) Option[T] {
	return func(descriptor T) {
		descriptor.setDatasetDID(did)
	}
}

type HasDescription interface {
	setDescription(string)
}

func WithDescription[T interface {
	HasDescription
	credential.Descriptor
}](description string) Option[T] {
	return func(descriptor T) {
		descriptor.setDescription(description)
	}
}

type HasFormat interface {
	setFormat(string)
}

func WithFormat[T interface {
	HasFormat
	credential.Descriptor
}](format string) Option[T] {
	return func(descriptor T) {
		descriptor.setFormat(format)
	}
}

type HasTags interface {
	setTags([]string)
}

func WithTags[T interface {
	HasTags
	credential.Descriptor
}](tags []string) Option[T] {
	return func(descriptor T) {
		descriptor.setTags(tags)
	}
}

type HasTopic interface {
	setTopic(string)
}

func WithTopic[T interface {
	HasTopic
	credential.Descriptor
}](topic string) Option[T] {
	return func(descriptor T) {
		descriptor.setTopic(topic)
	}
}

type HasIssuanceDate interface {
	setIssuanceDate(time.Time)
}

func WithIssuanceDate[T interface {
	HasIssuanceDate
	credential.Descriptor
}](t time.Time) Option[T] {
	return func(descriptor T) {
		descriptor.setIssuanceDate(t)
	}
}
