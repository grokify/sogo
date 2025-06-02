package cobrautil

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CobraToGoflagsRuntime fills a go-flags-style struct with values from a cobra.Command's flags
func CobraToGoflagsRuntime(cmd *cobra.Command, opts any) error {
	flagSet := cmd.Flags()

	v := reflect.ValueOf(opts)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("opts must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	// Map field long names from `long` struct tags to field indices
	fieldMap := map[string]int{}
	for i := range t.NumField() {
		field := t.Field(i)
		longName := field.Tag.Get("long")
		if longName != "" {
			fieldMap[longName] = i
		}
	}

	// Iterate over cobra flags
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if idx, ok := fieldMap[flag.Name]; ok {
			fieldVal := v.Field(idx)
			if fieldVal.CanSet() {
				switch fieldVal.Kind() {
				case reflect.String:
					if val, err := flagSet.GetString(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetString(val)
					}
				case reflect.Bool:
					if val, err := flagSet.GetBool(flag.Name); err != nil {
						panic(err)
					} else if val {
						fieldVal.Set(reflect.ValueOf([]bool{true}))
					}
				case reflect.Int:
					if val, err := flagSet.GetInt(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetInt(int64(val))
					}
				case reflect.Int64:
					if val, err := flagSet.GetInt64(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetInt(val)
					}
				case reflect.Uint:
					if val, err := flagSet.GetUint(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetUint(uint64(val))
					}
				case reflect.Uint64:
					if val, err := flagSet.GetUint64(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetUint(val)
					}
				case reflect.Float64:
					if val, err := flagSet.GetFloat64(flag.Name); err != nil {
						panic(err)
					} else {
						fieldVal.SetFloat(val)
					}
				case reflect.Slice:
					elemKind := fieldVal.Type().Elem().Kind()
					switch elemKind {
					case reflect.String:
						if val, err := flagSet.GetStringSlice(flag.Name); err != nil {
							panic(err)
						} else {
							fieldVal.Set(reflect.ValueOf(val))
						}
					case reflect.Bool:
						// Presence of the flag means "true"
						if flag.Changed {
							fieldVal.Set(reflect.ValueOf([]bool{true}))
						}
					}
				}
			}
		}
	})

	return nil
}
