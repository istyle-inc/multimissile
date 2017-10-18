package server

import (
	"bytes"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/fujiwara/fluent-agent-hydra/ltsv"
	"github.com/istyle-inc/multimissile"
	"github.com/istyle-inc/multimissile/jsonrpc"
	"github.com/istyle-inc/multimissile/wlog"
)

func accessLog(r *http.Request, rr *[]jsonrpc.Request, stime time.Time, status int) {
	etime := time.Now()
	ptime := math.Floor(etime.Sub(stime).Seconds()*1000) / 1000
	records := make(map[string]interface{})

	records["time"] = time.Now().Local().Format("2006/01/02 15:04:05 MST")
	records["addr"] = r.RemoteAddr
	records["status"] = status
	records["ptime"] = ptime
	records["length"] = r.ContentLength
	if msl.Config.LogLevel == "debug" {
		records["headers"] = r.Header
		records["body"] = *rr
	}
	buf := &bytes.Buffer{}
	encoder := ltsv.NewEncoder(buf)
	encoder.Encode(records)
	msl.AL.Out(wlog.Info, strings.TrimRight(buf.String(), "\n"))
}

func errorLog(level wlog.LogLevel, msg string, args ...interface{}) {
	msl.EL.Out(level, msg, args...)
}
