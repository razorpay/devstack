package cloud

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}

type Message struct {
	ID            string
	ReceiptHandle string
	Body          string
	Attributes    map[string]string
}

