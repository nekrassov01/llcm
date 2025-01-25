package llcm

import (
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/nekrassov01/mintab"
)

// OutputType represents the type of the output.
type Renderer[E Entry, D EntryData[E]] struct {
	Data       D
	OutputType OutputType
	w          io.Writer
}

// NewRenderer creates a new renderer with the specified parameters.
func NewRenderer[E Entry, D EntryData[E]](w io.Writer, data D, outputType OutputType) *Renderer[E, D] {
	return &Renderer[E, D]{
		Data:       data,
		OutputType: outputType,
		w:          w,
	}
}

// String returns the string representation of the renderer.
func (ren *Renderer[E, D]) String() string {
	b, _ := json.MarshalIndent(ren, "", "  ")
	return string(b)
}

// Render renders the output.
func (ren *Renderer[E, D]) Render() error {
	switch ren.OutputType {
	case OutputTypeJSON:
		return ren.toJSON()
	case OutputTypeText, OutputTypeCompressedText, OutputTypeMarkdown, OutputTypeBacklog:
		return ren.toTable()
	case OutputTypeTSV:
		return ren.toTSV()
	default:
		return nil
	}
}

func (ren *Renderer[E, D]) toJSON() error {
	b := json.NewEncoder(ren.w)
	b.SetIndent("", "  ")
	return b.Encode(ren.Data.Entries())
}

func (ren *Renderer[E, D]) toTable() error {
	var opt mintab.Option
	switch ren.OutputType {
	case OutputTypeText:
		opt = mintab.WithFormat(mintab.TextFormat)
	case OutputTypeCompressedText:
		opt = mintab.WithFormat(mintab.CompressedTextFormat)
	case OutputTypeMarkdown:
		opt = mintab.WithFormat(mintab.MarkdownFormat)
	case OutputTypeBacklog:
		opt = mintab.WithFormat(mintab.BacklogFormat)
	}
	table := mintab.New(ren.w, opt)
	if err := table.Load(ren.toInput()); err != nil {
		return err
	}
	table.Render()
	return nil
}

func (ren *Renderer[E, D]) toTSV() error {
	if len(ren.Data.Entries()) == 0 {
		return nil
	}
	w := csv.NewWriter(ren.w)
	w.Comma = '\t'
	if err := w.Write(ren.Data.Header()); err != nil {
		return err
	}
	for _, entry := range ren.Data.Entries() {
		if err := w.Write(entry.toTSV()); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (ren *Renderer[E, D]) toInput() mintab.Input {
	data := make([][]any, len(ren.Data.Entries()))
	for i, entry := range ren.Data.Entries() {
		data[i] = entry.toInput()
	}
	return mintab.Input{
		Header: ren.Data.Header(),
		Data:   data,
	}
}
