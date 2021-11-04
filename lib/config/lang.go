package config

type Strings struct {
	Lang  string
	Help  string
	Error errorstr
}

type errorstr struct {
	Title     string
	Unknown   string
	SubCmd    string
	Voice     string
	Joinfirst string
	Join      joinerrorstr
	Leave     leaveerrorstr
	Guild     guilderrorstr
	Config    configerrorstr
	Replace   replaceerrorstr
	Policy    policyerrorstr
	Skip      skiperrorstr
}

type guilderrorstr struct {
	Prefix  string
	MaxChar string
	Policy  string
}

type joinerrorstr struct {
	Already string
	Failed  string
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
	Empty  string
}

type policyerrorstr struct {
	Exists    string
	NotExists string
	Empty     string
}

type skiperrorstr struct {
	NotPlaying string
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
			Title:     "エラー",
			Unknown:   "不明なエラーが発生しました。\nこの問題は管理者に報告されます。",
			SubCmd:    "サブコマンドが不正です。",
			Voice:     "そのようなボイスは存在しません。",
			Joinfirst: "まずはVCに参加してください。",
			Join: joinerrorstr{
				Already: "既にVCに接続済です。",
				Failed:  "接続に失敗しました。権限設定をご確認ください。",
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
				Regex:  "正規表現が間違っています。",
				Del:    "削除対象が見つかりません。",
				Syntax: "コマンド文法エラーです。",
				Empty:  "空の置換リストを表示することはできません。",
			},
			Policy: policyerrorstr{
				Exists:    "ポリシーが既に存在しています。",
				NotExists: "ポリシーが存在しません。",
				Empty:     "空のポリシーを表示することはできません。",
			},
			Skip: skiperrorstr{
				NotPlaying: "読み上げ中ではありません。",
			},
		},
	}
	Lang["english"] = Strings{
		Lang: "english",

		Help: "Usage is available on the Wiki.",
		Error: errorstr{
			Title:     "Error",
			Unknown:   "Unknown Error!\nThis will be reported.",
			SubCmd:    "Invalid subcommand.",
			Voice:     "No such voice.",
			Joinfirst: "You must join VC first",
			Join: joinerrorstr{
				Already: "I've already joined.",
				Failed:  "Connection failed. please check your server permissions.",
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
				Del:    "So such key in the database.",
				Empty:  "Cannot print empty replace list.",
			},
			Policy: policyerrorstr{
				Exists:    "Policy already exists.",
				NotExists: "Policy doesn't exists.",
				Empty:     "Cannot print empty policy list.",
			},
			Skip: skiperrorstr{
				NotPlaying: "I'm not playing anything.",
			},
		},
	}
}
