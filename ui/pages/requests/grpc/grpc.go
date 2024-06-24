package grpc

import (
	"gioui.org/layout"
	"gioui.org/unit"
	giox "gioui.org/x/component"

	"github.com/chapar-rest/chapar/internal/domain"
	"github.com/chapar-rest/chapar/ui/chapartheme"
	"github.com/chapar-rest/chapar/ui/converter"
	"github.com/chapar-rest/chapar/ui/pages/requests/component"
	"github.com/chapar-rest/chapar/ui/widgets"
)

type Grpc struct {
	Prompt *widgets.Prompt

	Req *domain.Request

	Breadcrumb *component.Breadcrumb
	AddressBar *AddressBar

	Request  *Request
	Response *Response

	split widgets.SplitView

	onSave        func(id string)
	onDataChanged func(id string, data any)
	onSubmit      func(id string)
}

func (r *Grpc) SetOnTitleChanged(f func(title string)) {
	r.Breadcrumb.SetOnTitleChanged(f)
}

func (r *Grpc) SetDataChanged(changed bool) {
	r.Breadcrumb.SetDataChanged(changed)
}

func New(req *domain.Request, theme *chapartheme.Theme) *Grpc {
	r := &Grpc{
		Req:        req,
		Prompt:     widgets.NewPrompt("", "", ""),
		Breadcrumb: component.NewBreadcrumb(req.MetaData.ID, req.CollectionName, "gRPC", req.MetaData.Name),
		split: widgets.SplitView{
			Resize: giox.Resize{
				Ratio: 0.5,
			},
			BarWidth: unit.Dp(2),
		},
		AddressBar: NewAddressBar(theme, req.Spec.GRPC.ServerInfo.Address, req.Spec.GRPC.LasSelectedMethod, req.Spec.GRPC.Services),
		Request:    NewRequest(req, theme),
		Response:   NewResponse(theme),
	}

	r.setupHooks()

	return r
}

func (r *Grpc) setupHooks() {
	r.AddressBar.SetOnServerAddressChanged(func(address string) {
		r.Req.Spec.GRPC.ServerInfo.Address = address
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.AddressBar.SetOnMethodChanged(func(method string) {
		r.Req.Spec.GRPC.LasSelectedMethod = method
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.AddressBar.SetOnSubmit(func() {
		r.onSubmit(r.Req.MetaData.ID)
	})

	r.Breadcrumb.SetOnSave(func(id string) {
		r.onSave(id)
	})

	r.Request.Body.SetOnChanged(func(data string) {
		r.Req.Spec.GRPC.Body = data
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.Request.Auth.SetOnChange(func(auth domain.Auth) {
		r.Req.Spec.GRPC.Auth = auth
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.Request.Metadata.SetOnChanged(func(items []*widgets.KeyValueItem) {
		data := converter.KeyValueFromWidgetItems(items)
		r.Req.Spec.GRPC.Metadata = data
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.Request.Settings.SetOnChange(func(values map[string]any) {
		r.Req.Spec.GRPC.Settings = convertSettingsToItems(values)
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.Request.ServerInfo.FileSelector.SetOnChanged(func(filePath string) {
		if r.Req.Spec.GRPC.ServerInfo.ProtoFiles == nil {
			r.Req.Spec.GRPC.ServerInfo.ProtoFiles = make([]string, 0)
		}

		r.Req.Spec.GRPC.ServerInfo.ProtoFiles = append(r.Req.Spec.GRPC.ServerInfo.ProtoFiles, filePath)
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})

	r.Request.ServerInfo.SetOnChanged(func() {
		r.Req.Spec.GRPC.ServerInfo.ServerReflection = r.Request.ServerInfo.definitionFrom.Value == "reflection"
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	})
}

func convertSettingsToItems(values map[string]any) domain.Settings {
	out := domain.Settings{
		UseSSL:              false,
		TimeoutMilliseconds: 1000,
		NameOverride:        "",
	}

	if v, ok := values["useSSL"]; ok {
		out.UseSSL = v.(bool)
	}

	if v, ok := values["timeoutMilliseconds"]; ok {
		out.TimeoutMilliseconds = v.(int)
	}

	if v, ok := values["nameOverride"]; ok {
		out.NameOverride = v.(string)
	}

	return out
}

func (r *Grpc) SetOnProtoFileSelect(f func(id string)) {
	r.Request.ServerInfo.FileSelector.SetOnSelectFile(func() {
		f(r.Req.MetaData.ID)
	})
}

func (r *Grpc) SetProtoBodyFilePath(filePath string) {
	r.Request.ServerInfo.FileSelector.SetFileName(filePath)
	if r.onDataChanged != nil {
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	}
}

func (r *Grpc) SetOnDataChanged(f func(id string, data any)) {
	r.onDataChanged = f
}

func (r *Grpc) SetOnSubmit(f func(id string)) {
	r.onSubmit = f
}

func (r *Grpc) SetOnSave(f func(id string)) {
	r.onSave = f
	r.Breadcrumb.SetOnSave(f)
}

func (r *Grpc) HidePrompt() {
	r.Prompt.Hide()
}

func (r *Grpc) ShowPrompt(title, content, modalType string, onSubmit func(selectedOption string, remember bool), options ...widgets.Option) {
	r.Prompt.Type = modalType
	r.Prompt.Title = title
	r.Prompt.Content = content
	r.Prompt.SetOptions(options...)
	r.Prompt.WithoutRememberBool()
	r.Prompt.SetOnSubmit(onSubmit)
	r.Prompt.Show()
}

func (r *Grpc) SetOnReload(f func(id string)) {
	if r.Request.ServerInfo == nil {
		return
	}

	r.Request.ServerInfo.SetOnReload(func() {
		f(r.Req.MetaData.ID)
	})
}

func (r *Grpc) SetServices(services []domain.GRPCService) {
	r.Req.Spec.GRPC.Services = services
	if r.onDataChanged != nil {
		r.onDataChanged(r.Req.MetaData.ID, r.Req)
	}

	r.AddressBar.SetServices(services)
}

func (r *Grpc) ShowMethodsLoading() {
	r.Request.ServerInfo.IsLoading = true
}

func (r *Grpc) HideMethodsLoading() {
	r.Request.ServerInfo.IsLoading = false
}

func (r *Grpc) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return r.Prompt.Layout(gtx, theme)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15), Top: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return r.Breadcrumb.Layout(gtx, theme)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return r.AddressBar.Layout(gtx, theme)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return r.split.Layout(gtx, theme,
					func(gtx layout.Context) layout.Dimensions {
						return r.Request.Layout(gtx, theme)
					},
					func(gtx layout.Context) layout.Dimensions {
						return r.Response.Layout(gtx, theme)
					},
				)
			}),
		)
	})
}
