package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/cuijxin/k8s-dashboard/src/backend/args"
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/client"

	"github.com/spf13/pflag"
)

var (
	argInsecurePort        = pflag.Int("insecure-port", 9090, "The port to listen to for incoming HTTP requests.")
	argPort                = pflag.Int("port", 8443, "The secure port to listen to for incoming HTTPS requests.")
	argInsecureBindAddress = pflag.IP("insecure-bind-address", net.IPv4(127, 0, 0, 1), "The IP address on which to serve the --insecure-port (set to 127.0.0.1 for all interfaces).")
	argBindAddress         = pflag.IP("bind-address", net.IPv4(0, 0, 0, 0), "The IP address on which to serve the --port (set to 0.0.0.0 for all interfaces).")
	argDefaultCertDir      = pflag.String("default-cert-dir", "/certs", "Directory path containing '--tls-cert-file' and '--tls-key-file' files. Used also when auto-generating certificates flag is set.")
	argCertFile            = pflag.String("tls-cert-file", "", "File containing the default x509 Certificate for HTTPS.")
	argKeyFile             = pflag.String("tls-key-file", "", "File containing the default x509 private key matching --tls-cert-file.")
	argApiserverHost       = pflag.String("apiserver-host", "", "The address of the Kubernetes Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8080. If not specified, the assumption is that the binary runs inside a "+
		"Kubernetes cluster and local discovery is attempted.")
	argMetricsProvider = pflag.String("metrics-provider", "sidecar", "Select provider type for metrics. 'none' will not check metrics.")
	argHeapsterHost    = pflag.String("heapster-host", "", "The address of the Heapster Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8082. If not specified, the assumption is that the binary runs inside a "+
		"Kubernetes cluster and service proxy will be used.")
	argSidecarHost = pflag.String("sidecar-host", "", "The address of the Sidecar Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8000. If not specified, the assumption is that the binary runs inside a "+
		"Kubernetes cluster and service proxy will be used.")
	argKubeConfigFile     = pflag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	argTokenTTL           = pflag.Int("token-ttl", int(authApi.DefaultTokenTTL), "Expiration time (in seconds) of JWE tokens generated by dashboard. '0' never expires")
	argAuthenticationMode = pflag.StringSlice("authentication-mode", []string{authApi.Token.String()}, "Enables authentication options that will be reflected on login screen. Supported values: token, basic. "+
		"Note that basic option should only be used if apiserver has '--authorization-mode=ABAC' and '--basic-auth-file' flags set.")
	argMetricClientCheckPeriod   = pflag.Int("metric-client-check-period", 30, "Time in seconds that defines how often configured metric client health check should be run.")
	argAutoGenerateCertificates  = pflag.Bool("auto-generate-certificates", false, "When set to true, Dashboard will automatically generate certificates used to serve HTTPS. (default false)")
	argEnableInsecureLogin       = pflag.Bool("enable-insecure-login", false, "When enabled, Dashboard login view will also be shown when Dashboard is not served over HTTPS. (default false)")
	argEnableSkip                = pflag.Bool("enable-skip-login", false, "When enabled, the skip button on the login page will be shown. (default false)")
	argSystemBanner              = pflag.String("system-banner", "", "When non-empty displays message to Dashboard users. Accepts simple HTML tags.")
	argSystemBannerSeverity      = pflag.String("system-banner-severity", "INFO", "Severity of system banner. Should be one of 'INFO|WARNING|ERROR'.")
	argAPILogLevel               = pflag.String("api-log-level", "INFO", "Level of API request logging. Should be one of 'INFO|NONE|DEBUG'.")
	argDisableSettingsAuthorizer = pflag.Bool("disable-settings-authorizer", false, "When enabled, Dashboard settings page will not require user to be logged in and authorized to access settings page. (default false)")
	argNamespace                 = pflag.String("namespace", getEnv("POD_NAMESPACE", "kube-system"), "When non-default namespace is used, create encryption key in the specified namespace.")
	localeConfig                 = pflag.String("locale-config", "./locale_conf.json", "File containing the configuration of locales")
)

func main() {
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = flag.CommandLine.Parse(make([]string, 0))

	initArgHolder()

	if args.Holder.GetApiServerHost() != "" {
		log.Printf("Using apiserver-host location: %s", args.Holder.GetApiServerHost())
	}
	if args.Holder.GetKubeConfigFile() != "" {
		log.Printf("Using kubeconfig file: %s", args.Holder.GetKubeConfigFile())
	}
	if args.Holder.GetNamespace() != "" {
		log.Printf("Using namespace: %s", args.Holder.GetNamespace())
	}

	clientManager := client.NewClientManager(args.Holder.GetKubeConfigFile(), args.Holder.GetApiServerHost())
	versionInfo, err := clientManager.InsecureClient().Discovery().ServerVersion()
	if err != nil {
		handleFatalInitError(err)
	}

	log.Printf("Successful initial request to the apiserver, version: %s", versionInfo.String())
}

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetInsecurePort(*argInsecurePort)
	builder.SetPort(*argPort)
	builder.SetTokenTTL(*argTokenTTL)
	builder.SetMetricClientCheckPeriod(*argMetricClientCheckPeriod)
	builder.SetInsecureBindAddress(*argInsecureBindAddress)
	builder.SetBindAddress(*argBindAddress)
	builder.SetDefaultCertDir(*argDefaultCertDir)
	builder.SetCertFile(*argCertFile)
	builder.SetKeyFile(*argKeyFile)
	builder.SetApiServerHost(*argApiserverHost)
	builder.SetMetricsProvider(*argMetricsProvider)
	builder.SetHeapsterHost(*argHeapsterHost)
	builder.SetSidecarHost(*argSidecarHost)
	builder.SetKubeConfigFile(*argKubeConfigFile)
	builder.SetSystemBanner(*argSystemBanner)
	builder.SetSystemBannerSeverity(*argSystemBannerSeverity)
	builder.SetAPILogLevel(*argAPILogLevel)
	builder.SetAuthenticationMode(*argAuthenticationMode)
	builder.SetAutoGenerateCertificates(*argAutoGenerateCertificates)
	builder.SetDisableSettingsAuthorizer(*argDisableSettingsAuthorizer)
	builder.SetEnableSkipLogin(*argEnableSkip)
	builder.SetNamespace(*argNamespace)
	builder.SetLocaleConfig(*localeConfig)
}

func handleFatalInitError(err error) {
	log.Fatalf("Error while initializing connection to Kubernetes apiserver."+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service account's configuration) or the "+
		"--apiserver-host param points to a server that does not exist. Reason: %s\n"+
		"Refer to our FAQ and wiki pages for more information: "+
		"https://github.com/kubernetes/dashboard/wiki/FAQ", err)
}

func handleFatalInitServingCertError(err error) {
	log.Fatalf("Error while loading dashboard server certificates. Reason: %s", err)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}
