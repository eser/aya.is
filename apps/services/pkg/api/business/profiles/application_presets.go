package profiles

// PresetField defines a single field in an application form preset.
type PresetField struct {
	Label       string `json:"label"`
	FieldType   string `json:"field_type"` // short_text, long_text, url
	Placeholder string `json:"placeholder"`
	IsRequired  bool   `json:"is_required"`
}

// ApplicationPresetDef defines a preset template for application forms.
type ApplicationPresetDef struct {
	Key    string        `json:"key"`
	Label  string        `json:"label"`
	Fields []PresetField `json:"fields"`
}

// presetTranslations maps locale → preset key → translated label + field labels.
type presetTx struct {
	Label  string
	Fields []string // parallel to base fields array
}

//nolint:gochecknoglobals
var presetTranslations = map[string]map[string]presetTx{
	"en": {
		"developer_community": {
			Label:  "Developer Community",
			Fields: []string{"Why do you want to join?", "How did you hear about us?"},
		},
		"open_source_project": {
			Label:  "Open Source Project",
			Fields: []string{"Which area interests you?", "Relevant experience"},
		},
		"professional_organization": {
			Label:  "Professional Organization",
			Fields: []string{"Current role & company", "Why do you want to join?"},
		},
	},
	"tr": {
		"developer_community": {
			Label:  "Geliştirici Topluluğu",
			Fields: []string{"Neden katılmak istiyorsunuz?", "Bizi nereden duydunuz?"},
		},
		"open_source_project": {
			Label:  "Açık Kaynak Projesi",
			Fields: []string{"Hangi alan ilginizi çekiyor?", "İlgili deneyimleriniz"},
		},
		"professional_organization": {
			Label: "Profesyonel Organizasyon",
			Fields: []string{
				"Mevcut pozisyon ve şirket",
				"Neden katılmak istiyorsunuz?",
				"Mevcut bir üyeden referans?",
			},
		},
	},
	"de": {
		"developer_community": {
			Label:  "Entwickler-Community",
			Fields: []string{"Warum möchten Sie beitreten?", "Wie haben Sie von uns erfahren?"},
		},
		"open_source_project": {
			Label:  "Open-Source-Projekt",
			Fields: []string{"Welcher Bereich interessiert Sie?", "Relevante Erfahrung"},
		},
		"professional_organization": {
			Label:  "Professionelle Organisation",
			Fields: []string{"Aktuelle Position & Unternehmen", "Warum möchten Sie beitreten?"},
		},
	},
	"es": {
		"developer_community": {
			Label:  "Comunidad de desarrolladores",
			Fields: []string{"¿Por qué quieres unirte?", "¿Cómo nos conociste?"},
		},
		"open_source_project": {
			Label:  "Proyecto de código abierto",
			Fields: []string{"¿Qué área te interesa?", "Experiencia relevante"},
		},
		"professional_organization": {
			Label:  "Organización profesional",
			Fields: []string{"Puesto actual y empresa", "¿Por qué quieres unirte?"},
		},
	},
	"fr": {
		"developer_community": {
			Label: "Communauté de développeurs",
			Fields: []string{
				"Pourquoi souhaitez-vous rejoindre ?",
				"Comment avez-vous entendu parler de nous ?",
			},
		},
		"open_source_project": {
			Label:  "Projet open source",
			Fields: []string{"Quel domaine vous intéresse ?", "Expérience pertinente"},
		},
		"professional_organization": {
			Label:  "Organisation professionnelle",
			Fields: []string{"Poste actuel et entreprise", "Pourquoi souhaitez-vous rejoindre ?"},
		},
	},
	"it": {
		"developer_community": {
			Label:  "Comunità di sviluppatori",
			Fields: []string{"Perché vuoi unirti?", "Come hai saputo di noi?"},
		},
		"open_source_project": {
			Label:  "Progetto open source",
			Fields: []string{"Quale area ti interessa?", "Esperienza rilevante"},
		},
		"professional_organization": {
			Label:  "Organizzazione professionale",
			Fields: []string{"Ruolo attuale e azienda", "Perché vuoi unirti?"},
		},
	},
	"ja": {
		"developer_community": {
			Label:  "開発者コミュニティ",                         //nolint:gosmopolitan // Japanese
			Fields: []string{"参加したい理由は？", "どこで知りましたか？"}, //nolint:gosmopolitan // Japanese
		},
		"open_source_project": {
			Label:  "オープンソースプロジェクト",
			Fields: []string{"興味のある分野は？", "関連する経験"}, //nolint:gosmopolitan // Japanese
		},
		"professional_organization": {
			Label:  "専門組織",                            //nolint:gosmopolitan // Japanese
			Fields: []string{"現在の役職と会社", "参加したい理由は？"}, //nolint:gosmopolitan // Japanese
		},
	},
	"ko": {
		"developer_community": {
			Label:  "개발자 커뮤니티",
			Fields: []string{"가입하려는 이유는?", "어떻게 알게 되셨나요?"},
		},
		"open_source_project": {
			Label:  "오픈소스 프로젝트",
			Fields: []string{"관심 있는 분야는?", "관련 경험"},
		},
		"professional_organization": {
			Label:  "전문 조직",
			Fields: []string{"현재 직책 및 회사", "가입하려는 이유는?"},
		},
	},
	"nl": {
		"developer_community": {
			Label:  "Ontwikkelaarsgemeenschap",
			Fields: []string{"Waarom wilt u lid worden?", "Hoe heeft u van ons gehoord?"},
		},
		"open_source_project": {
			Label:  "Open source-project",
			Fields: []string{"Welk gebied interesseert u?", "Relevante ervaring"},
		},
		"professional_organization": {
			Label:  "Professionele organisatie",
			Fields: []string{"Huidige functie en bedrijf", "Waarom wilt u lid worden?"},
		},
	},
	"pt-PT": {
		"developer_community": {
			Label:  "Comunidade de desenvolvedores",
			Fields: []string{"Por que deseja participar?", "Como soube de nós?"},
		},
		"open_source_project": {
			Label:  "Projeto de código aberto",
			Fields: []string{"Qual área lhe interessa?", "Experiência relevante"},
		},
		"professional_organization": {
			Label:  "Organização professional",
			Fields: []string{"Cargo atual e empresa", "Por que deseja participar?"},
		},
	},
	"ru": {
		"developer_community": {
			Label:  "Сообщество разработчиков",
			Fields: []string{"Почему вы хотите присоединиться?", "Как вы о нас узнали?"},
		},
		"open_source_project": {
			Label:  "Проект с открытым кодом",
			Fields: []string{"Какая область вас интересует?", "Соответствующий опыт"},
		},
		"professional_organization": {
			Label:  "Профессиональная организация",
			Fields: []string{"Текущая должность и компания", "Почему вы хотите присоединиться?"},
		},
	},
	"ar": {
		"developer_community": {
			Label:  "مجتمع المطورين",
			Fields: []string{"لماذا تريد الانضمام؟", "كيف سمعت عنا؟"},
		},
		"open_source_project": {
			Label:  "مشروع مفتوح المصدر",
			Fields: []string{"ما المجال الذي يهمك؟", "الخبرة ذات الصلة"},
		},
		"professional_organization": {
			Label:  "منظمة مهنية",
			Fields: []string{"الدور الحالي والشركة", "لماذا تريد الانضمام؟"},
		},
	},
	"zh-CN": {
		"developer_community": {
			Label:  "开发者社区",                           //nolint:gosmopolitan // Chinese
			Fields: []string{"为什么想加入？", "你是怎么知道我们的？"}, //nolint:gosmopolitan // Chinese
		},
		"open_source_project": {
			Label:  "开源项目",                         //nolint:gosmopolitan // Chinese
			Fields: []string{"你对哪个领域感兴趣？", "相关经验"}, //nolint:gosmopolitan // Chinese
		},
		"professional_organization": {
			Label:  "专业组织",                         //nolint:gosmopolitan // Chinese
			Fields: []string{"当前职位和公司", "为什么想加入？"}, //nolint:gosmopolitan // Chinese
		},
	},
}

// basePresets defines the structure (field types, required flags) — labels come from translations.
//
//nolint:gochecknoglobals
var basePresets = []struct {
	Key    string
	Fields []struct {
		FieldType  string
		IsRequired bool
	}
}{
	{
		Key: "developer_community",
		Fields: []struct {
			FieldType  string
			IsRequired bool
		}{
			{FieldType: "long_text", IsRequired: true},
			{FieldType: "short_text", IsRequired: false},
		},
	},
	{
		Key: "open_source_project",
		Fields: []struct {
			FieldType  string
			IsRequired bool
		}{
			{FieldType: "short_text", IsRequired: true},
			{FieldType: "long_text", IsRequired: false},
		},
	},
	{
		Key: "professional_organization",
		Fields: []struct {
			FieldType  string
			IsRequired bool
		}{
			{FieldType: "short_text", IsRequired: true},
			{FieldType: "long_text", IsRequired: true},
		},
	},
}

// ApplicationPresets is a lookup map used by UpsertApplicationForm for auto-populating fields.
// Uses English labels as default when the locale-aware ListApplicationPresets is not used.
//
//nolint:gochecknoglobals
var ApplicationPresets map[string]ApplicationPresetDef

//nolint:gochecknoinits
func init() {
	ApplicationPresets = make(map[string]ApplicationPresetDef)

	for _, bp := range basePresets {
		ApplicationPresets[bp.Key] = buildPreset(bp.Key, "en")
	}
}

func buildPreset(key string, locale string) ApplicationPresetDef {
	translation, found := presetTranslations[locale]
	if !found {
		translation = presetTranslations["en"]
	}

	ptx, found := translation[key]
	if !found {
		ptx = presetTranslations["en"][key]
	}

	var base *struct {
		Key    string
		Fields []struct {
			FieldType  string
			IsRequired bool
		}
	}

	for i := range basePresets {
		if basePresets[i].Key == key {
			base = &basePresets[i]

			break
		}
	}

	if base == nil {
		return ApplicationPresetDef{Key: key, Label: key, Fields: nil}
	}

	fields := make([]PresetField, 0, len(base.Fields))

	for i, baseField := range base.Fields {
		label := ""
		if i < len(ptx.Fields) {
			label = ptx.Fields[i]
		}

		fields = append(fields, PresetField{
			Label:       label,
			FieldType:   baseField.FieldType,
			IsRequired:  baseField.IsRequired,
			Placeholder: "",
		})
	}

	return ApplicationPresetDef{
		Key:    key,
		Label:  ptx.Label,
		Fields: fields,
	}
}

// ListApplicationPresets returns all available presets translated to the given locale.
func ListApplicationPresets(locale string) []ApplicationPresetDef {
	result := make([]ApplicationPresetDef, 0, len(basePresets))

	for _, bp := range basePresets {
		result = append(result, buildPreset(bp.Key, locale))
	}

	return result
}
