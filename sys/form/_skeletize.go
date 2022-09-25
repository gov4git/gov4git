package form

import (
	"context"
	"reflect"
)

func Skeletize(ctx context.Context, v any) any {
	return skeletize(ctx, reflect.ValueOf(v)).Interface()
}

func skeletize(ctx context.Context, v reflect.Value) reflect.Value {
	if f, ok := v.Interface().(Form); ok {
		return reflect.ValueOf(f.Skeletize(ctx))
	}
	if f, ok := v.Addr().Interface().(Form); ok {
		return reflect.ValueOf(f.Skeletize(ctx))
	}
	switch v.Kind() {
	case reflect.Bool:
		return v
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v
	case reflect.Uintptr:
		panic("type not supported")
	case reflect.Float32, reflect.Float64:
		return v
	case reflect.Complex64, reflect.Complex128:
		panic("type not supported")
	case reflect.Array:
		return skeletizeSliceOrArray(ctx, v)
	case reflect.Chan:
		panic("type not supported")
	case reflect.Func:
		panic("type not supported")
	case reflect.Interface:
		return skeletizeInterface(ctx, v)
	case reflect.Map:
		return skeletizeMap(ctx, v)
	case reflect.Pointer:
		return skeletizePointer(ctx, v)
	case reflect.Slice:
		return skeletizeSliceOrArray(ctx, v)
	case reflect.String:
		return v
	case reflect.Struct:
		return skeletizeStruct(ctx, v)
	default:
	}
	panic("type invalid or unknown")
}

var (
	typeOfMap   = reflect.TypeOf(map[string]any{})
	typeOfSlice = reflect.TypeOf([]any{})
)

const skeletizeFieldTag = "gitty"

func skeletizeSliceOrArray(ctx context.Context, v reflect.Value) reflect.Value {
	s := reflect.MakeSlice(typeOfSlice, v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		// skip nils
		w := v.Index(i)
		if (w.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) && v.IsNil() {
			continue
		}
		s.Index(i).Set(skeletize(ctx, w))
	}
	return s
}

func skeletizeStruct(ctx context.Context, v reflect.Value) reflect.Value {
	m := reflect.MakeMap(typeOfMap)
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		nm := t.Field(i).Tag.Get(skeletizeFieldTag)
		if nm == "" || nm == "-" {
			continue
		}
		w := skeletize(ctx, v.Field(i))
		// skip nils
		if (w.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) && v.IsNil() {
			continue
		}
		m.SetMapIndex(reflect.ValueOf(nm), w)
	}
	return m
}

var anyNil = reflect.ValueOf((any)(nil))

func skeletizePointer(ctx context.Context, v reflect.Value) reflect.Value {
	if v.IsNil() {
		return anyNil
	}
	return skeletize(ctx, v.Elem())
}

func skeletizeMap(ctx context.Context, v reflect.Value) reflect.Value {
	if v.Type().Key().Kind() != reflect.String {
		panic("skeletal maps must have string keys")
	}
	m := reflect.MakeMap(typeOfMap)
	for _, k := range v.MapKeys() {
		w := skeletize(ctx, v.MapIndex(k))
		// skip nils
		if (w.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) && v.IsNil() {
			continue
		}
		m.SetMapIndex(k, w)
	}
	return m
}

func skeletizeInterface(ctx context.Context, v reflect.Value) reflect.Value {
	if v.IsNil() {
		return anyNil
	}
	return skeletize(ctx, v.Elem())
}
