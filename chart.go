package llcm

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/pkg/browser"
)

var (
	// MaxPieChartItems is the maximum number of items in a pie chart.
	MaxPieChartItems = 11

	// MaxBarChartItems is the maximum number of items in a bar chart.
	MaxBarChartItems = 31
)

func renderChart(chart components.Charter) error {
	var (
		title = "llcm"
		fname = fmt.Sprintf("%s.html", title)
		i     = 1
	)
	for {
		if _, err := os.Stat(fname); err != nil {
			if os.IsNotExist(err) {
				break
			}
			return err
		}
		fname = fmt.Sprintf("%s%d.html", title, i)
		i++
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	page := components.NewPage()
	page.SetPageTitle(title)
	page.AddCharts(chart)
	if err := page.Render(io.MultiWriter(f)); err != nil {
		return err
	}
	browser.OpenFile(fname) //nolint:errcheck
	return nil
}

func getPieItems[E Entry](entries []E) (string, []opts.PieData) {
	var (
		othersTotal int64
		items       = make([]opts.PieData, 0, MaxPieChartItems)
		title       = "Stored bytes of log groups"
	)
	for i, entry := range entries {
		m := entry.DataSet()
		v := m["storedBytes"]
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
	return title, items
}

func newPieChart(title string, items []opts.PieData) *charts.Pie {
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
			Title: title,
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

func renderPieChart(pie *charts.Pie) error {
	if pie == nil {
		return nil
	}
	return renderChart(pie)
}

func getBarTitle[E Entry](entries []E) (string, string) {
	var (
		title    = "The simulation of reductions in log groups"
		subtitle = ""
	)
	for _, entry := range entries {
		if subtitle != "" {
			break
		}
		var (
			m = entry.DataSet()
			d = m["desiredState"]
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
	return title, subtitle
}

func getBarItems[E Entry](entries []E) ([]string, []opts.BarData, []opts.BarData) {
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
			rmb = m["remainingBytes"]
			rdb = m["reducibleBytes"]
		)
		if m["storedBytes"] == 0 {
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

func newBarChart(title, subtitle string, names []string, remainings, reducibles []opts.BarData) *charts.Bar {
	if len(remainings) == 0 && len(reducibles) == 0 {
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
			Title:    title,
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

func renderBarChart(bar *charts.Bar) error {
	if bar == nil {
		return nil
	}
	return renderChart(bar)
}
