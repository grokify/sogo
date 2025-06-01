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
					val, _ := flagSet.GetString(flag.Name)
					fieldVal.SetString(val)
				case reflect.Bool:
					val, _ := flagSet.GetBool(flag.Name)
					fieldVal.SetBool(val)
				case reflect.Int:
					val, _ := flagSet.GetInt(flag.Name)
					fieldVal.SetInt(int64(val))
				case reflect.Int64:
					val, _ := flagSet.GetInt64(flag.Name)
					fieldVal.SetInt(val)
				case reflect.Uint:
					val, _ := flagSet.GetUint(flag.Name)
					fieldVal.SetUint(uint64(val))
				case reflect.Uint64:
					val, _ := flagSet.GetUint64(flag.Name)
					fieldVal.SetUint(val)
				case reflect.Float64:
					val, _ := flagSet.GetFloat64(flag.Name)
					fieldVal.SetFloat(val)
				case reflect.Slice:
					if fieldVal.Type().Elem().Kind() == reflect.String {
						val, _ := flagSet.GetStringSlice(flag.Name)
						fieldVal.Set(reflect.ValueOf(val))
					}
				}
			}
		}
	})

	return nil
}
