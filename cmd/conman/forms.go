package main

import (
	"fyne.io/fyne/widget"
)

// AuthForm ...
type AuthForm struct {
	*widget.Box

	Strategy *Strategy
	Entries  map[string]*widget.Entry
}

// GetCreds ...
func (f *AuthForm) GetCreds() map[string]string {
	res := map[string]string{}
	for k, f := range f.Entries {
		res[k] = f.Text
	}
	return res
}

// NewAuthForm ...
func NewAuthForm() *AuthForm {
	ss := []string{}
	for _, s := range Strategies {
		ss = append(ss, s.Text)
	}

	siteURL := widget.NewEntry()
	strats := widget.NewSelectEntry(ss)

	return genForm(siteURL, strats, "")
}

func genForm(siteURL *widget.Entry, strats *widget.SelectEntry, strat string) *AuthForm {
	form := &AuthForm{Box: widget.NewVBox(), Entries: map[string]*widget.Entry{}}

	// siteURL.PlaceHolder = "Site URL"
	form.Append(NewFormField("Site URL", siteURL))
	// strats.PlaceHolder = "Strategy"
	form.Append(NewFormField("Strategy", strats))
	setStrategy(form, "")

	strats.OnChanged = func(s string) {
		setStrategy(form, s)
	}

	return form
}

func setStrategy(form *AuthForm, strat string) {
	form.ExtendBaseWidget(form)
	form.Children = form.Children[0:2]
	for _, s := range Strategies {
		if s.Text == strat {
			form.Strategy = s
			for _, f := range s.Fields {
				e := widget.NewEntry()
				if f[1] == "password" {
					e = widget.NewPasswordEntry()
				}
				form.Entries[f[1]] = e
				form.Append(NewFormField(f[0], e))

				// e.PlaceHolder = f[0]
				// form.Append(NewFormField("", e))
			}
		}
	}
	form.Append(widget.NewButton("Check auth", func() {}))
	// form.ExtendBaseWidget(form)
	form.Refresh()
}
