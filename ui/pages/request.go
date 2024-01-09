package pages

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mirzakhany/chapar/ui/widgets"
)

type Request struct {
	theme *material.Theme

	actions *Actions

	textEditor widget.Editor
	searchBox  *widgets.TextField

	split widgets.SplitView
}

func NewRequest(theme *material.Theme) *Request {
	search := widgets.NewTextField(theme, "Search...")
	search.SetIcon(widgets.SearchIcon, widgets.IconPositionEnd)
	return &Request{
		theme:     theme,
		actions:   NewActions(theme),
		searchBox: search,
		split: widgets.SplitView{
			Ratio:         -0.64,
			MinLeftSize:   420,
			MaxLeftSize:   800,
			BarWidth:      unit.Dp(2),
			BarColor:      color.NRGBA{R: 0x2b, G: 0x2d, B: 0x31, A: 0xff},
			BarColorHover: theme.Palette.ContrastBg,
		},
	}
}

func (r *Request) list(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Top:    10,
		Bottom: 0,
		Left:   10,
		Right:  10,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return r.actions.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return r.searchBox.Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.List(r.theme, &widget.List{
					List: layout.List{
						Axis: layout.Vertical,
					},
				},
				).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Body1(r.theme, "Request").Layout(gtx)
					})
				})
			}),
		)
	})
}

func (r *Request) Layout(gtx layout.Context) layout.Dimensions {

	return r.split.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return r.list(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H3(r.theme, "Request").Layout(gtx)
			})
		},
	)

	//return layout.Flex{Axis: layout.Horizontal, WeightSum: 12}.Layout(gtx,
	//	layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
	//		return r.list(gtx)
	//	}),
	//	widgets.VerticalFullLine(),
	//	layout.Flexed(9, func(gtx layout.Context) layout.Dimensions {
	//		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//			return material.H3(r.theme, "Request").Layout(gtx)
	//		})
	//	}),
	//)
}
