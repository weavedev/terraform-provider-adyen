package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mapWebhooksAdditionalSettingsEventCodes(input []string) []attr.Value {
	output := make([]attr.Value, 0, len(input))
	for _, code := range input {
		output = append(output, types.StringValue(code))
	}
	return output
}

func mapWebhooksAdditionalSettingsProperties(input map[string]bool) map[string]attr.Value {
	output := make(map[string]attr.Value)
	for k, v := range input {
		output[k] = types.BoolValue(v)
	}
	return output
}

func mapWebhooksLinks(self *string, generateHmac *string, merchant *string, testWebhook *string) map[string]attr.Value {
	return map[string]attr.Value{
		"self": types.ObjectValueMust(map[string]attr.Type{
			"href": types.StringType,
		}, map[string]attr.Value{
			"href": types.StringPointerValue(self),
		}),
		"generate_hmac": types.ObjectValueMust(map[string]attr.Type{
			"href": types.StringType,
		}, map[string]attr.Value{
			"href": types.StringPointerValue(generateHmac),
		}),
		"merchant": types.ObjectValueMust(map[string]attr.Type{
			"href": types.StringType,
		}, map[string]attr.Value{
			"href": types.StringPointerValue(merchant),
		}),
		"test_webhook": types.ObjectValueMust(map[string]attr.Type{
			"href": types.StringType,
		}, map[string]attr.Value{
			"href": types.StringPointerValue(testWebhook),
		}),
	}
}

var linksAttributeMap = map[string]attr.Type{
	"self": types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"href": types.StringType,
		},
	},
	"generate_hmac": types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"href": types.StringType,
		},
	},
	"merchant": types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"href": types.StringType,
		},
	},
	"test_webhook": types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"href": types.StringType,
		},
	},
}

var additionalSettingsAttributeMap = map[string]attr.Type{
	"include_event_codes": types.ListType{
		ElemType: types.StringType,
	},
	"exclude_event_codes": types.ListType{
		ElemType: types.StringType,
	},
	"properties": types.MapType{
		ElemType: types.BoolType,
	},
}
