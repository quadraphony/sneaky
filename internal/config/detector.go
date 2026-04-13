package config

import "sort"

var singboxHintKeys = map[string]struct{}{
	"dns":          {},
	"endpoints":    {},
	"experimental": {},
	"inbounds":     {},
	"log":          {},
	"outbounds":    {},
	"route":        {},
}

func Detect(input Input) (Metadata, error) {
	object, err := decodeJSONObject(input)
	if err != nil {
		return Metadata{}, err
	}

	keys := make([]string, 0, len(object))
	singboxHints := 0
	for key := range object {
		keys = append(keys, key)
		if _, ok := singboxHintKeys[key]; ok {
			singboxHints++
		}
	}
	sort.Strings(keys)

	if _, ok := object["outbounds"]; ok {
		return Metadata{
			Format:       FormatJSON,
			AdapterID:    AdapterSingbox,
			TopLevelKeys: keys,
		}, nil
	}

	if singboxHints > 0 {
		return Metadata{}, &ValidationError{
			Code:    ErrCodeAmbiguousFormat,
			Source:  input.Source,
			Message: "config has sing-box-like keys but no definitive sing-box outbound definition",
		}
	}

	return Metadata{}, &ValidationError{
		Code:    ErrCodeUnknownFormat,
		Source:  input.Source,
		Message: "config format is not recognized within the current detection scope",
	}
}
