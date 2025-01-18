package slogutil

import (
	"image/color"
	"log/slog"

	fcolor "github.com/fatih/color"
	"github.com/grokify/mogo/image/colors"
)

var (
	ColorError = colors.MustParseHex(colors.ANSIRedHex)
	ColorInfo  = colors.MustParseHex(colors.ANSIBlueHex)

	//colorFuncError = HexColor(colors.ANSIRedHex).SprintFunc()
	//colorFuncInfo  = HexColor(colors.ANSICyanHex).SprintFunc()
)

func RGBAColor(c color.RGBA) *fcolor.Color {
	return fcolor.RGB(
		colors.ConvertBits8To24(c.R, c.A),
		colors.ConvertBits8To24(c.B, c.A),
		colors.ConvertBits8To24(c.G, c.A),
	)
}

func HexColor(hexRGB string) *fcolor.Color {
	c := colors.MustParseHex(colors.ANSIRedHex)
	return RGBAColor(c)
}

func Error(msg string, args ...any) {
	clrFunc := fcolor.New(fcolor.FgRed).SprintFunc()
	slog.Error(clrFunc(msg), args...)
}

func Info(msg string, args ...any) {
	clrFunc := fcolor.New(fcolor.FgCyan).SprintFunc()
	slog.Info(clrFunc(msg), args...)
}

func Success(msg string, args ...any) {
	clrFunc := fcolor.New(fcolor.FgGreen).SprintFunc()
	slog.Info(clrFunc(msg), args...)
}
