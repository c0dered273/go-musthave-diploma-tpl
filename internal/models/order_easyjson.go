// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	decimal "github.com/shopspring/decimal"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(in *jlexer.Lexer, out *OrderDTO) {
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
		case "number":
			out.ID = uint64(in.Uint64())
		case "status":
			out.Status = OrderStatus(in.Int())
		case "accrual":
			if in.IsNull() {
				in.Skip()
				out.Amount = nil
			} else {
				if out.Amount == nil {
					out.Amount = new(decimal.Decimal)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Amount).UnmarshalJSON(data))
				}
			}
		case "uploaded_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UploadedAt).UnmarshalJSON(data))
			}
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
func easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(out *jwriter.Writer, in OrderDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"number\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Int(int(in.Status))
	}
	if in.Amount != nil {
		const prefix string = ",\"accrual\":"
		out.RawString(prefix)
		out.Raw((*in.Amount).MarshalJSON())
	}
	{
		const prefix string = ",\"uploaded_at\":"
		out.RawString(prefix)
		out.Raw((in.UploadedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OrderDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrderDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OrderDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrderDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(l, v)
}