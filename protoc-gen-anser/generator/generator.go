package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// proto package is introduced; the generated code references
// a constant, proto.ProtoPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 1

// Generator is the type whose methods generate the output, stored in the associated response structure.
type Generator struct {
	Request  *plugin.CodeGeneratorRequest  // The input.
	Response *plugin.CodeGeneratorResponse // The output.

	Param map[string]string // Command-line parameters.
}

// New creates a new generator and allocates the request and response protobufs.
func New() *Generator {
	g := new(Generator)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)
	return g
}

// Error reports a problem, including an error, and exits the program.
func (g *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc-gen-anser: error:", s)
	os.Exit(1)
}

// Fail reports a problem and exits the program.
func (g *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("protoc-gen-anser: error:", s)
	os.Exit(1)
}

// CommandLineParameters breaks the comma-separated list of key=value pairs
// in the parameter (a member of the request protobuf) into a key/value map.
// It then sets file name mappings defined by those entries.
func (g *Generator) CommandLineParameters(parameter string) {
	g.Param = make(map[string]string)
	for _, p := range strings.Split(parameter, ",") {
		if i := strings.Index(p, "="); i < 0 {
			g.Param[p] = ""
		} else {
			g.Param[p[0:i]] = p[i+1:]
		}
	}
}

func (g *Generator) GenerateAllFiles() {
	m := map[string]interface{}{}
	if len(g.Request.ProtoFile) != 1 {
		g.Fail(fmt.Sprintf("only support gen 1 proto file, give %d", len(g.Request.ProtoFile)))
	}
	file := g.Request.ProtoFile[0]
	if len(file.Service) != 1 {
		g.Fail(fmt.Sprintf("only support gen 1 service in file %s, give %d", file.GetName(), len(file.Service)))
	}
	service := file.Service[0]
	LowerServiceName := strings.ToLower(service.GetName())
	m["ServiceName"] = service.GetName()
	m["LowerServiceName"] = LowerServiceName
	services := []interface{}{}
	for _, method := range service.Method {
		s := map[string]interface{}{}
		s["MethodName"] = method.GetName()
		s["InputTypeName"] = lastName(method.GetInputType())
		s["OutputTypeName"] = lastName(method.GetOutputType())
		services = append(services, s)
	}
	m["Services"] = services

	prefix := "/Users/dworld/go/src/github.com/anseres/anser/protoc-gen-anser/generator/tmpl"
	data, err := ioutil.ReadFile(prefix + "/client.go.tmpl")
	if err != nil {
		g.Error(err, "read client.go.tmpl error")
	}
	t := template.Must(template.New("client.go.tmpl").Parse(string(data)))
	buf := &bytes.Buffer{}
	err = t.Execute(buf, m)
	if err != nil {
		g.Error(err, "execute client.go.tmpl error")
	}
	g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(LowerServiceName + "_client.go"),
		Content: proto.String(buf.String()),
	})
}

// And now lots of helper functions.

var escapeChars = [256]byte{
	'a': '\a', 'b': '\b', 'f': '\f', 'n': '\n', 'r': '\r', 't': '\t', 'v': '\v', '\\': '\\', '"': '"', '\'': '\'', '?': '?',
}

// unescape reverses the "C" escaping that protoc does for default values of bytes fields.
// It is best effort in that it effectively ignores malformed input. Seemingly invalid escape
// sequences are conveyed, unmodified, into the decoded result.
func unescape(s string) string {
	// NB: Sadly, we can't use strconv.Unquote because protoc will escape both
	// single and double quotes, but strconv.Unquote only allows one or the
	// other (based on actual surrounding quotes of its input argument).

	var out []byte
	for len(s) > 0 {
		// regular character, or too short to be valid escape
		if s[0] != '\\' || len(s) < 2 {
			out = append(out, s[0])
			s = s[1:]
		} else if c := escapeChars[s[1]]; c != 0 {
			// escape sequence
			out = append(out, c)
			s = s[2:]
		} else if s[1] == 'x' || s[1] == 'X' {
			// hex escape, e.g. "\x80
			if len(s) < 4 {
				// too short to be valid
				out = append(out, s[:2]...)
				s = s[2:]
				continue
			}
			v, err := strconv.ParseUint(s[2:4], 16, 8)
			if err != nil {
				out = append(out, s[:4]...)
			} else {
				out = append(out, byte(v))
			}
			s = s[4:]
		} else if '0' <= s[1] && s[1] <= '7' {
			// octal escape, can vary from 1 to 3 octal digits; e.g., "\0" "\40" or "\164"
			// so consume up to 2 more bytes or up to end-of-string
			n := len(s[1:]) - len(strings.TrimLeft(s[1:], "01234567"))
			if n > 3 {
				n = 3
			}
			v, err := strconv.ParseUint(s[1:1+n], 8, 8)
			if err != nil {
				out = append(out, s[:1+n]...)
			} else {
				out = append(out, byte(v))
			}
			s = s[1+n:]
		} else {
			// bad escape, just propagate the slash as-is
			out = append(out, s[0])
			s = s[1:]
		}
	}

	return string(out)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// CamelCaseSlice is like CamelCase, but the argument is a slice of strings to
// be joined with "_".
func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

// dottedSlice turns a sliced name into a dotted name.
func dottedSlice(elem []string) string { return strings.Join(elem, ".") }

// Is this field optional?
func isOptional(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_OPTIONAL
}

// Is this field required?
func isRequired(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REQUIRED
}

// Is this field repeated?
func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

// Is this field a scalar numeric type?
func isScalar(field *descriptor.FieldDescriptorProto) bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return true
	default:
		return false
	}
}

// badToUnderscore is the mapping function used to generate Go names from package names,
// which can be dotted in the input .proto file.  It replaces non-identifier characters such as
// dot or dash with underscore.
func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

// baseName returns the last path element of the name, with the last dotted suffix removed.
func baseName(name string) string {
	// First, find the last element
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	// Now drop the suffix
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[0:i]
	}
	return name
}

func lastName(name string) string {
	// First, find the last element
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// The SourceCodeInfo message describes the location of elements of a parsed
// .proto file by way of a "path", which is a sequence of integers that
// describe the route from a FileDescriptorProto to the relevant submessage.
// The path alternates between a field number of a repeated field, and an index
// into that repeated field. The constants below define the field numbers that
// are used.
//
// See descriptor.proto for more information about this.
const (
	// tag numbers in FileDescriptorProto
	packagePath = 2 // package
	messagePath = 4 // message_type
	enumPath    = 5 // enum_type
	// tag numbers in DescriptorProto
	messageFieldPath   = 2 // field
	messageMessagePath = 3 // nested_type
	messageEnumPath    = 4 // enum_type
	messageOneofPath   = 8 // oneof_decl
	// tag numbers in EnumDescriptorProto
	enumValuePath = 2 // value
)
