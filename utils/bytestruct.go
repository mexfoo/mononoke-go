package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(reader io.Reader, order binary.ByteOrder, v interface{}, version int) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	storedValues := make(map[string]reflect.Value)

	if err := readData(reader, order, reflect.StructField{}, val, storedValues, version); err != nil {
		return err
	}

	return nil
}

//nolint:gocognit,gocyclo,cyclop,funlen // yes, too complex to understand...
func readData(reader io.Reader, order binary.ByteOrder,
	structField reflect.StructField, val reflect.Value, storedValues map[string]reflect.Value, version int) error {
	if val.Kind() != reflect.Struct {
		storedValues[structField.Name] = val
	}

	if value, ok := structField.Tag.Lookup("version"); ok {
		ver1, ver2 := getVersionAsIntFromTag(value)
		if version < ver1 || version > ver2 {
			return nil // Skip.
		}
	}

	checkKind := val.Kind()
	if kind, kindOk := structField.Tag.Lookup("subtype"); kindOk { //nolint:nestif // Fine.
		subVersion, versOk := structField.Tag.Lookup("subversion")
		if !versOk {
			return fmt.Errorf("cannot find subversion for field %s", structField.Name)
		}
		ver1, ver2 := getVersionAsIntFromTag(subVersion)
		if version >= ver1 && version <= ver2 {
			if subTypeNum, convErr := strconv.Atoi(kind); convErr == nil {
				checkKind = reflect.Kind(subTypeNum) //nolint:gosec // This is fine.
			} else {
				return convErr
			}
		}
	}

	switch checkKind { //nolint:exhaustive // too many to handle
	case reflect.Struct:
		// We always enter here first since we want to unmarshall a struct.
		t := val.Type()
		for i := range val.NumField() {
			structF := t.Field(i)
			if v := val.Field(i); v.CanSet() {
				if err := readData(reader, order, structF, v, storedValues, version); err != nil {
					return err
				}
			}
		}

	case reflect.Bool:
		var value bool
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetBool(value)
		} else {
			return err
		}
	case reflect.Int:
		var value int
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int8:
		var value int8
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int16:
		var value int16
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int32:
		var value int32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int64:
		var value int64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(value)
		} else {
			return err
		}
	case reflect.Uint:
		var value uint
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint8:
		var value uint8
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint16:
		var value uint16
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint32:
		var value uint32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint64:
		var value uint64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(value)
		} else {
			return err
		}
	case reflect.Float32:
		var value float32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetFloat(float64(value))
		} else {
			return err
		}
	case reflect.Float64:
		var value float64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetFloat(float64(value))
		} else {
			return err
		}
	case reflect.String:
		if err := unmarshalString(reader, order, storedValues, structField, val); err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		if err := unmarshalArray(reader, order, storedValues, structField, val, version); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type: %s", val.Kind().String())
	}

	return nil
}

//nolint:gocognit,funlen // yes, too complex to understand...
func unmarshalArray(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value,
	field reflect.StructField, value reflect.Value, version int) error {
	var arrayLen int

	if v, ok := field.Tag.Lookup("byteSize"); ok {
		if mapValue, exists := storedValues[v]; exists {
			arrayLen = int(mapValue.Uint()) //nolint:gosec // This is fine.
		}
	}

	// If we have multiple versions, check for a loop tag to loop (X := range loop)
	// checks if we have lenX and versionX, then compares those
	if loop, loopOk := field.Tag.Lookup("loop"); loopOk { //nolint:nestif // This is fine.
		loopNum, err := strconv.Atoi(loop)
		if err != nil {
			return fmt.Errorf("field: %s loop tag: %s", field.Name, err.Error())
		}

		for i := range loopNum {
			lenField, lenOk := field.Tag.Lookup(fmt.Sprintf("len%d", i))
			if !lenOk {
				return fmt.Errorf("field: %s no len%d tag found", field.Name, i)
			}
			length, atoiErr := strconv.Atoi(lenField)
			if atoiErr != nil {
				if lenField, exists := storedValues[lenField]; exists {
					arrayLen = int(lenField.Uint()) //nolint:gosec // This is fine.
				} else {
					return fmt.Errorf("field: %s len%d tag: %s", field.Name, i, atoiErr.Error())
				}
			}
			verField, versOk := field.Tag.Lookup(fmt.Sprintf("version%d", i))
			if !versOk {
				return fmt.Errorf("field: %s no version%d tag found", field.Name, i)
			}
			ver1, ver2 := getVersionAsIntFromTag(verField)
			if version >= ver1 && version <= ver2 {
				arrayLen = length
				break
			}
		}
	}

	// If we have an array instead of a slice
	if arrayLen == 0 && value.Cap() >= 0 {
		arrayLen = value.Cap()
	}

	if arrayLen <= 0 {
		return errors.New("no valid string length found")
	}

	// Do we need to loop value based on struct?
	data := make([]byte, arrayLen)

	typ := field.Type.Elem().Kind()
	switch typ { //nolint:exhaustive // too many to handle
	case reflect.String:
		return fmt.Errorf("does not support type with array: %s ", field.Type.Elem().Kind())

	case reflect.Struct:
		if err := binary.Read(reader, order, &data); err != nil {
			return err
		}

		slice := reflect.MakeSlice(value.Type(), arrayLen, arrayLen)
		sliceReader := bytes.NewBuffer(data)
		index := 0
		for ; sliceReader.Len() != 0; index++ {
			sliceStruct := reflect.New(value.Type().Elem()).Elem()
			t := value.Type().Elem()
			for i := range t.NumField() {
				structF := t.Field(i)
				if v := sliceStruct.Field(i); v.CanSet() {
					if err := readData(sliceReader, order, structF, v, storedValues, version); err != nil {
						return err
					}
				}
			}
			v := slice.Index(index)
			v.Set(sliceStruct)
		}
		value.Set(slice.Slice(0, index))
	default:
		if err := binary.Read(reader, order, &data); err == nil {
			if value.Kind() == reflect.Array {
				reflect.Copy(value, reflect.ValueOf(data))
			} else {
				value.SetBytes(data)
			}
		}
	}
	return nil
}

func unmarshalString(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value,
	field reflect.StructField, value reflect.Value) error {
	if v, ok := field.Tag.Lookup("byteSize"); ok {
		var size int

		if mapValue, exists := storedValues[v]; exists {
			size = int(mapValue.Uint()) //nolint:gosec // This is fine
		}

		data := make([]byte, size)
		if err := binary.Read(reader, order, &data); err == nil {
			value.SetString(string(data))
		}
	} else {
		return errors.New("missing byte tag")
	}

	return nil
}

func hexToInt(hex string) (int, error) {
	hex = strings.ReplaceAll(hex, "0x", "")
	n, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func getVersionAsIntFromTag(versionString string) (int, int) {
	if strings.HasPrefix(versionString, ".") {
		ver2, err := hexToInt(versionString[1:])
		if err != nil {
			panic(err.Error())
		}
		return 0, ver2
	}
	// e.g. 0x090500:, equals "version >= 0x090500"
	if strings.HasSuffix(versionString, ".") {
		ver1, err := hexToInt(versionString[:len(versionString)-1])
		if err != nil {
			panic(err.Error())
		}
		return ver1, 0
	}
	// e.g. 0x090500:0x090800, equals "version >= 0x090500 && version <= 0x090800"
	if strings.Contains(versionString, ".") {
		indexOf := strings.Index(versionString, ".")
		ver1, err := hexToInt(versionString[:indexOf])
		if err != nil {
			panic(err.Error())
		}

		ver2, err := hexToInt(versionString[indexOf+1:])
		if err != nil {
			panic(err.Error())
		}
		return ver1, ver2
	}

	// e.g. 0x090500, equals "version == 0x090500"
	ver, err := hexToInt(versionString)
	if err != nil {
		panic(err.Error())
	}
	return ver, ver
}

func Marshal(order binary.ByteOrder, v interface{}, version int) ([]byte, error) {
	valueOfField := reflect.ValueOf(v)
	typeOfField := reflect.TypeOf(v)

	var buf bytes.Buffer

	for i := range valueOfField.NumField() {
		// Skip
		// Has dynamic size due to packet version.
		/*if _, ok := t.Tag.Lookup("dynamic"); ok {
			size := getSizeForTag(t, version)
			value := f.Slice(0, size)
			if err := binary.Write(&buf, order, value); err != nil {
				return nil, err
			}
			continue
		}*/

		f := valueOfField.Field(i)
		t := typeOfField.Field(i)

		checkKind := t.Type.Kind()
		if kind, kindOk := t.Tag.Lookup("subtype"); kindOk { //nolint:nestif // Fine.
			subVersion, versOk := t.Tag.Lookup("subversion")
			if !versOk {
				return nil, fmt.Errorf("cannot find subversion for field %s", t.Name)
			}
			ver1, ver2 := getVersionAsIntFromTag(subVersion)
			if version >= ver1 && version <= ver2 {
				if subTypeNum, convErr := strconv.Atoi(kind); convErr == nil {
					checkKind = reflect.Kind(subTypeNum) //nolint:gosec // This is fine.
				} else {
					return nil, convErr
				}
			}
		}

		data, err := writeByteData(order, checkKind, f, version)
		if err != nil {
			return nil, err
		}
		if err = binary.Write(&buf, order, data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

//nolint:gocognit,gocyclo,cyclop,funlen // yes, too complex to understand...
func writeByteData(order binary.ByteOrder, kind reflect.Kind, f reflect.Value, version int) ([]byte, error) {
	var buf bytes.Buffer
	switch kind { //nolint:exhaustive // too many to handle
	case reflect.Bool:
		if err := binary.Write(&buf, order, f.Bool()); err != nil {
			return nil, err
		}
	case reflect.Int:
		if err := binary.Write(&buf, order, f.Int()); err != nil {
			return nil, err
		}
	case reflect.Int8:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, int8(f.Int())); err != nil {
			return nil, err
		}
	case reflect.Int16:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, int16(f.Int())); err != nil {
			return nil, err
		}
	case reflect.Int32:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, int32(f.Int())); err != nil {
			return nil, err
		}
	case reflect.Int64:
		if err := binary.Write(&buf, order, f.Int()); err != nil {
			return nil, err
		}
	case reflect.Uint:
		if err := binary.Write(&buf, order, f.Uint()); err != nil {
			return nil, err
		}
	case reflect.Uint8:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, uint8(f.Uint())); err != nil {
			return nil, err
		}
	case reflect.Uint16:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, uint16(f.Uint())); err != nil {
			return nil, err
		}
	case reflect.Uint32:
		//nolint:gosec // This has to be this way
		if err := binary.Write(&buf, order, uint32(f.Uint())); err != nil {
			return nil, err
		}
	case reflect.Uint64:
		if err := binary.Write(&buf, order, f.Uint()); err != nil {
			return nil, err
		}
	case reflect.Float32:
		if err := binary.Write(&buf, order, float32(f.Float())); err != nil {
			return nil, err
		}
	case reflect.Float64:
		if err := binary.Write(&buf, order, float64(f.Float())); err != nil {
			return nil, err
		}
	case reflect.String:
		if err := binary.Write(&buf, order, []byte(f.String())); err != nil {
			return nil, err
		}
	case reflect.Slice, reflect.Array:
		for i := range f.Len() {
			if reflect.String == f.Index(i).Kind() {
				return nil, errors.New("cant handle string slices")
			}

			data, err := writeByteData(order, f.Index(i).Kind(), f.Index(i), version)

			if err != nil {
				return nil, err
			}
			if err = binary.Write(&buf, order, data); err != nil {
				return nil, err
			}
		}

	case reflect.Struct:
		data, err := Marshal(order, f.Interface(), version)

		if err != nil {
			return nil, err
		}
		if err = binary.Write(&buf, order, data); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("does not support type : %s ", kind)
	}

	return buf.Bytes(), nil
}
