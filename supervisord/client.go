package supervisord

import (
	"github.com/shivmohith/go-supervisord"
)

type Client interface {
	GetAllProcessInfo() ([]supervisord.ProcessInfo, error)
	TailProcessStdoutLog(name string, offset int, length int) (*supervisord.LogSegment, error)
	TailProcessStderrLog(name string, offset int, length int) (*supervisord.LogSegment, error)
}

func NewClient() (Client, error) {
	c, err := supervisord.NewClient("http://127.0.0.1:9001/RPC2")
	if err != nil {
		return nil, err
	}

	return c, nil
}
