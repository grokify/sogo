package httpsimple

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apex/gateway"
	"github.com/buaazp/fasthttprouter"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/sogo/net/http/anyhttp"
	"github.com/valyala/fasthttp"
)

const (
	EngineAWSLambda = "awslambda"
	EngineNetHTTP   = "nethttp"
	EngineFastHTTP  = "fasthttp"
)

type SimpleServer interface {
	PortInt() int
	HTTPEngine() string
	Router() http.Handler
	RouterFast() *fasthttprouter.Router
}

func Serve(svc SimpleServer) {
	engine := strings.ToLower(strings.TrimSpace(svc.HTTPEngine()))
	if len(engine) == 0 {
		engine = EngineNetHTTP
	}
	switch engine {
	case EngineNetHTTP:
		svr := httputilmore.NewServerTimeouts(portAddress(svc.PortInt()), svc.Router(), 1*time.Second)
		log.Fatal(svr.ListenAndServe())
	case EngineAWSLambda:
		log.Fatal(
			gateway.ListenAndServe(
				portAddress(svc.PortInt()),
				svc.Router()))
	case EngineFastHTTP:
		router := svc.RouterFast()
		if router == nil {
			log.Fatal("E_NO_FASTROUTER_FOR_ENGINE_FASTHTTP")
		}
		log.Fatal(
			fasthttp.ListenAndServe(
				portAddress(svc.PortInt()),
				router.Handler))
	default:
		log.Fatalf("E_ENGINE_NOT_FOUND [%s]", engine)
	}
}

func portAddress(port int) string { return ":" + strconv.Itoa(port) }

type TestResponse struct {
	Command    string    `json:"command"`
	RequestURL string    `json:"requestURL"`
	Time       time.Time `json:"time"`
}

func HandleTestFastHTTP(ctx *fasthttp.RequestCtx) {
	HandleTestAnyEngine(anyhttp.NewResReqFastHTTP(ctx))
}

func HandleTestNetHTTP(res http.ResponseWriter, req *http.Request) {
	HandleTestAnyEngine(anyhttp.NewResReqNetHTTP(res, req))
}

func HandleTestAnyEngine(aRes anyhttp.Response, aReq anyhttp.Request) {
	test := TestResponse{
		Command:    "pong",
		RequestURL: string(aReq.RequestURI()),
		Time:       time.Now().UTC()}
	bytes, _ := json.Marshal(test)
	_, err := aRes.SetBodyBytes(bytes)
	if err != nil {
		aRes.SetStatusCode(http.StatusInternalServerError)
	} else {
		aRes.SetStatusCode(http.StatusOK)
	}
}

type Handler interface {
	HandleNetHTTP(res http.ResponseWriter, req *http.Request)
	HandleFastHTTP(ctx *fasthttp.RequestCtx)
	HandleAnyHTTP(aRes anyhttp.Response, aReq anyhttp.Request)
}
