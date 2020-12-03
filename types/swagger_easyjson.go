// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package types

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

func easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes(in *jlexer.Lexer, out *Swagger) {
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
		case "openapi":
			out.Openapi = string(in.String())
		case "components":
			if in.IsNull() {
				in.Skip()
				out.Components = nil
			} else {
				if out.Components == nil {
					out.Components = new(Components)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Components).UnmarshalJSON(data))
				}
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
func easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes(out *jwriter.Writer, in Swagger) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"openapi\":"
		out.RawString(prefix[1:])
		out.String(string(in.Openapi))
	}
	{
		const prefix string = ",\"components\":"
		out.RawString(prefix)
		if in.Components == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.Components).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Swagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Swagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Swagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Swagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes(l, v)
}
func easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes1(in *jlexer.Lexer, out *Components) {
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
		case "schemas":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Schemas = make(map[string]*Object)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 *Object
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(Object)
						}
						if data := in.Raw(); in.Ok() {
							in.AddError((*v1).UnmarshalJSON(data))
						}
					}
					(out.Schemas)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
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
func easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes1(out *jwriter.Writer, in Components) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"schemas\":"
		out.RawString(prefix[1:])
		if in.Schemas == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Schemas {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				if v2Value == nil {
					out.RawString("null")
				} else {
					out.Raw((*v2Value).MarshalJSON())
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Components) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Components) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5ac9ea7aEncodeGithubComLutomasSwagger2MdTypes1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Components) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Components) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5ac9ea7aDecodeGithubComLutomasSwagger2MdTypes1(l, v)
}
