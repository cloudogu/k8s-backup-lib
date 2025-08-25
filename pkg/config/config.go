package config

const (
	StageDevelopment = "development"
	StageProduction  = "production"
)

var Stage = StageProduction

func IsStageDevelopment() bool {
	return Stage == StageDevelopment
}
