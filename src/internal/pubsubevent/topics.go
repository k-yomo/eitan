package pubsubevent

type TopicName string

const (
	UserRegisteredTopicName TopicName = "account.user-registered"
)

func (t TopicName) String() string {
	return string(t)
}
