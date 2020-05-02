package server

type message struct {
	Version      string   `json:"version"`
	Response     response `json:"response"`
	IsSessionEnd bool     `json:"is_session_end"`
}
type response struct {
	OpenMic bool    `json:"open_mic"`
	ToSpeak toSpeak `json:"to_speak"`
}
type toSpeak struct {
	TextType int    `json:"type"`
	Text     string `json:"text"`
}

func defaultMessage() message {
	toSpeak := toSpeak{
		TextType: 0,
		Text:     "Hello",
	}
	response := response{
		OpenMic: true,
		ToSpeak: toSpeak,
	}
	message := message{
		Version:      "1.0",
		Response:     response,
		IsSessionEnd: false,
	}
	return message
}
func (message *message) setTalkContent(content string) {
	message.Response.ToSpeak.Text = content
}
func (message *message) openMic(open bool){
	message.Response.OpenMic = open
}
type receive struct {
	Version string      `json:"version"`
	Session interface{} `json:"session"`
	Request interface{} `json:"request"`
	Query   string      `json:"query"`
	Context interface{} `json:"context"`
}
