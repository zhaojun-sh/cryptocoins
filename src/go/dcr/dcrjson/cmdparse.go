// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcrjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func makeParams(rt reflect.Type, rv reflect.Value) []interface{} {
	numFields := rt.NumField()
	params := make([]interface{}, 0, numFields)
	for i := 0; i < numFields; i++ {
		rtf := rt.Field(i)
		rvf := rv.Field(i)
		if rtf.Type.Kind() == reflect.Ptr {
			if rvf.IsNil() {
				break
			}
			rvf.Elem()
		}
		params = append(params, rvf.Interface())
	}

	return params
}

func MarshalCmd(rpcVersion string, id interface{}, cmd interface{}) ([]byte, error) {
	rt := reflect.TypeOf(cmd)
	registerLock.RLock()
	method, ok := concreteTypeToMethod[rt]
	registerLock.RUnlock()
	if !ok {
		str := fmt.Sprintf("%q is not registered", method)
		return nil, makeError(ErrUnregisteredMethod, str)
	}

	rv := reflect.ValueOf(cmd)
	if rv.IsNil() {
		str := "the specified command is nil"
		return nil, makeError(ErrInvalidType, str)
	}

	params := makeParams(rt.Elem(), rv.Elem())

	rawCmd, err := NewRequest(rpcVersion, id, method, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(rawCmd)
}

func checkNumParams(numParams int, info *methodInfo) error {
	if numParams < info.numReqParams || numParams > info.maxParams {
		if info.numReqParams == info.maxParams {
			str := fmt.Sprintf("wrong number of params (expected "+
				"%d, received %d)", info.numReqParams,
				numParams)
			return makeError(ErrNumParams, str)
		}

		str := fmt.Sprintf("wrong number of params (expected "+
			"between %d and %d, received %d)", info.numReqParams,
			info.maxParams, numParams)
		return makeError(ErrNumParams, str)
	}

	return nil
}

func populateDefaults(numParams int, info *methodInfo, rv reflect.Value) {
	// When there are no more parameters left in the supplied parameters,
	// any remaining struct fields must be optional.  Thus, populate them
	// with their associated default value as needed.
	for i := numParams; i < info.maxParams; i++ {
		rvf := rv.Field(i)
		if defaultVal, ok := info.defaults[i]; ok {
			rvf.Set(defaultVal)
		}
	}
}

func UnmarshalCmd(r *Request) (interface{}, error) {
	registerLock.RLock()
	rtp, ok := methodToConcreteType[r.Method]
	info := methodToInfo[r.Method]
	registerLock.RUnlock()
	if !ok {
		str := fmt.Sprintf("%q is not registered", r.Method)
		return nil, makeError(ErrUnregisteredMethod, str)
	}
	rt := rtp.Elem()
	rvp := reflect.New(rt)
	rv := rvp.Elem()

	numParams := len(r.Params)
	if err := checkNumParams(numParams, &info); err != nil {
		return nil, err
	}

	for i := 0; i < numParams; i++ {
		rvf := rv.Field(i)
		concreteVal := rvf.Addr().Interface()
		if err := json.Unmarshal(r.Params[i], &concreteVal); err != nil {
			fieldName := strings.ToLower(rt.Field(i).Name)
			if jerr, ok := err.(*json.UnmarshalTypeError); ok {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"be type %v (got %v)", i+1, fieldName,
					jerr.Type, jerr.Value)
				return nil, makeError(ErrInvalidType, str)
			}

			str := fmt.Sprintf("parameter #%d '%s' failed to "+
				"unmarshal: %v", i+1, fieldName, err)
			return nil, makeError(ErrInvalidType, str)
		}
	}

	if numParams < info.maxParams {
		populateDefaults(numParams, &info, rv)
	}

	return rvp.Interface(), nil
}

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Float32, reflect.Float64:

		return true
	}

	return false
}

func typesMaybeCompatible(dest reflect.Type, src reflect.Type) bool {
	if dest == src {
		return true
	}

	srcKind := src.Kind()
	destKind := dest.Kind()
	if isNumeric(destKind) && isNumeric(srcKind) {
		return true
	}

	if srcKind == reflect.String {
		if isNumeric(destKind) {
			return true
		}

		switch destKind {
		case reflect.Bool:
			return true

		case reflect.String:
			return true

		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			return true
		}
	}

	return false
}

func baseType(arg reflect.Type) (reflect.Type, int) {
	var numIndirects int
	for arg.Kind() == reflect.Ptr {
		arg = arg.Elem()
		numIndirects++
	}
	return arg, numIndirects
}

func assignField(paramNum int, fieldName string, dest reflect.Value, src reflect.Value) error {
	destBaseType, destIndirects := baseType(dest.Type())
	srcBaseType, srcIndirects := baseType(src.Type())
	if !typesMaybeCompatible(destBaseType, srcBaseType) {
		str := fmt.Sprintf("parameter #%d '%s' must be type %v (got "+
			"%v)", paramNum, fieldName, destBaseType, srcBaseType)
		return makeError(ErrInvalidType, str)
	}

	if destBaseType == srcBaseType && srcIndirects >= destIndirects {
		for i := 0; i < srcIndirects-destIndirects; i++ {
			src = src.Elem()
		}
		dest.Set(src)
		return nil
	}

	destIndirectsRemaining := destIndirects
	if destIndirects > srcIndirects {
		indirectDiff := destIndirects - srcIndirects
		for i := 0; i < indirectDiff; i++ {
			dest.Set(reflect.New(dest.Type().Elem()))
			dest = dest.Elem()
			destIndirectsRemaining--
		}
	}

	if destBaseType == srcBaseType {
		dest.Set(src)
		return nil
	}

	for i := 0; i < destIndirectsRemaining; i++ {
		dest.Set(reflect.New(dest.Type().Elem()))
		dest = dest.Elem()
	}

	for src.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:

		switch dest.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:

			srcInt := src.Int()
			if dest.OverflowInt(srcInt) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}

			dest.SetInt(srcInt)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:

			srcInt := src.Int()
			if srcInt < 0 || dest.OverflowUint(uint64(srcInt)) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetUint(uint64(srcInt))

		default:
			str := fmt.Sprintf("parameter #%d '%s' must be type "+
				"%v (got %v)", paramNum, fieldName, destBaseType,
				srcBaseType)
			return makeError(ErrInvalidType, str)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:

		switch dest.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:

			srcUint := src.Uint()
			if srcUint > uint64(1<<63)-1 {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			if dest.OverflowInt(int64(srcUint)) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetInt(int64(srcUint))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:

			srcUint := src.Uint()
			if dest.OverflowUint(srcUint) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetUint(srcUint)

		default:
			str := fmt.Sprintf("parameter #%d '%s' must be type "+
				"%v (got %v)", paramNum, fieldName, destBaseType,
				srcBaseType)
			return makeError(ErrInvalidType, str)
		}

	case reflect.Float32, reflect.Float64:
		destKind := dest.Kind()
		if destKind != reflect.Float32 && destKind != reflect.Float64 {
			str := fmt.Sprintf("parameter #%d '%s' must be type "+
				"%v (got %v)", paramNum, fieldName, destBaseType,
				srcBaseType)
			return makeError(ErrInvalidType, str)
		}

		srcFloat := src.Float()
		if dest.OverflowFloat(srcFloat) {
			str := fmt.Sprintf("parameter #%d '%s' overflows "+
				"destination type %v", paramNum, fieldName,
				destBaseType)
			return makeError(ErrInvalidType, str)
		}
		dest.SetFloat(srcFloat)

	case reflect.String:
		switch dest.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(src.String())
			if err != nil {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"parse to a %v", paramNum, fieldName,
					destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetBool(b)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:

			srcInt, err := strconv.ParseInt(src.String(), 0, 0)
			if err != nil {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"parse to a %v", paramNum, fieldName,
					destBaseType)
				return makeError(ErrInvalidType, str)
			}
			if dest.OverflowInt(srcInt) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetInt(srcInt)

		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:

			srcUint, err := strconv.ParseUint(src.String(), 0, 0)
			if err != nil {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"parse to a %v", paramNum, fieldName,
					destBaseType)
				return makeError(ErrInvalidType, str)
			}
			if dest.OverflowUint(srcUint) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetUint(srcUint)

		case reflect.Float32, reflect.Float64:
			srcFloat, err := strconv.ParseFloat(src.String(), 0)
			if err != nil {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"parse to a %v", paramNum, fieldName,
					destBaseType)
				return makeError(ErrInvalidType, str)
			}
			if dest.OverflowFloat(srcFloat) {
				str := fmt.Sprintf("parameter #%d '%s' "+
					"overflows destination type %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.SetFloat(srcFloat)

		case reflect.String:
			dest.SetString(src.String())

		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			concreteVal := dest.Addr().Interface()
			err := json.Unmarshal([]byte(src.String()), &concreteVal)
			if err != nil {
				str := fmt.Sprintf("parameter #%d '%s' must "+
					"be valid JSON which unmarshals to a %v",
					paramNum, fieldName, destBaseType)
				return makeError(ErrInvalidType, str)
			}
			dest.Set(reflect.ValueOf(concreteVal).Elem())
		}
	}

	return nil
}

func NewCmd(method string, args ...interface{}) (interface{}, error) {
	registerLock.RLock()
	rtp, ok := methodToConcreteType[method]
	info := methodToInfo[method]
	registerLock.RUnlock()
	if !ok {
		str := fmt.Sprintf("%q is not registered", method)
		return nil, makeError(ErrUnregisteredMethod, str)
	}

	numParams := len(args)
	if err := checkNumParams(numParams, &info); err != nil {
		return nil, err
	}

	rvp := reflect.New(rtp.Elem())
	rv := rvp.Elem()
	rt := rtp.Elem()

	for i := 0; i < numParams; i++ {
		rvf := rv.Field(i)
		fieldName := strings.ToLower(rt.Field(i).Name)
		err := assignField(i+1, fieldName, rvf, reflect.ValueOf(args[i]))
		if err != nil {
			return nil, err
		}
	}

	return rvp.Interface(), nil
}
