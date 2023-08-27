// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(in *jlexer.Lexer, out *OrdersDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(OrdersDTO, 0, 1)
			} else {
				*out = OrdersDTO{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 OrderDTO
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(out *jwriter.Writer, in OrdersDTO) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v OrdersDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrdersDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OrdersDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrdersDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels(l, v)
}
func easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(in *jlexer.Lexer, out *OrderDTO) {
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
			out.ID = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "accrual":
			out.Amount = float64(in.Float64())
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
func easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(out *jwriter.Writer, in OrderDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"number\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	if in.Amount != 0 {
		const prefix string = ",\"accrual\":"
		out.RawString(prefix)
		out.Float64(float64(in.Amount))
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
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrderDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OrderDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrderDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels1(l, v)
}
func easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(in *jlexer.Lexer, out *AccrualOrderDTO) {
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
		case "order":
			out.ID = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "accrual":
			out.Accrual = float64(in.Float64())
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
func easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(out *jwriter.Writer, in AccrualOrderDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"accrual\":"
		out.RawString(prefix)
		out.Float64(float64(in.Accrual))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AccrualOrderDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AccrualOrderDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson120d1ca2EncodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AccrualOrderDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AccrualOrderDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson120d1ca2DecodeGithubComC0dered273GoMusthaveDiplomaTplInternalModels2(l, v)
}
