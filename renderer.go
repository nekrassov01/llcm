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
func NewRenderer[E Entry, D EntryData[E]](w io.Writer, data D) *Renderer[E, D] {
	return &Renderer[E, D]{
		Data:       data,
		OutputType: OutputTypeJSON,
		w:          w,
	}
}

// SetOutputType sets the output type for rendering.
func (ren *Renderer[E, D]) SetOutputType(outputType string) error {
	if outputType == "" {
		return nil
	}
	t, err := ParseOutputType(outputType)
	if err != nil {
		return err
	}
	ren.OutputType = t
	return nil
}

// String returns the string representation of the renderer.
func (ren *Renderer[E, D]) String() string {
	b, _ := json.MarshalIndent(ren, "", "  ")
	return string(b)
}

// Render renders the output.
func (ren *Renderer[E, D]) Render() error {
	switch ren.OutputType {
	case OutputTypeJSON, OutputTypePrettyJSON:
		return ren.toJSON()
	case OutputTypeText, OutputTypeCompressedText, OutputTypeMarkdown, OutputTypeBacklog:
		return ren.toTable()
	case OutputTypeTSV:
		return ren.toTSV()
	case OutputTypeChart:
		return ren.toChart()
	default:
		return nil
	}
}

func (ren *Renderer[E, D]) toJSON() error {
	b := json.NewEncoder(ren.w)
	if ren.OutputType == OutputTypePrettyJSON {
		b.SetIndent("", "  ")
	}
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
	if err := table.Load(ren.getInput()); err != nil {
		return err
	}
	table.Render()
	return nil
}

func (ren *Renderer[E, D]) toTSV() error {
	entries := ren.Data.Entries()
	if len(entries) == 0 {
		return nil
	}
	w := csv.NewWriter(ren.w)
	w.Comma = '\t'
	if err := w.Write(ren.Data.Header()); err != nil {
		return err
	}
	for _, entry := range entries {
		if err := w.Write(entry.toTSV()); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (ren *Renderer[E, D]) toChart() error {
	if len(ren.Data.Entries()) == 0 {
		return nil
	}
	return ren.Data.Chart()
}

func (ren *Renderer[E, D]) getInput() mintab.Input {
	var (
		entries = ren.Data.Entries()
		data    = make([][]any, len(entries))
	)
	if len(entries) == 0 {
		return mintab.Input{}
	}
	for i, entry := range entries {
		data[i] = entry.toInput()
	}
	return mintab.Input{
		Header: ren.Data.Header(),
		Data:   data,
	}
}
