// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC80ae7adDecodeYuiPluginsAnimeCharacterConfig(in *jlexer.Lexer, out *RemoteCallS) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Test":
			out.Test = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeYuiPluginsAnimeCharacterConfig(out *jwriter.Writer, in RemoteCallS) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Test\":"
		out.RawString(prefix[1:])
		out.String(string(in.Test))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoteCallS) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeYuiPluginsAnimeCharacterConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoteCallS) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeYuiPluginsAnimeCharacterConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoteCallS) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeYuiPluginsAnimeCharacterConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoteCallS) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeYuiPluginsAnimeCharacterConfig(l, v)
}
