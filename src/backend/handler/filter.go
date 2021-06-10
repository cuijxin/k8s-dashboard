package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	utilnet "k8s.io/apimachinery/pkg/util/net"

	"github.com/cuijxin/k8s-dashboard/src/backend/args"
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/xsrftoken"
)

const (
	originalForwardedForHeader = "X-Original-Forwarded-For"
	forwardedForHeader         = "X-Forwarded-For"
	realIPHeader               = "X-Real-Ip"
)

// InstallFilters installs defined filter for given web service.
func InstallFilters(ws *restful.WebService, manager clientapi.clientManager) {
	ws.Filter(requestAndResponseLogger)
	ws.Filter(metricsFilter)
	ws.Filter(validateXSRFFilter(manager.CSRFKey()))
	ws.Filter(restrictedResourcesFilter)
}

// Filter used to restrict access to dashboard exclusive resource,
// i.e. secret used to store dashboard encryption key.
func restrictedResourcesFilter(request *restful.Request,
	response *restful.Response,
	chain *restful.FilterChain) {
	if !authApi.ShouldRejectRequest(request.Request.URL.String()) {
		chain.ProcessFilter(request, response)
		return
	}

	err := errors.NewUnauthorized(errors.MsgDashboardExclusiveResourceError)
	response.WriteHeaderAndEntity(int(err.ErrStatus.Code), err.Error())
}

// web-service filter function used for request and response logging.
func requestAndResponseLogger(request *restful.Request, response *restful.Response,
	chain *restful.FilterChain) {
	if args.Holder.GetAPILogLevel() != "NONE" {
		log.Printf(formatRequestLog(request))
	}

	chain.ProcessFilter(request, response)

	if args.Holder.GetAPILogLevel() != "NONE" {
		log.Printf(formatResponseLog(response, request))
	}
}

// formatRequestLog formats request log string.
func formatRequestLog(request *restful.Request) string {
	uri := ""
	content := "{}"

	if request.Request.URL != nil {
		uri = request.Request.URL.RequestURI()
	}

	byteArr, err := ioutil.ReadAll(request.Request.Body)
	if err == nil {
		content = string(byteArr)
	}

	// Restore request body so we can read it again in regular request handlers
	request.Request.Body = ioutil.NopCloser(bytes.NewReader(byteArr))

	// Is DEBUG level logging enabled? Yes?
	// Great now let's filter out any content from sensitive URLs
	if args.Holder.GetAPILogLevel() != "DEBUG" && checkSensitiveURL(&uri) {
		content = "{ contents hidden }"
	}

	return fmt.Sprintf(RequestLogString, time.Now().Format(time.RFC3339), request.Request.Proto,
		request.Request.Method, uri, getRemoteAddr(request.Request), content)
}

// formatResponseLog formats response log string.
func formatResponseLog(response *restful.Response, request *restful.Request) string {
	return fmt.Sprintf(ResponseLogString, time.Now().Format(time.RFC3339),
		getRemoteAddr(request.Request), response.StatusCode())
}

// checkSensitiveUrl checks if a string matches against a sensitive URL
// true if sensitive. false if not.
func checkSensitiveURL(url *string) bool {
	var s struct{}
	var sensitiveUrls = make(map[string]struct{})
	sensitiveUrls["/api/v1/login"] = s
	sensitiveUrls["/api/v1/csrftoken/login"] = s
	sensitiveUrls["/api/v1/token/refresh"] = s

	if _, ok := sensitiveUrls[*url]; ok {
		return true
	}

	return false
}

func metricsFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resource := mapUrlToResource(req.SelectedRoutePath())
	httpClient := utilnet.GetHTTPClient(req.Request)

	fmt.Println("%s", httpClient)

	chain.ProcessFilter(req, resp)

	if resource != nil {
		// monitor(
		// 	req.Request.Method,
		// 	*resource, httpClient,
		// 	resp.Header().Get("Content-Type"),
		// 	resp.StatusCode(),
		// 	time.Now(),
		// )
	}
}

func validateXSRFFilter(csrfKey string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		resource := mapUrlToResource(req.SelectedRoutePath())

		if resource == nil || (shouldDoCsrfValidation(req) &&
			!xsrftoken.Valid(req.HeaderParameter("X-CSRF-TOKEN"), csrfKey, "none", *resource)) {
			err := errors.NewInvalid("CSRF validation failed")
			log.Print(err)
			resp.AddHeader("Content-Type", "text/plain")
			resp.WriteErrorString(http.StatusUnauthorized, err.Error()+"\n")
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

// Post requests should set correct X-CSRF-TOKEN header, all other requests
// should either not edit anything or be already safe to CSRF attacks (PUT
// and DELETE)
func shouldDoCsrfValidation(req *restful.Request) bool {
	if req.Request.Method != http.MethodPost {
		return false
	}

	// Validation handlers are idempotent functions, and not actual data
	// modification operations
	if strings.HasPrefix(req.SelectedRoutePath(), "/api/v1/appdeployment/validate/") {
		return false
	}

	return true
}

// mapUrlToResource extracts the resource from the URL path /api/v1/<resource>.
// Ignores potential subresources.
func mapUrlToResource(url string) *string {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return nil
	}
	return &parts[3]
}

// getRemoteAddr extracts the remote address of the request, taking into
// account proxy headers.
func getRemoteAddr(r *http.Request) string {
	if ip := getRemoteIPFromForwardHeader(r, originalForwardedForHeader); ip != "" {
		return ip
	}

	if ip := getRemoteIPFromForwardHeader(r, forwardedForHeader); ip != "" {
		return ip
	}

	if realIP := strings.TrimSpace(r.Header.Get(realIPHeader)); realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func getRemoteIPFromForwardHeader(r *http.Request, header string) string {
	ips := strings.Split(r.Header.Get(header), ",")
	return strings.TrimSpace(ips[0])
}
