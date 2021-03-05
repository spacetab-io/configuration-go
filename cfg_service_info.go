package config

type ApplicationInfo struct {
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Version   string `json:"version"`
	About     string `json:"about"`
	Docs      string `json:"docs"`
	Contacts  string `json:"contacts"`
	Copyright string `json:"copyright"`
}
