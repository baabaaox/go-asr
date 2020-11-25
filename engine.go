package asr

type (
	// IEngine interface
	IEngine interface {
		Prepare() (err error)
		Recognize(filename string) (text string, err error)
	}
	// Engine struct
	Engine struct {
		IEngine
	}
)
