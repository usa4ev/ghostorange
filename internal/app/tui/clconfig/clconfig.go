package clconfig

import (
	"flag"
	"os"
)

type Config struct{
	srvAddr string
	logPath string
}

type (
	option func(*Config)
)

func WithAddress(addr string)option{
	return func(c *Config){
		c.srvAddr = addr
	}
}

func New(opts... option)*Config{

	c:= Config{}

	for _,opt := range opts{
		opt(&c)
	}

	fs := flag.NewFlagSet("myFS", flag.ContinueOnError)
	if !fs.Parsed() {
		fs.StringVar(&c.srvAddr, "a", c.srvAddr, "the service address")
		fs.StringVar(&c.logPath, "l", c.logPath, "path to write log")

		fs.Parse(os.Args[1:])
	}

	return &c

}

func (c Config) SrvAddr()string{
	return c.srvAddr
}

func (c Config) LogPath()string{
	return c.logPath
}