package pubsubevent

type TopicName string

const (
	UserRegisteredTopicName           TopicName = "account.user-registered"
	EmailConfirmationCreatedTopicName TopicName = "account.email-confirmation-created"
)

func (t TopicName) String() string {
	return string(t)
}
