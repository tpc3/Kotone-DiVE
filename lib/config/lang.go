package config

type Strings struct {
	Lang  string
	Error errorstr
}

type errorstr struct {
	Title   string
	Unknown string
	Join    joinerrorstr
	Leave   leaveerrorstr
	Guild   guilderrorstr
	Config  configerrorstr
}

type guilderrorstr struct {
	Prefix  string
	MaxChar string
	Voice   string
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
	SubCmd string
	Value  string
}

var (
	Lang map[string]Strings
)

func loadLang() {
	Lang = map[string]Strings{}
	Lang["japanese"] = Strings{
		Lang: "japanese",
		Error: errorstr{
			Title:   "エラー",
			Unknown: "不明なエラーが発生しました。\nこの問題は管理者に報告されます。",
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
				Voice:   "そのようなボイスは存在しません。",
				Policy:  "ポリシーはallowかdenyである必要があります。",
			},
			Config: configerrorstr{
				SubCmd: "サブコマンドが不正です。",
				Value:  "不正な設定値です。",
			},
		},
	}
	Lang["english"] = Strings{
		Lang: "english",
		Error: errorstr{
			Title:   "Error",
			Unknown: "Unknown Error!\nThis will be reported.",
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
				Voice:   "No such voice.",
				Policy:  "Policy should be allow or deny.",
			},
			Config: configerrorstr{
				SubCmd: "Invalid subcommand.",
				Value:  "Invalid value.",
			},
		},
	}
}
