package picker

import (
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
)

func iconColWidth(config model.Config) int {
	return icon.ColumnWidth(config)
}

func iconWidth(icn string) int {
	return icon.Width(icn)
}
