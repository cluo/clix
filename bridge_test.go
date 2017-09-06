package clix

import (
	"fmt"
	"github.com/urfave/cli"
	"testing"
)

var config struct {
	Listen       string `json:"listen" name:"listen,l" value:":0003" env:"LISTEN_ADDR" usage:"balalallaalalala"`
	Target       string `json:"target"`
	Key          string `json:"key"`
	Crypt        string `json:"crypt"`
	Mode         string `json:"mode"`
	MTU          int    `json:"mtu"`
	SndWnd       int    `json:"sndwnd"`
	RcvWnd       int    `json:"rcvwnd"`
	DataShard    int    `json:"datashard"`
	ParityShard  int    `json:"parityshard"`
	DSCP         int    `json:"dscp"`
	NoComp       bool   `json:"nocomp"`
	AckNodelay   bool   `json:"acknodelay"`
	NoDelay      int    `json:"nodelay"`
	Interval     int    `json:"interval"`
	Resend       int    `json:"resend"`
	NoCongestion int    `json:"nc"`
	SockBuf      int    `json:"sockbuf"`
	KeepAlive    int    `json:"keepalive"`
	Log          string `json:"log"`
	SnmpLog      string `json:"snmplog"`
	SnmpPeriod   int    `json:"snmpperiod"`
	Pprof        bool   `json:"pprof"`
}

func TestMakeFlags(t *testing.T) {
	flags := MakeFlags(&config)
	println(fmt.Sprintf("%+v", flags))

	c := &config
	fn := MakeParser(&c)
	fn(&cli.Context{})
	println(fmt.Sprintf("%+v", config))
}
