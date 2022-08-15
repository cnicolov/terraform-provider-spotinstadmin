package main

const (
	providerTokenAttrKey    = "token"
	providerEmailAttrKey    = "email"
	providerPasswordAttrKey = "password"
)

const (
	envSpotinstTokenKey    = "SPOTINST_TOKEN"
	envSpotinstEmailKey    = "SPOTINST_EMAIL"
	envSpotinstPasswordKey = "SPOTINST_PASSWORD"
)

const (
	providerName                 = "spotinstadmin"
	accountResourceName          = providerName + "_account"
	linkAccountResourceName      = providerName + "_account_aws_link"
	programmaticUserResourceName = providerName + "_programmatic_user"
)

const (
	accountResourceNameAttrKey               = "name"
	accountResourceExternalIdAttrKey         = "external_id"
	accountResourceProviderExternalIdAttrKey = "provider_external_id"
	accountResourceOrganizationIdAttrKey     = "organization_id"
)

const (
	linkAccountResourceAccountIDAttrKey          = "account_id"
	linkAccountResourceRoleArnAttrKey            = "aws_role_arn"
	linkAccountResourceProviderExternalIdAttrKey = "provider_external_id"
	linkAccountResourceOrganizationIdAttrKey     = "organization_id"
)

const (
	userResourceAccountIDAttrKey   = "account_id"
	userResourceNameAttrKey        = "name"
	userResourceDescriptionAttrKey = "description"
	userResourceAccessTokenAttrKey = "access_token"
)
