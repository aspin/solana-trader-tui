package store

type App struct {
	Err      error
	UI       UI
	Settings Settings
}

type UI struct {
	WindowWidth  int
	WindowHeight int
}

type Settings struct {
	AuthHeader        string
	PrivateKey        string
	PublicKey         string
	OpenOrdersAddress string
}

func (a App) NeedsInit() bool {
	return a.Settings.AuthHeader == ""
}
