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
	programmaticUserResourceName = providerName + "_programmatic_user"
)

const (
	accountResourceNameAttrKey       = "name"
	accountResourceRoleArnAttrKey    = "aws_role_arn"
	accountResourceExternalIDAttrKey = "aws_external_id"
)

const (
	userResourceAccountIDAttrKey   = "account_id"
	userResourceNameAttrKey        = "name"
	userResourceDescriptionAttrKey = "description"
	userResourceAccessTokenAttrKey = "access_token"
)
