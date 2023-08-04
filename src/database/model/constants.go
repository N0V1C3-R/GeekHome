package model

type UserRole string

type ThirdPartyServiceName string

type BlogClassification string

const (
	SuperAdmin UserRole = "SUPER_ADMIN"
	Admin      UserRole = "ADMIN"
	Client     UserRole = "USER"

	ChatGPT ThirdPartyServiceName = "ChatGPT"
	DeepL   ThirdPartyServiceName = "DeepL"
	Judge0  ThirdPartyServiceName = "Judge0"
)
