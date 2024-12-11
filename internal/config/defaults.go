package config

// Sync, editor, plugin configs can use default type values
func defaultConfig() OolongConfig {
	return OolongConfig{
		NotesDirPaths:     []string{},
		AllowedExtensions: []string{},
		IgnoreDirectories: []string{},
		OpenCommand:       []string{},
		LinkerConfig:      defaultLinkerConfig(),
		GraphConfig:       defaultGraphConfig(),
		SyncConfig:        OolongSyncConfig{},
		PluginsConfig:     OolongPluginConfig{},
	}
}

func defaultLinkerConfig() OolongLinkerConfig {
	return OolongLinkerConfig{
		NGramRange: []int{1, 2, 3},
		StopWords:  []string{},
	}
}

func defaultGraphConfig() OolongGraphConfig {
	// TEST: these values with new users (setting these dynamically would likely be ideal using some sort of percentages)
	return OolongGraphConfig{
		MinNodeWeight: 2.0,
		MaxNodeWeight: 10.0,
		MinLinkWeight: 0.1,
		DefaultMode:   "2d",
	}
}
