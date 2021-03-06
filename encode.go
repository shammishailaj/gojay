package gojay

// MarshalObject returns the JSON encoding of v.
//
// It takes a struct implementing Marshaler to a JSON slice of byte
// it returns a slice of bytes and an error.
// Example with an Marshaler:
//	type TestStruct struct {
//		id int
//	}
//	func (s *TestStruct) MarshalObject(enc *gojay.Encoder) {
//		enc.AddIntKey("id", s.id)
//	}
//	func (s *TestStruct) IsNil() bool {
//		return s == nil
//	}
//
// 	func main() {
//		test := &TestStruct{
//			id: 123456,
//		}
//		b, _ := gojay.Marshal(test)
// 		fmt.Println(b) // {"id":123456}
//	}
func MarshalObject(v MarshalerObject) ([]byte, error) {
	enc := NewEncoder()
	enc.grow(200)
	enc.writeByte('{')
	v.MarshalObject(enc)
	enc.writeByte('}')
	defer enc.addToPool()
	return enc.buf, nil
}

// MarshalArray returns the JSON encoding of v.
//
// It takes an array or a slice implementing Marshaler to a JSON slice of byte
// it returns a slice of bytes and an error.
// Example with an Marshaler:
// 	type TestSlice []*TestStruct
//
// 	func (t TestSlice) MarshalArray(enc *Encoder) {
//		for _, e := range t {
//			enc.AddObject(e)
//		}
//	}
//
//	func main() {
//		test := &TestSlice{
//			&TestStruct{123456},
//			&TestStruct{7890},
// 		}
// 		b, _ := Marshal(test)
//		fmt.Println(b) // [{"id":123456},{"id":7890}]
//	}
func MarshalArray(v MarshalerArray) ([]byte, error) {
	enc := NewEncoder()
	enc.grow(200)
	enc.writeByte('[')
	v.(MarshalerArray).MarshalArray(enc)
	enc.writeByte(']')
	defer enc.addToPool()
	return enc.buf, nil
}

// Marshal returns the JSON encoding of v.
//
// Marshal takes interface v and encodes it according to its type.
// Basic example with a string:
// 	b, err := gojay.Marshal("test")
//	fmt.Println(b) // "test"
//
// If v implements Marshaler or Marshaler interface
// it will call the corresponding methods.
//
// If a struct, slice, or array is passed and does not implement these interfaces
// it will return a a non nil InvalidTypeError error.
// Example with an Marshaler:
//	type TestStruct struct {
//		id int
//	}
//	func (s *TestStruct) MarshalObject(enc *gojay.Encoder) {
//		enc.AddIntKey("id", s.id)
//	}
//	func (s *TestStruct) IsNil() bool {
//		return s == nil
//	}
//
// 	func main() {
//		test := &TestStruct{
//			id: 123456,
//		}
//		b, _ := gojay.Marshal(test)
// 		fmt.Println(b) // {"id":123456}
//	}
func Marshal(v interface{}) ([]byte, error) {
	var b []byte
	var err error = InvalidTypeError("Unknown type to Marshal")
	switch vt := v.(type) {
	case MarshalerObject:
		enc := NewEncoder()
		enc.writeByte('{')
		vt.MarshalObject(enc)
		enc.writeByte('}')
		b = enc.buf
		defer enc.addToPool()
		return b, nil
	case MarshalerArray:
		enc := NewEncoder()
		enc.writeByte('[')
		vt.MarshalArray(enc)
		enc.writeByte(']')
		b = enc.buf
		defer enc.addToPool()
		return b, nil
	case string:
		enc := NewEncoder()
		b, err = enc.encodeString(vt)
		defer enc.addToPool()
	case bool:
		enc := NewEncoder()
		err = enc.AddBool(vt)
		b = enc.buf
		defer enc.addToPool()
	case int:
		enc := NewEncoder()
		b, err = enc.encodeInt(int64(vt))
		defer enc.addToPool()
	case int64:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(vt)
	case int32:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case int16:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case int8:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case uint64:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case uint32:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case uint16:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeInt(int64(vt))
	case uint8:
		enc := NewEncoder()
		b, err = enc.encodeInt(int64(vt))
		defer enc.addToPool()
	case float64:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeFloat(vt)
	case float32:
		enc := NewEncoder()
		defer enc.addToPool()
		return enc.encodeFloat(float64(vt))
	}
	return b, err
}

// MarshalerObject is the interface to implement for struct to be encoded
type MarshalerObject interface {
	MarshalObject(enc *Encoder)
	IsNil() bool
}

// MarshalerArray is the interface to implement
// for a slice or an array to be encoded
type MarshalerArray interface {
	MarshalArray(enc *Encoder)
}

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	buf []byte
}

func (enc *Encoder) getPreviousRune() (byte, bool) {
	last := len(enc.buf) - 1
	if last < 0 {
		return 0, false
	}
	return enc.buf[last], true
}
