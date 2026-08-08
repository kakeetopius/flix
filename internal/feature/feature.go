// Package feature defines the go types for the different features ie movie/tv show
package feature

type Type int

const (
	FeatureTypeMovie Type = iota
	FeatureTypeSerie
)

func (t Type) String() string {
	switch t {
	case FeatureTypeMovie:
		return "movie"
	case FeatureTypeSerie:
		return "serie"
	}

	return ""
}
