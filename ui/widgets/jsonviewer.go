package widgets

import (
	"strings"

	"gioui.org/text"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mirzakhany/chapar/ui/chapartheme"
)

const chunkSize = 100

type JsonViewer struct {
	//	data string

	lines   []string
	chunks  []string
	editors []*widget.Editor

	//selectables []*widget.Selectable

	list *widget.List
}

func NewJsonViewer() *JsonViewer {
	return &JsonViewer{
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
}

func (j *JsonViewer) SetData(data string) {
	//j.data = data
	j.lines = strings.Split(data, "\n")

	j.chunks = make([]string, 0)
	for i := 0; i < len(j.lines); i += chunkSize {
		end := i + chunkSize
		if end > len(j.lines) {
			end = len(j.lines)
		}
		j.chunks = append(j.chunks, strings.Join(j.lines[i:end], "\n"))
	}

	j.editors = make([]*widget.Editor, len(j.chunks))
	for i := range j.editors {
		j.editors[i] = new(widget.Editor)
		j.editors[i].Submit = false
		j.editors[i].SingleLine = false
		j.editors[i].SetText(j.chunks[i])
		j.editors[i].WrapPolicy = text.WrapGraphemes
		j.editors[i].ReadOnly = true
	}

	//j.selectables = make([]*widget.Selectable, len(j.lines))
	//for i := range j.selectables {
	//	j.selectables[i] = &widget.Selectable{}
	//}
}

func (j *JsonViewer) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
	border := widget.Border{
		Color:        theme.BorderColor,
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}

	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(3).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.List(theme.Material(), j.list).Layout(gtx, len(j.chunks), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					//	return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					//		l := material.Label(theme.Material(), theme.TextSize, fmt.Sprintf("%d", i+1))
					//		l.Font.Weight = font.Medium
					//		l.Color = theme.TextColor
					//		l.SelectionColor = theme.TextSelectionColor
					//		l.Alignment = text.End
					//		return l.Layout(gtx)
					//	})
					//}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							//l := material.Label(theme.Material(), theme.TextSize, j.lines[i])
							////l.State = j.selectables[i]
							//l.SelectionColor = theme.TextSelectionColor
							////l.TextSize = unit.Sp(14)
							//return l.Layout(gtx)
							ed := material.EditorStyle{}
							ed = material.Editor(theme.Material(), j.editors[i], "")
							ed.SelectionColor = theme.TextSelectionColor
							return ed.Layout(gtx)
						})
					}),
				)
			})
		})
	})
}
