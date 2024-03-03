package app

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/mirzakhany/chapar/internal/bus"
	"github.com/mirzakhany/chapar/internal/notify"
	"github.com/mirzakhany/chapar/ui"
	"github.com/mirzakhany/chapar/ui/fonts"
	"github.com/mirzakhany/chapar/ui/pages/console"
	"github.com/mirzakhany/chapar/ui/pages/envs"
	"github.com/mirzakhany/chapar/ui/pages/requests"
	"github.com/mirzakhany/chapar/ui/widgets"
)

type C = layout.Context
type D = layout.Dimensions

type UI struct {
	app   *ui.Application
	Theme *material.Theme

	sideBar *Sidebar
	header  *Header

	requestsPage *requests.Requests
	envsPage     *envs.Envs
	consolePage  *console.Console
	notification *widgets.Notification

	tipsOpen bool
}

// New creates a new UI using the Go Fonts.
func New(app *ui.Application) (*UI, error) {
	u := &UI{
		app: app,
	}
	fontCollection, err := fonts.Prepare()
	if err != nil {
		return nil, err
	}

	bus.Init()

	u.Theme = material.NewTheme()
	u.Theme.Shaper = text.NewShaper(text.WithCollection(fontCollection))
	// set foreground color
	u.Theme.Palette.Fg = color.NRGBA{R: 0xD7, G: 0xDA, B: 0xDE, A: 0xff}
	// set background color
	u.Theme.Palette.Bg = color.NRGBA{R: 0x20, G: 0x22, B: 0x24, A: 0xff}

	u.Theme.TextSize = unit.Sp(14)
	// console need to be initialized before other pages as its listening for logs
	u.consolePage = console.New()
	u.header = NewHeader()
	u.sideBar = NewSidebar(u.Theme)

	u.requestsPage, err = requests.New(u.Theme)
	if err != nil {
		return nil, err
	}

	u.envsPage, err = envs.New(u.Theme)
	if err != nil {
		return nil, err
	}

	u.notification = &widgets.Notification{}

	return u, nil
}

func (u *UI) Run(w *ui.Window) error {
	// ops are the operations from the UI
	var ops op.Ops

	for {
		switch e := w.NextEvent().(type) {
		// this is sent when the application should re-render.
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// render and handle UI.
			u.Layout(gtx, gtx.Constraints.Max.X)
			// render and handle the operations from the UI.
			e.Frame(gtx.Ops)
		// this is sent when the application is closed.
		case app.DestroyEvent:
			return e.Err
		}
	}
}

// Layout displays the main program layout.
func (u *UI) Layout(gtx layout.Context, windowWidth int) layout.Dimensions {
	// Create a full-screen rectangle for the background
	bgRect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, u.Theme.Palette.Bg, clip.Rect(bgRect).Op())

	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return u.header.Layout(gtx, u.Theme)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return u.sideBar.Layout(gtx, u.Theme)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							switch u.sideBar.SelectedIndex() {
							case 0:
								return u.requestsPage.Layout(gtx, u.Theme)
							case 1:
								return u.envsPage.Layout(gtx, u.Theme)
							case 4:
								return u.consolePage.Layout(gtx, u.Theme)
							}
							return layout.Dimensions{}
						}),
					)
				}),
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return notify.NotificationController.Layout(gtx, u.Theme, windowWidth)
		}),
	)
}