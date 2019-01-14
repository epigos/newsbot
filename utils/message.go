package utils

// ButtonType for buttons, it can be ButtonTypeWebURL or ButtonTypePostback
type ButtonType string

// AttachmentType describes attachment type in GenericMessage
type AttachmentType string

// TemplateType of template in GenericMessage
type TemplateType string

// NotificationType for sent messages
type NotificationType string

// ContentType for quick replies
type ContentType string

// Message interface that represents all type of messages that we can send to Facebook Messenger
type Message interface {
	serialize()
	String() string
}

func (m TextMessage) serialize()       {} // Message interface
func (m GenericMessage) serialize()    {} // Message interface
func (m QuickReplyMessage) serialize() {} // Message interface

const (
	// ButtonTypeWebURL is type for web links
	ButtonTypeWebURL = ButtonType("web_url")

	//ButtonTypePostback is type for postback buttons that sends data back to webhook
	ButtonTypePostback = ButtonType("postback")

	// ButtonTypeShare is type for sharing of message
	ButtonTypeShare = ButtonType("element_share")

	// AttachmentTypeTemplate for template attachments
	AttachmentTypeTemplate = AttachmentType("template")

	// TemplateTypeGeneric for generic message templates
	TemplateTypeGeneric = TemplateType("generic")

	// NotificationTypeRegular for regular notification type
	NotificationTypeRegular = NotificationType("REGULAR")

	// NotificationTypeSilentPush for silent push
	NotificationTypeSilentPush = NotificationType("SILENT_PUSH")

	// NotificationTypeNoPush for no push
	NotificationTypeNoPush = NotificationType("NO_PUSH")

	// QuickReplyTextType quick reply text type
	QuickReplyTextType = ContentType("text")
)

// Recipient represents facebook recipient
type Recipient struct {
	ID string `json:"id"`
}

// TextMessage struct used for sending text messages to messenger
type TextMessage struct {
	Message          textMessageContent `json:"message"`
	Recipient        Recipient          `json:"recipient"`
	NotificationType NotificationType   `json:"notification_type,omitempty"`
}

// GenericMessage struct used for sending structural messages to messenger (messages with images, links, and buttons)
type GenericMessage struct {
	Message          genericMessageContent `json:"message"`
	Recipient        Recipient             `json:"recipient"`
	NotificationType NotificationType      `json:"notification_type,omitempty"`
}

// QuickReplyMessage struct used for sending quick replies
type QuickReplyMessage struct {
	Message   quickReplyMessageContent `json:"message"`
	Recipient Recipient                `json:"recipient"`
}

type quickReplyMessageContent struct {
	Text         string       `json:"text,omitempty"`
	QuickReplies []quickReply `json:"quick_replies"`
}

type quickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

type textMessageContent struct {
	Text string `json:"text,omitempty"`
}

type genericMessageContent struct {
	Attachment *attachment `json:"attachment,omitempty"`
}

type attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload payload `json:"payload,omitempty"`
}

type payload struct {
	TemplateType string     `json:"template_type,omitempty"`
	Sharable     bool       `json:"sharable,omitempty"`
	Elements     []*Element `json:"elements,omitempty"`
}

// Element in Generic Message template attachment
type Element struct {
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle,omitempty"`
	ItemURL  string    `json:"item_url,omitempty"`
	ImageURL string    `json:"image_url,omitempty"`
	Buttons  []*Button `json:"buttons,omitempty"`
}

// Button on Generic Message template element
type Button struct {
	Type    ButtonType `json:"type"`
	URL     string     `json:"url,omitempty"`
	Title   string     `json:"title,omitempty"`
	Payload string     `json:"payload,omitempty"`
}

func (m *TextMessage) String() string {
	return m.Message.Text
}

func (m *GenericMessage) String() string {
	return m.Message.Attachment.Type
}

func (m *QuickReplyMessage) String() string {
	return m.Message.Text
}

// AddNewElement adds element to Generic template message with defined title, subtitle, link url and image url
// Title param is mandatory. If not used set "" for other params and nil for buttons param
// Generic messages can have up to 10 elements which are scolled horizontaly in Facebook messenger
func (m *GenericMessage) AddNewElement(title, subtitle, itemURL, imageURL string, buttons []*Button) {
	m.AddElement(newElement(title, subtitle, itemURL, imageURL, buttons))
}

// AddElement adds element e to Generic Message
// Generic messages can have up to 10 elements which are scolled horizontaly in Facebook messenger
func (m *GenericMessage) AddElement(e *Element) {
	m.Message.Attachment.Payload.Elements = append(m.Message.Attachment.Payload.Elements, e)
}

// NewElement creates new element with defined title, subtitle, link url and image url
// Title param is mandatory. If not used set "" for other params and nil for buttons param
// Instead of calling this function you can also initialize Element struct, depends what you prefere
func NewElement(title, subtitle, itemURL, imageURL string, buttons []*Button) *Element {
	return newElement(title, subtitle, itemURL, imageURL, buttons)
}

func newElement(title, subtitle, itemURL, imageURL string, buttons []*Button) *Element {
	return &Element{
		Title:    title,
		Subtitle: subtitle,
		ItemURL:  itemURL,
		ImageURL: imageURL,
		Buttons:  buttons,
	}
}

// NewWebURLButton creates new web url button
func NewWebURLButton(title, URL string) *Button {
	return &Button{
		Type:  ButtonTypeWebURL,
		Title: title,
		URL:   URL,
	}
}

// NewPostbackButton creates new postback button that sends payload string back to webhook when pressed
func NewPostbackButton(title, payload string) *Button {
	return &Button{
		Type:    ButtonTypePostback,
		Title:   title,
		Payload: payload,
	}
}

// NewShareButton creates new share button that sends payload string back to webhook when pressed
func NewShareButton() *Button {
	return &Button{
		Type: ButtonTypeShare,
	}
}

// AddWebURLButton creates and adds web link URL button to the element
func (e *Element) AddWebURLButton(title, URL string) {
	b := NewWebURLButton(title, URL)
	e.Buttons = append(e.Buttons, b)
}

// AddPostbackButton creates and adds button that sends payload string back to webhook when pressed
func (e *Element) AddPostbackButton(title, payload string) {
	b := NewPostbackButton(title, payload)
	e.Buttons = append(e.Buttons, b)
}

// AddShareButton creates and adds share button that sends payload string back to webhook when pressed
func (e *Element) AddShareButton() {
	b := NewShareButton()
	e.Buttons = append(e.Buttons, b)
}

// NewTextMessage creates new text message for userID
// This function is here for convenient reason, you will
// probably use shorthand version SentTextMessage which sends message immediatly
func NewTextMessage(userID string, text string) *TextMessage {
	return &TextMessage{
		Recipient: Recipient{ID: userID},
		Message:   textMessageContent{Text: text},
	}
}

// NewGenericMessage creates new Generic Template message for userID
// Generic template messages are used for structured messages with images, links, buttons and postbacks
func NewGenericMessage(userID string) *GenericMessage {
	return &GenericMessage{
		Recipient: Recipient{ID: userID},
		Message: genericMessageContent{
			Attachment: &attachment{
				Type:    "template",
				Payload: payload{TemplateType: "generic", Sharable: true},
			},
		},
	}
}

// NewQuickReply creates a new quick reply message for userID
func NewQuickReply(userID, text string) *QuickReplyMessage {
	return &QuickReplyMessage{
		Recipient: Recipient{ID: userID},
		Message: quickReplyMessageContent{
			Text:         text,
			QuickReplies: []quickReply{},
		},
	}
}

// AddQuickReply adds quick reply to the quick reply message
func (m *QuickReplyMessage) AddQuickReply(title, payload, contentType, imageURL string) {
	reply := quickReply{
		Title:       title,
		ContentType: contentType,
		Payload:     payload,
		ImageURL:    imageURL,
	}
	m.Message.QuickReplies = append(m.Message.QuickReplies, reply)
}

// AddTextQuickReply adds text quick reply to the quick reply message
func (m *QuickReplyMessage) AddTextQuickReply(title, payload string) {
	reply := quickReply{
		Title:       title,
		Payload:     payload,
		ContentType: "text",
	}
	m.Message.QuickReplies = append(m.Message.QuickReplies, reply)
}

// NewSubscribeMenu returns new subscribe quick replies
func NewSubscribeMenu(userID string) *QuickReplyMessage {
	// add quick reply
	reply := NewQuickReply(userID, "Here are some options ⬇️")
	reply.AddTextQuickReply("Subscribe", "Subscribe")
	reply.AddTextQuickReply("Other Topics", "Topics")
	return reply
}
