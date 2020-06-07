package main

// Strategy ...
type Strategy struct {
	Text    string
	Code    string
	Targets []string
	Fields  [][]string
}

// Strategies ...
var Strategies = []*Strategy{
	{
		Text:    "Add-in",
		Code:    "addin",
		Targets: []string{"spo"},
		Fields: [][]string{
			{"Client ID", "clientId"},
			{"Client Secret", "clientSecret"},
		},
	},
	{
		Text:    "NTLM",
		Code:    "ntlm",
		Targets: []string{"on-prem"},
		Fields: [][]string{
			{"User Name", "username"},
			{"Password", "password"},
		},
	},
	{
		Text:    "SAML",
		Code:    "saml",
		Targets: []string{"spo"},
		Fields: [][]string{
			{"User Name", "username"},
			{"Password", "password"},
		},
	},
}
