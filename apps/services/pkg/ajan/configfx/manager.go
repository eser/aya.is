package configfx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotStruct                  = errors.New("not a struct")
	ErrMissingRequiredConfigValue = errors.New("missing required config value")
)

type ConfigManager struct{}

var _ ConfigLoader = (*ConfigManager)(nil)

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (cl *ConfigManager) LoadMeta(i any) (ConfigItemMeta, error) {
	r := reflect.ValueOf(i).Elem() //nolint:varnamelen

	children, err := reflectMeta(r)
	if err != nil {
		return ConfigItemMeta{}, err
	}

	return ConfigItemMeta{
		Name:            "root",
		Field:           r,
		Type:            nil,
		IsRequired:      false,
		HasDefaultValue: false,
		DefaultValue:    "",

		Children: children,
	}, nil
}

// ------------------------
// Load Methods
// ------------------------

func (cl *ConfigManager) LoadMap(resources ...ConfigResource) (*map[string]any, error) {
	target := make(map[string]any)

	for _, resource := range resources {
		err := resource(&target)
		if err != nil {
			return nil, err
		}
	}

	return &target, nil
}

func (cl *ConfigManager) Load(i any, resources ...ConfigResource) error {
	meta, err := cl.LoadMeta(i)
	if err != nil {
		return err
	}

	target, err := cl.LoadMap(resources...)
	if err != nil {
		return err
	}

	err = reflectSet(meta, "", target)
	if err != nil {
		return err
	}

	return nil
}

func (cl *ConfigManager) LoadDefaults(i any) error {
	return cl.Load(
		i,
		cl.FromJSONFile("config.json"),
		cl.FromEnvFile(".env", true),
		cl.FromSystemEnv(true),
	)
}

func reflectMeta(r reflect.Value) ([]ConfigItemMeta, error) { //nolint:varnamelen
	result := make([]ConfigItemMeta, 0)

	if r.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"%w (type=%s)",
			ErrNotStruct,
			r.Type().String(),
		)
	}

	for i := range r.NumField() {
		structField := r.Field(i)
		structFieldType := r.Type().Field(i)

		if structFieldType.Anonymous {
			children, err := reflectMeta(structField)
			if err != nil {
				return nil, err
			}

			if children != nil {
				result = append(result, children...)
			}

			continue
		}

		tag, hasTag := structFieldType.Tag.Lookup(TagConf)
		if !hasTag {
			continue
		}

		_, isRequired := structFieldType.Tag.Lookup(TagRequired)
		defaultValue, hasDefaultValue := structFieldType.Tag.Lookup(TagDefault)

		var children []ConfigItemMeta = nil

		if structFieldType.Type.Kind() == reflect.Struct {
			var err error

			children, err = reflectMeta(structField)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, ConfigItemMeta{
			Name:            tag,
			Field:           structField,
			Type:            structFieldType.Type,
			IsRequired:      isRequired,
			HasDefaultValue: hasDefaultValue,
			DefaultValue:    defaultValue,

			Children: children,
		})
	}

	return result, nil
}

func reflectSet( //nolint:cyclop,gocognit,funlen
	meta ConfigItemMeta,
	prefix string,
	target *map[string]any,
) error {
	for _, child := range meta.Children {
		key := prefix + child.Name

		if child.Type.Kind() == reflect.Map { //nolint:nestif
			// Create a new map
			newMap := reflect.MakeMap(child.Type)

			// Find all keys that start with our prefix
			prefix := key + Separator
			for targetKey := range *target {
				// if !strings.HasPrefix(targetKey, prefix) {
				if !strings.HasPrefix(strings.ToLower(targetKey), strings.ToLower(prefix)) {
					continue
				}

				// Extract the map key from the flattened key
				mapKey := targetKey[len(prefix):]
				if idx := strings.Index(mapKey, Separator); idx != -1 {
					mapKey = mapKey[:idx]
				}

				// Create and set the map value
				valueType := child.Type.Elem()
				mapValue := reflect.New(valueType).Elem()

				if valueType.Kind() == reflect.String {
					value, valueOk := (*target)[targetKey].(string)

					if valueOk {
						mapValue.SetString(value)
					}
				}

				// Recursively set the fields of the map value
				subMeta := ConfigItemMeta{
					Name:            mapKey,
					Field:           mapValue,
					Type:            valueType,
					IsRequired:      child.IsRequired,
					HasDefaultValue: child.HasDefaultValue,
					DefaultValue:    child.DefaultValue,

					Children: nil,
				}

				if valueType.Kind() == reflect.Struct {
					children, _ := reflectMeta(mapValue)
					subMeta.Children = children
				}

				err := reflectSet(subMeta, prefix+mapKey+Separator, target)
				if err != nil {
					return err
				}

				// Set the value in the map
				newMap.SetMapIndex(reflect.ValueOf(mapKey), mapValue)
			}

			child.Field.Set(newMap)

			continue
		}

		if child.Type.Kind() == reflect.Struct {
			err := reflectSet(child, key+Separator, target)
			if err != nil {
				return err
			}

			continue
		}

		// Check if the target map has the key with the child name
		value, valueOk := (*target)[key].(string)
		if !valueOk {
			if child.HasDefaultValue {
				reflectSetField(child.Field, child.Type, child.DefaultValue)

				continue
			}

			if child.IsRequired {
				return fmt.Errorf(
					"%w (key=%q, child_name=%q, child_type=%s)",
					ErrMissingRequiredConfigValue,
					key,
					child.Name,
					child.Type.String(),
				)
			}

			continue
		}

		reflectSetField(child.Field, child.Type, value)
	}

	return nil
}

func reflectSetField( //nolint:cyclop,funlen
	field reflect.Value,
	fieldType reflect.Type,
	value string,
) {
	var finalValue reflect.Value

	switch fieldType {
	case reflect.TypeFor[string]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int]():
		intValue, _ := strconv.Atoi(value)
		finalValue = reflect.ValueOf(intValue)
	case reflect.TypeFor[int8]():
		int64Value, _ := strconv.ParseInt(value, 10, 8)
		int8Value := int8(int64Value)
		finalValue = reflect.ValueOf(int8Value)
	case reflect.TypeFor[int16]():
		int64Value, _ := strconv.ParseInt(value, 10, 16)
		int16Value := int16(int64Value)
		finalValue = reflect.ValueOf(int16Value)
	case reflect.TypeFor[int32]():
		int64Value, _ := strconv.ParseInt(value, 10, 32)
		int32Value := int32(int64Value)
		finalValue = reflect.ValueOf(int32Value)
	case reflect.TypeFor[int64]():
		int64Value, _ := strconv.ParseInt(value, 10, 64)
		finalValue = reflect.ValueOf(int64Value)
	case reflect.TypeFor[uint]():
		uint64Value, _ := strconv.ParseUint(value, 10, 64)
		uintValue := uint(uint64Value)
		finalValue = reflect.ValueOf(uintValue)
	case reflect.TypeFor[uint8]():
		uint64Value, _ := strconv.ParseUint(value, 10, 8)
		uint8Value := uint8(uint64Value)
		finalValue = reflect.ValueOf(uint8Value)
	case reflect.TypeFor[uint16]():
		uint64Value, _ := strconv.ParseUint(value, 10, 16)
		uint16Value := uint16(uint64Value)
		finalValue = reflect.ValueOf(uint16Value)
	case reflect.TypeFor[uint32]():
		uint64Value, _ := strconv.ParseUint(value, 10, 32)
		uint32Value := uint32(uint64Value)
		finalValue = reflect.ValueOf(uint32Value)
	case reflect.TypeFor[uint64]():
		uint64Value, _ := strconv.ParseUint(value, 10, 64)
		finalValue = reflect.ValueOf(uint64Value)
	case reflect.TypeFor[float32]():
		floatValue, _ := strconv.ParseFloat(value, 32)
		finalValue = reflect.ValueOf(floatValue)
	case reflect.TypeFor[float64]():
		floatValue, _ := strconv.ParseFloat(value, 64)
		finalValue = reflect.ValueOf(floatValue)
	case reflect.TypeFor[bool]():
		boolValue, _ := strconv.ParseBool(value)
		finalValue = reflect.ValueOf(boolValue)
	case reflect.TypeFor[time.Duration]():
		durationValue, _ := time.ParseDuration(value)
		finalValue = reflect.ValueOf(durationValue)
	default:
		return
	}

	if field.Kind() == reflect.Ptr {
		// Handle pointer types by allocating a new instance
		ptr := reflect.New(fieldType.Elem())
		ptr.Elem().Set(finalValue)
		field.Set(ptr)

		return
	}

	// Set the field directly
	field.Set(finalValue)
}
