package model

type UserRole int

type ThirdPartyServiceName string

type BlogClassification string

type FavoritesFolderType bool

const (
	SuperAdmin UserRole = 8
	Admin      UserRole = 4
	Vip        UserRole = 2
	Client     UserRole = 1

	ChatGPT ThirdPartyServiceName = "ChatGPT"
	DeepL   ThirdPartyServiceName = "DeepL"
	Judge0  ThirdPartyServiceName = "Judge0"

	ProgrammingLanguages      BlogClassification = "Programming Languages"
	OperatingSystem           BlogClassification = "Operating System"
	Database                  BlogClassification = "Database"
	News                      BlogClassification = "News"
	NetworkSecurity           BlogClassification = "Network Security"
	NMiscellaneousDiscussions BlogClassification = "NMiscellaneous Discussions"
)
