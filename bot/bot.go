package bot

type Bot interface {
	RunOnce() msg
}

type botDoneMsg Bot

type botErrorMsg error
