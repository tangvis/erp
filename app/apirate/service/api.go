package service

type APP interface {
	Allow(userID, path string) bool
	InitPublic(publicLimitSetting map[string]int)
}
