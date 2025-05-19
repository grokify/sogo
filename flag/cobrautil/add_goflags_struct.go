package cobrautil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
	"github.com/spf13/cobra"
)

// AddFlags augments a `cobra.Command` with flags with by
// struct definitions used by `github.com/jessevdk/go-flags`.
func AddFlags(cmd *cobra.Command, opts any) error {
	if cmd == nil {
		return errors.New("cobra command cannot be nil")
	} else if opts == nil {
		return nil
	}
	val := reflect.ValueOf(opts).Elem()
	typ := val.Type()

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		flagLong := strings.TrimSpace(field.Tag.Get("long"))
		flagShort := strings.TrimSpace(field.Tag.Get("short"))
		desc := field.Tag.Get("description")
		defaultStr := field.Tag.Get("default")
		required := strings.ToLower(strings.TrimSpace(field.Tag.Get("required")))

		if strings.TrimSpace(flagLong) == "" {
			continue // skip fields without a long flag
		}

		switch field.Type.Kind() {
		case reflect.String:
			ptr := fieldVal.Addr().Interface().(*string)
			if flagShort != "" {
				cmd.Flags().StringVarP(ptr, flagLong, flagShort, defaultStr, desc)
			} else {
				cmd.Flags().StringVar(ptr, flagLong, defaultStr, desc)
			}
		case reflect.Int:
			def := 0
			if _, err := fmt.Sscanf(defaultStr, "%d", &def); err != nil {
				return err
			}
			ptr := fieldVal.Addr().Interface().(*int)
			if flagShort != "" {
				cmd.Flags().IntVarP(ptr, flagLong, flagShort, def, desc)
			} else {
				cmd.Flags().IntVar(ptr, flagLong, def, desc)
			}
		case reflect.Bool:
			def := false
			if _, err := fmt.Sscanf(defaultStr, "%t", &def); err != nil {
				return err
			}
			ptr := fieldVal.Addr().Interface().(*bool)
			if flagShort != "" {
				cmd.Flags().BoolVarP(ptr, flagLong, flagShort, def, desc)
			} else {
				cmd.Flags().BoolVar(ptr, flagLong, def, desc)
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				defSlice := []string{}
				if defaultStr != "" {
					defSlice = stringsutil.SplitTrimSpace(defaultStr, ",", true)
				}
				ptr := fieldVal.Addr().Interface().(*[]string)
				if flagShort != "" {
					cmd.Flags().StringSliceVarP(ptr, flagLong, flagShort, defSlice, desc)
				} else {
					cmd.Flags().StringSliceVar(ptr, flagLong, defSlice, desc)
				}
			}
		}
		if required == "true" {
			if err := cmd.MarkFlagRequired(flagLong); err != nil {
				return err
			}
		}
	}

	return nil
}
