package got

import "github.com/nlopes/slack"

type Message struct {
	got  *Got
	raw  *slack.MessageEvent
	text string
}

func (got *Got) newMessage(e *slack.MessageEvent) *Message {
	return &Message{
		got:  got,
		raw:  e,
		text: got.mentionRE.ReplaceAllString(e.Text, ""),
	}
}

func (m *Message) Channel() string {
	return m.raw.Channel
}

func (m *Message) FullText() string {
	return m.raw.Text
}

func (m *Message) IsMention() bool {
	return m.got.mentionRE.MatchString(m.raw.Text)
}

func (m *Message) RawMessageEvent() *slack.MessageEvent {
	return m.raw
}

func (m *Message) Reply(msg string) {
	m.Send("<@" + m.UserID() + "> " + msg)
}

func (m *Message) Send(msg string) {
	m.got.rtm.SendMessage(m.got.rtm.NewOutgoingMessage(msg, m.Channel()))
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) UserID() string {
	return m.raw.User
}
