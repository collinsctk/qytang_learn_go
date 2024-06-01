package plot_cpu

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// XY is a struct representing a single point in the plot.
type XY struct {
	X float64
	Y float64
}

// XYs is a slice of XY points.
type XYs []XY

// Len returns the length of the XYs slice.
func (xy XYs) Len() int {
	return len(xy)
}

// XY returns the x and y values at the specified index.
func (xy XYs) XY(i int) (float64, float64) {
	return xy[i].X, xy[i].Y
}

// PlotData 绘制图表并保存到文件
func PlotData(deviceData map[string]XYs, title, xlabel, ylabel, imagePath string) error {
	// 创建绘图
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel
	p.Y.Max = 100
	p.Y.Min = 0

	// 添加数据点
	i := 0
	for deviceIP, data := range deviceData {
		line, err := plotter.NewLine(data)
		if err != nil {
			return err
		}
		line.Color = plotutil.Color(i)
		p.Add(line)
		p.Legend.Add(deviceIP, line)
		i++
	}

	// 将图例放置在右上角
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.XOffs = -vg.Length(5) // 让图例往左移动一些

	// 保存绘图到文件
	if err := p.Save(10*vg.Inch, 5*vg.Inch, imagePath); err != nil {
		return err
	}

	return nil
}
