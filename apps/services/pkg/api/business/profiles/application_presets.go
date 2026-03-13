package profiles

// PresetField defines a single field in an application form preset.
type PresetField struct {
	Label       string
	FieldType   string // short_text, long_text, url
	IsRequired  bool
	Placeholder string
}

// ApplicationPresetDef defines a preset template for application forms.
type ApplicationPresetDef struct {
	Key    string
	Label  string
	Fields []PresetField
}

// ApplicationPresets contains the built-in application form templates.
//
//nolint:gochecknoglobals
var ApplicationPresets = map[string]ApplicationPresetDef{
	"developer_community": {
		Key:   "developer_community",
		Label: "Developer Community",
		Fields: []PresetField{
			{
				Label:       "GitHub/GitLab username",
				FieldType:   "short_text",
				IsRequired:  true,
				Placeholder: "e.g. octocat",
			},
			{
				Label:       "Why do you want to join?",
				FieldType:   "long_text",
				IsRequired:  true,
				Placeholder: "",
			},
			{
				Label:       "How did you hear about us?",
				FieldType:   "short_text",
				IsRequired:  false,
				Placeholder: "",
			},
		},
	},
	"open_source_project": {
		Key:   "open_source_project",
		Label: "Open Source Project",
		Fields: []PresetField{
			{
				Label:       "GitHub username",
				FieldType:   "short_text",
				IsRequired:  true,
				Placeholder: "e.g. octocat",
			},
			{
				Label:       "Which area interests you?",
				FieldType:   "short_text",
				IsRequired:  true,
				Placeholder: "",
			},
			{
				Label:       "Relevant experience",
				FieldType:   "long_text",
				IsRequired:  false,
				Placeholder: "",
			},
		},
	},
	"professional_organization": {
		Key:   "professional_organization",
		Label: "Professional Organization",
		Fields: []PresetField{
			{
				Label:       "Current role & company",
				FieldType:   "short_text",
				IsRequired:  true,
				Placeholder: "",
			},
			{
				Label:       "LinkedIn profile",
				FieldType:   "url",
				IsRequired:  false,
				Placeholder: "https://linkedin.com/in/...",
			},
			{
				Label:       "Why do you want to join?",
				FieldType:   "long_text",
				IsRequired:  true,
				Placeholder: "",
			},
			{
				Label:       "Referral from existing member?",
				FieldType:   "short_text",
				IsRequired:  false,
				Placeholder: "",
			},
		},
	},
}

// ListApplicationPresets returns all available preset keys in a stable order.
func ListApplicationPresets() []ApplicationPresetDef {
	return []ApplicationPresetDef{
		ApplicationPresets["developer_community"],
		ApplicationPresets["open_source_project"],
		ApplicationPresets["professional_organization"],
	}
}
