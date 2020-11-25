package asr

type (
	// IEngine interface
	IEngine interface {
		prepare() (err error)
		recognize(filename string) (text string, err error)
	}
	// Engine struct
	Engine struct {
		IEngine
	}
)
