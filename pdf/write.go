package pdf

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func WriteFileMerged(filename string, srcFilenames []string, conf *model.Configuration) error {
	return api.MergeCreateFile(
		srcFilenames,
		filename,
		false,
		conf)
}
