package config

type Strings struct {
	Lang  string
	Help  string
	Error errorstr
}

type errorstr struct {
	Title   string
	Unknown string
	SubCmd  string
	Voice   string
	Join    joinerrorstr
	Leave   leaveerrorstr
	Guild   guilderrorstr
	Config  configerrorstr
	Replace replaceerrorstr
	Policy  policyerrorstr
}

type guilderrorstr struct {
	Prefix  string
	MaxChar string
	Policy  string
}

type joinerrorstr struct {
	Already   string
	Joinfirst string
	Failed    string
}

type leaveerrorstr struct {
	None string
}

type configerrorstr struct {
	Value string
}

type replaceerrorstr struct {
	Syntax string
	Regex  string
	Del    string
}

type policyerrorstr struct {
	Exists    string
	NotExists string
}

var (
	Lang map[string]Strings
)

func loadLang() {
	Lang = map[string]Strings{}
	Lang["japanese"] = Strings{
		Lang: "japanese",
		Help: "Botの使い方に関しては、Wikiをご覧ください。",
		Error: errorstr{
			Title:   "エラー",
			Unknown: "不明なエラーが発生しました。\nこの問題は管理者に報告されます。",
			SubCmd:  "サブコマンドが不正です。",
			Voice:   "そのようなボイスは存在しません。",
			Join: joinerrorstr{
				Already:   "既にVCに接続済です。",
				Joinfirst: "まずはVCに参加してください。",
				Failed:    "接続に失敗しました。権限設定をご確認ください。",
			},
			Leave: leaveerrorstr{
				None: "VCに参加していません。",
			},
			Guild: guilderrorstr{
				Prefix:  "プレフィクスは1文字である必要があります。",
				MaxChar: "最大文字数は0文字以上2000文字以下である必要があります。",
				Policy:  "ポリシーはallowかdenyである必要があります。",
			},
			Config: configerrorstr{
				Value: "不正な設定値です。",
			},
			Replace: replaceerrorstr{
				Regex: "正規表現が間違っています。",
				Del:   "削除対象が見つかりません。",
			},
			Policy: policyerrorstr{
				Exists:    "ポリシーが既に存在しています。",
				NotExists: "ポリシーが存在しません。",
			},
		},
	}
	Lang["english"] = Strings{
		Lang: "english",

		Help: "Usage is available on the Wiki.",
		Error: errorstr{
			Title:   "Error",
			Unknown: "Unknown Error!\nThis will be reported.",
			SubCmd:  "Invalid subcommand.",
			Voice:   "No such voice.",
			Join: joinerrorstr{
				Already:   "I've already joined.",
				Joinfirst: "You must join VC first",
				Failed:    "Connection failed. please check your server permissions.",
			},
			Leave: leaveerrorstr{
				None: "No VC to leave.",
			},
			Guild: guilderrorstr{
				Prefix:  "Prefix must be a single char.",
				MaxChar: "Max char must be between 0 and 2000.",
				Policy:  "Policy should be allow or deny.",
			},
			Config: configerrorstr{
				Value: "Invalid value.",
			},
			Replace: replaceerrorstr{
				Syntax: "Command syntax error.",
				Regex:  "Invalid regex.",
				Del:    "So such key in the database",
			},
			Policy: policyerrorstr{
				Exists:    "Policy already exists.",
				NotExists: "Policy doesn't exists.",
			},
		},
	}
}
