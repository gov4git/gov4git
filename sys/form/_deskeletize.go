package form

import (
	"context"
	"errors"
	"reflect"
)

func Deskeletize(ctx context.Context, v any, skeleton any) error {
	return deskeletize(ctx, reflect.ValueOf(v), reflect.ValueOf(skeleton))
}

func deskeletize(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if f, ok := v.Interface().(Form); ok {
		return f.Deskeletize(ctx, skeleton)
	}
	if f, ok := v.Addr().Interface().(Form); ok {
		return f.Deskeletize(ctx, skeleton)
	}
	switch v.Kind() {
	case reflect.Bool:
		return deskeletizePrimitive(ctx, v, skeleton)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return deskeletizePrimitive(ctx, v, skeleton)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return deskeletizePrimitive(ctx, v, skeleton)
	case reflect.Uintptr:
		panic("type not supported")
	case reflect.Float32, reflect.Float64:
		return deskeletizePrimitive(ctx, v, skeleton)
	case reflect.Complex64, reflect.Complex128:
		panic("type not supported")
	case reflect.Array:
		return deskeletizeArray(ctx, v, skeleton)
	case reflect.Chan:
		panic("type not supported")
	case reflect.Func:
		panic("type not supported")
	case reflect.Interface:
		return deskeletizeInterface(ctx, v, skeleton)
	case reflect.Map:
		return deskeletizeMap(ctx, v, skeleton)
	case reflect.Pointer:
		return deskeletizePointer(ctx, v, skeleton)
	case reflect.Slice:
		return deskeletizeSlice(ctx, v, skeleton)
	case reflect.String:
		return deskeletizePrimitive(ctx, v, skeleton)
	case reflect.Struct:
		return deskeletizeStruct(ctx, v, skeleton)
	default:
	}
	panic("type invalid or unknown")
}

var ErrDeskeletize = errors.New("invalid format")

func deskeletizePrimitive(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if !skeleton.Type().AssignableTo(v.Type()) {
		return ErrDeskeletize
	}
	v.Set(skeleton)
	return nil
}

func deskeletizeArray(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if skeleton.Kind() != reflect.Slice || skeleton.Len() != v.Len() {
		return ErrDeskeletize
	}
	for i := 0; i < v.Len(); i++ {
		err := deskeletize(ctx, v.Index(i), skeleton.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func deskeletizeSlice(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if skeleton.Kind() != reflect.Slice {
		return ErrDeskeletize
	}
	w := reflect.MakeSlice(v.Type(), skeleton.Len(), skeleton.Len())
	for i := 0; i < skeleton.Len(); i++ {
		if err := deskeletize(ctx, w.Index(i), skeleton.Index(i)); err != nil {
			return err
		}
	}
	v.Set(w)
	return nil
}

func deskeletizeStruct(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if skeleton.Kind() != reflect.Map {
		return ErrDeskeletize
	}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		nm := t.Field(i).Tag.Get(skeletizeFieldTag)
		if nm == "" || nm == "-" {
			continue
		}
		subskel := skeleton.MapIndex(reflect.ValueOf(nm))
		if err := deskeletize(ctx, v.Field(i), subskel); err != nil {
			return err
		}
	}
	return nil
}

func deskeletizeMap(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	if v.Type().Key().Kind() != reflect.String {
		panic("only strings can be map keys")
	}
	w := reflect.MakeMap(reflect.TypeOf(v))
	for _, k := range skeleton.MapKeys() {
		mk := reflect.New(v.Type().Key()).Elem()
		mv := reflect.New(v.Type().Elem()).Elem()
		mk.Set(k)
		if err := deskeletize(ctx, mv, skeleton.MapIndex(k)); err != nil {
			return err
		}
		w.MapIndex(mk).Set(mv)
	}
	v.Set(w)
	return nil
}

func deskeletizePointer(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	XXX
}

func deskeletizeInterface(ctx context.Context, v reflect.Value, skeleton reflect.Value) error {
	panic("deskeletizing into an interface is not supported")
}
