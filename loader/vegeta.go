package loader

import (
	"bytes"
	"github.com/globocom/vegeta/lib"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	"gopkg.in/square/go-jose.v1/json"
)

type attackOpts struct {
	target      string
	http2       bool
	duration    time.Duration
	timeout     time.Duration
	rate        uint64
	workers     uint64
	connections int
	redirects   int
	headers     headers
	laddr       localAddr
	keepalive   bool
	statsd      *statsdOpts
}

type headers struct{ http.Header }

type localAddr struct{ *net.IPAddr }

type statsdOpts struct {
	enable bool
	host   string
	port   uint64
	prefix string
}

func NewVegeta() *Vegeta {
	return &Vegeta{name: "vegeta", running: false,
		opts: &attackOpts{
			headers: headers{http.Header{}},
			laddr:   localAddr{&vegeta.DefaultLocalAddr},
			statsd:  &statsdOpts{}}}
}

type Vegeta struct {
	name    string
	running bool
	params  map[string]interface{}
	opts    *attackOpts
}

func (v *Vegeta) GetName() string {
	return v.name
}

func (v *Vegeta) Start(params map[string]interface{}) {
	defer func() { v.running = false }()
	v.running = true
	v.params = params
	v.decodeParams()

	go v.realStart()
}

func (v *Vegeta) Check() bool {
	return v.running
}

func (v *Vegeta) Params() map[string]interface{} {
	return v.params
}

func (v *Vegeta) realStart() error {
	var (
		tr  vegeta.Targeter
		err error
	)

	src := bytes.NewReader([]byte(v.opts.target))
	body := []byte{}
	if tr, err = vegeta.NewEagerTargeter(src, body, v.opts.headers.Header); err != nil {
		return err
	}

	atk := vegeta.NewAttacker(
		vegeta.Redirects(v.opts.redirects),
		vegeta.Timeout(v.opts.timeout),
		vegeta.LocalAddr(*v.opts.laddr.IPAddr),
		vegeta.Workers(v.opts.workers),
		vegeta.KeepAlive(v.opts.keepalive),
		vegeta.Connections(v.opts.connections),
		vegeta.HTTP2(v.opts.http2),
		vegeta.StatsdEnabled(v.opts.statsd.enable),
		vegeta.StatsdHost(v.opts.statsd.host),
		vegeta.StatsdPort(v.opts.statsd.port),
		vegeta.StatsdPrefix(v.opts.statsd.prefix),
	)

	res := atk.Attack(tr, v.opts.rate, v.opts.duration)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <-sig:
			atk.Stop()
			return nil
		case _, ok := <-res:
			if !ok {
				return nil
			}
		}
	}

}

func (v *Vegeta) decodeParams() {
	jsonStr, _ := json.Marshal(v.params)
	json.Unmarshal(jsonStr, &v.opts)
}
