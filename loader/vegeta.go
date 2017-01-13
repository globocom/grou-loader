package loader

import (
	"bytes"
	"fmt"
	"github.com/globocom/vegeta/lib"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	name       string
	running    bool
	params     map[string]interface{}
	opts       *attackOpts
}

func (v *Vegeta) GetName() string {
	return v.name
}

func (v *Vegeta) Start(params map[string]interface{}) {
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
	defer func() { v.running = false }()

	var (
		tr  vegeta.Targeter
		err error
	)

	src := bytes.NewReader([]byte(v.opts.target))
	body := []byte{}
	if tr, err = vegeta.NewEagerTargeter(src, body, v.opts.headers.Header); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return err
	}

	atk := vegeta.NewAttacker(
		vegeta.Redirects(v.opts.redirects),
		vegeta.Timeout(v.opts.timeout),
		vegeta.Workers(v.opts.workers),
		vegeta.KeepAlive(v.opts.keepalive),
		vegeta.Connections(v.opts.connections),
		vegeta.HTTP2(v.opts.http2),
		vegeta.StatsdEnabled(v.opts.statsd.enable),
		vegeta.StatsdHost(v.opts.statsd.host),
		vegeta.StatsdPort(v.opts.statsd.port),
		vegeta.StatsdPrefix(v.opts.statsd.prefix),
	)

	v.running = true

	res := atk.Attack(tr, v.opts.rate, v.opts.duration)

	fmt.Fprint(os.Stderr, "ATTACK\n")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <-sig:
			fmt.Fprint(os.Stderr, "STOPPED\n")
			atk.Stop()
			v.running = false
			return nil
		case _, ok := <-res:
			if !ok {
				return nil
			}
		}
	}
}

func (v *Vegeta) decodeParams() {
	v.opts.target, _ = v.params["target"].(string)
	v.opts.http2, _ = v.params["http2"].(bool)
	durationFloat, _ := v.params["duration"].(float64)
	v.opts.duration = time.Duration(durationFloat) * time.Second
	timeoutFloat, _ := v.params["timeout"].(float64)
	v.opts.timeout = time.Duration(timeoutFloat) * time.Second
	rateFloat, _ := v.params["rate"].(float64)
	v.opts.rate = uint64(rateFloat)
	workersFloat, _ := v.params["workers"].(float64)
	v.opts.workers = uint64(workersFloat)
	connectionsFloat, _ := v.params["connections"].(float64)
	v.opts.connections = int(connectionsFloat)
	redirectsFloat, _ := v.params["redirects"].(float64)
	v.opts.redirects = int(redirectsFloat)
	v.opts.headers, _ = v.params["headers"].(headers)
	v.opts.laddr, _ = v.params["laddr"].(localAddr)
	v.opts.keepalive, _ = v.params["keepalive"].(bool)
	v.opts.statsd, _ = v.params["statsd"].(statsdOpts)
}
