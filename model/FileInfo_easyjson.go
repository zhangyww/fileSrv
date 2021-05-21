// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

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

func easyjson35847e5aDecodeFileSrvModel(in *jlexer.Lexer, out *FileInfo) {
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
		case "FilePath":
			out.FilePath = string(in.String())
		case "Op":
			out.Op = FileOp(in.Int())
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
func easyjson35847e5aEncodeFileSrvModel(out *jwriter.Writer, in FileInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"FilePath\":"
		out.RawString(prefix[1:])
		out.String(string(in.FilePath))
	}
	{
		const prefix string = ",\"Op\":"
		out.RawString(prefix)
		out.Int(int(in.Op))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FileInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson35847e5aEncodeFileSrvModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FileInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson35847e5aEncodeFileSrvModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FileInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson35847e5aDecodeFileSrvModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FileInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson35847e5aDecodeFileSrvModel(l, v)
}