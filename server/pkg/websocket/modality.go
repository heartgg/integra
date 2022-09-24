package websocket

type Modality string

const (
	IE     Modality = "IE"
	Fluoro Modality = "Fluoro"
	XRAY   Modality = "XRAY"
	CT     Modality = "CT"
	IR     Modality = "IR"
	MRI    Modality = "MRI"
	US     Modality = "US"
	Dexa   Modality = "Dexa"
	NucMed Modality = "NucMed"
)

var (
	modalityMap = map[string]Modality{
		"IE":     IE,
		"Fluoro": Fluoro,
		"XRAY":   XRAY,
		"CT":     CT,
		"IR":     IR,
		"MRI":    MRI,
		"US":     US,
		"Dexa":   Dexa,
		"NucMed": NucMed,
	}
)

func ParseModality(modality string) Modality {
	return modalityMap[modality]
}