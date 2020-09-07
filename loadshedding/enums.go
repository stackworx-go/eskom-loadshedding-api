package loadshedding

// Province Export
type Province int

const (
	EasternCape Province = iota + 1
	FreeState
	Gauteng
	KwazuluNatal
	Limpop
	Mpumalanga
	NorthWest
	NorthernCape
	WesternCape
)

// Stage Export
type Stage int

const (
	StageUnknown    = -1
	NotLoadShedding = iota + 1
	Stage1
	Stage2
	Stage3
	Stage4
	Stage5
	Stage6
	Stage7
	Stage8
)

// ConvertStage export
func ConvertStage(stage string) Stage {
	switch stage {
	case "1":
		return NotLoadShedding
	case "2":
		return Stage1
	case "3":
		return Stage2
	case "4":
		return Stage3
	case "5":
		return Stage4
	case "6":
		return Stage5
	case "7":
		return Stage6
	case "8":
		return Stage7
	case "9":
		return Stage8
	}

	return StageUnknown
}
