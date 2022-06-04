package notifier

type INotifier interface {
	Notify(message string) error
}

type notifier struct {
	url string
}

func NewNotifier(url string) INotifier {
	return &notifier{
		url: url,
	}
}

func (n *notifier) Notify(message string) error {
	//TODO implement me
	return nil
}
