package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setFromEnvOrDefaults(e *ConfigVars) {

	e.set(&e.Tags, "PROBR_TAGS", "")
	e.set(&e.AuditEnabled, "PROBR_AUDIT_ENABLED", "true")
	e.set(&e.OutputType, "PROBR_OUTPUT_TYPE", "IO")
	e.set(&e.CucumberDir, "PROBR_CUCUMBER_DIR", "cucumber_output")
	e.set(&e.AuditDir, "PROBR_AUDIT_DIR", "audit_output")
	e.set(&e.LogLevel, "PROBR_LOG_LEVEL", "ERROR")
	e.set(&e.OverwriteHistoricalAudits, "OVERWRITE_AUDITS", "true")

	e.set(&e.ServicePacks.Kubernetes.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	e.set(&e.ServicePacks.Kubernetes.KubeContext, "KUBE_CONTEXT", "")
	e.set(&e.ServicePacks.Kubernetes.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
	e.set(&e.ServicePacks.Kubernetes.AuthorisedContainerRegistry, "PROBR_AUTHORISED_REGISTRY", "docker.io")
	e.set(&e.ServicePacks.Kubernetes.UnauthorisedContainerRegistry, "PROBR_UNAUTHORISED_REGISTRY", "")
	e.set(&e.ServicePacks.Kubernetes.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")

	e.set(&e.CloudProviders.Azure.ClientID, "AZURE_CLIENT_ID", "")
	e.set(&e.CloudProviders.Azure.ClientSecret, "AZURE_CLIENT_SECRET", "")
	e.set(&e.CloudProviders.Azure.TenantID, "AZURE_TENANT_ID", "")
	e.set(&e.CloudProviders.Azure.LocationDefault, "AZURE_LOCATION_DEFAULT", "")
	e.set(&e.CloudProviders.Azure.Identity.DefaultNamespaceAI, "DEFAULT_NS_AZURE_IDENTITY", "probr-defaultns-ai")
	e.set(&e.CloudProviders.Azure.Identity.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-defaultns-aib")
}

func getDefaultKubeConfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// set fetches the env var or sets the default value as needed for the specified field from ConfigVars
func (e *ConfigVars) set(field interface{}, varName string, defaultValue interface{}) {
	switch v := field.(type) {
	default:
		log.Fatalf("unexpected type for %v, %T", varName, v)
	case *string:
		if *field.(*string) == "" {
			*field.(*string) = os.Getenv(varName)
		}
		if *field.(*string) == "" {
			*field.(*string) = defaultValue.(string)
		}
	case *[]string:
		if len(*field.(*[]string)) == 0 {
			t := os.Getenv(varName) // if []string, env var should be comma separated values
			if len(t) > 0 {
				*field.(*[]string) = append(*field.(*[]string), strings.Split(t, ",")...)
			}
		}
		if len(*field.(*[]string)) == 0 {
			*field.(*[]string) = defaultValue.([]string)
		}
	}

}
