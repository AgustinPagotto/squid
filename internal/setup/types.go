package setup

type Language int

const (
	LangGo Language = iota + 1
	LangNodeJs
	LangPython
	LangReact
	LangReactNative
	LangNextJs
)

func (l Language) filename() string {
	return map[Language]string{
		LangGo:          "go.json",
		LangNodeJs:      "nodejs.json",
		LangPython:      "python.json",
		LangReact:       "react.json",
		LangReactNative: "react-native.json",
		LangNextJs:      "nextjs.json",
	}[l]
}

type Shell int

const (
	ShellBash Shell = iota + 1
	ShellZsh
	ShellFish
)

func (s Shell) configFile() string {
	return map[Shell]string{
		ShellBash: ".bashrc",
		ShellZsh:  ".zshrc",
		ShellFish: "config.fish",
	}[s]
}

func (s Shell) evalString() string {
	return map[Shell]string{
		ShellBash: bashEvalConfig,
		ShellZsh:  zshEvalConfig,
		ShellFish: fishEvalConfig,
	}[s]
}

type Action int

const (
	ActionSelect Action = iota + 1
	ActionBack
	ActionExit
)
