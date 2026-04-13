package sneaky

import "sneaky-core/internal/config"

type ConfigMetadata struct {
	Format    string
	AdapterID string
}

func InspectConfigPath(path string) (ConfigMetadata, error) {
	input, err := config.LoadFile(path)
	if err != nil {
		return ConfigMetadata{}, err
	}

	metadata, err := config.DetectAndValidate(input)
	if err != nil {
		return ConfigMetadata{}, err
	}

	return ConfigMetadata{
		Format:    string(metadata.Format),
		AdapterID: metadata.AdapterID,
	}, nil
}

func InspectConfigBytes(raw []byte) (ConfigMetadata, error) {
	input, err := config.Parse(raw, "inline")
	if err != nil {
		return ConfigMetadata{}, err
	}

	metadata, err := config.DetectAndValidate(input)
	if err != nil {
		return ConfigMetadata{}, err
	}

	return ConfigMetadata{
		Format:    string(metadata.Format),
		AdapterID: metadata.AdapterID,
	}, nil
}
