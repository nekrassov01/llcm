package llcm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/pkg/browser"
)

var (
	// PageTitle is the title of the HTML page.
	PageTitle = "llcm"

	// BaseName is the base name of the HTML file.
	BaseName = "llcm"

	// MaxPieChartItems is the maximum number of items in a pie chart.
	MaxPieChartItems = 11

	// MaxBarChartItems is the maximum number of items in a bar chart.
	MaxBarChartItems = 31

	// PieChartTitle is the title of the pie chart.
	PieChartTitle = "Stored bytes of log groups"

	// BarChartTitle is the title of the bar chart.
	BarChartTitle = "The simulation of reductions in log groups"
)

func render(chart components.Charter) error {
	var (
		fname = fmt.Sprintf("%s.html", BaseName)
		i     = 1
	)
	for {
		if _, err := os.Stat(fname); err != nil {
			if os.IsNotExist(err) {
				break
			}
			return err
		}
		fname = fmt.Sprintf("%s%d.html", BaseName, i)
		i++
	}
	f, err := os.Create(filepath.Clean(fname))
	if err != nil {
		return err
	}
	page := components.NewPage()
	page.SetPageTitle(PageTitle)
	page.AddCharts(chart)
	if err := page.Render(io.MultiWriter(f)); err != nil {
		return err
	}
	_ = browser.OpenFile(fname)
	return nil
}

func getPieItems[E Entry](entries []E) []opts.PieData {
	if len(entries) == 0 {
		return nil
	}
	var (
		othersTotal int64
		items       = make([]opts.PieData, 0, MaxPieChartItems)
	)
	for i, entry := range entries {
		m := entry.DataSet()
		v := m[storedBytesLabel]
		if v == 0 {
			continue
		}
		if i < MaxPieChartItems-1 {
			item := opts.PieData{
				Name:  entry.Name(),
				Value: v,
			}
			items = append(items, item)
		} else {
			othersTotal += v
		}
	}
	if othersTotal > 0 {
		item := opts.PieData{
			Name:  "others",
			Value: othersTotal,
		}
		items = append(items, item)
	}
	return items
}

func newPieChart(items []opts.PieData) *charts.Pie {
	if len(items) == 0 {
		return nil
	}
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "light",
			Width:  "1280px",
			Height: "720px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: PieChartTitle,
			Left:  "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Orient: "vertical",
			X:      "right",
			Y:      "bottom",
		}),
	)
	pie.AddSeries("", items)
	pie.SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{
			Show:      opts.Bool(true),
			Position:  "inside",
			Formatter: "{d}%",
		}),
	)
	return pie
}

func getBarSubtitle[E Entry](entries []E) string {
	if len(entries) == 0 {
		return ""
	}
	subtitle := ""
	for _, entry := range entries {
		if subtitle != "" {
			break
		}
		var (
			m = entry.DataSet()
			d = m[desiredStateLabel]
		)
		switch d {
		case 0:
			subtitle = "Desired state: Delete log groups"
		case 9999:
			subtitle = "Desired state: Delete retention policy"
		default:
			subtitle = fmt.Sprintf("Desired state: Change retention to %d days", d)
		}
	}
	return subtitle
}

func getBarItems[E Entry](entries []E) ([]string, []opts.BarData, []opts.BarData) {
	if len(entries) == 0 {
		return nil, nil, nil
	}
	var (
		rmOthersTotal int64
		rdOthersTotal int64
		lnames        = make([]string, 0, MaxBarChartItems)
		rmbytes       = make([]opts.BarData, 0, MaxBarChartItems)
		rdbytes       = make([]opts.BarData, 0, MaxBarChartItems)
	)
	for i, entry := range entries {
		var (
			m   = entry.DataSet()
			rmb = m[remainingBytesLabel]
			rdb = m[reducibleBytesLabel]
		)
		if m[storedBytesLabel] == 0 {
			continue
		}
		if i < MaxBarChartItems-1 {
			lnames = append(lnames, entry.Name())
			rmbytes = append(rmbytes, opts.BarData{Value: rmb})
			rdbytes = append(rdbytes, opts.BarData{Value: rdb})
		} else {
			rmOthersTotal += rmb
			rdOthersTotal += rdb
		}
	}
	if rmOthersTotal > 0 || rdOthersTotal > 0 {
		lnames = append(lnames, "others")
		rmbytes = append(rmbytes, opts.BarData{Value: rmOthersTotal})
		rdbytes = append(rdbytes, opts.BarData{Value: rdOthersTotal})
	}
	return lnames, rmbytes, rdbytes
}

func newBarChart(subtitle string, names []string, remainings, reducibles []opts.BarData) *charts.Bar {
	if len(remainings) == 0 || len(reducibles) == 0 {
		return nil
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "light",
			Width:  "1600px",
			Height: "900px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    BarChartTitle,
			Subtitle: subtitle,
			Left:     "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Orient: "vertical",
			X:      "right",
			Y:      "top",
		}),
		charts.WithGridOpts(opts.Grid{
			ContainLabel: opts.Bool(true),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Rotate: 45,
			},
			SplitLine: &opts.SplitLine{
				Show: opts.Bool(true),
			},
		}),
	)
	bar.SetXAxis(names)
	bar.AddSeries("Remaining bytes", remainings)
	bar.AddSeries("Reducible bytes", reducibles)
	bar.SetSeriesOptions(
		charts.WithBarChartOpts(opts.BarChart{
			Stack: "stack",
		}),
	)
	return bar
}
