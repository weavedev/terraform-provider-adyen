package provider

import (
	"bytes"
	"os"
	"testing"
	"text/template"
)

func testProviderClientFromTmpl(t *testing.T) string {
	tmplString := `
	provider "adyen" {
		api_key = "{{.ApiKey}}"
		environment = "{{.Environment}}"
		merchant_account = "{{.MerchantAccount}}"
	}

	`

	tmpl, err := template.New("providerClient").Parse(tmplString)
	if err != nil {
		t.Fatal("could not create template string from env vars")
	}
	varMap := map[string]interface{}{
		"ApiKey":          os.Getenv("ADYEN_API_KEY"),
		"Environment":     os.Getenv("ADYEN_API_ENVIRONMENT"),
		"MerchantAccount": os.Getenv("ADYEN_API_MERCHANT_ACCOUNT"),
	}

	var renderedConfig bytes.Buffer
	if err := tmpl.Execute(&renderedConfig, varMap); err != nil {
		panic("Failed to render template: " + err.Error())
	}

	return renderedConfig.String()
}
