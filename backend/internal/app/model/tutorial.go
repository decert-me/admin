package model

type Tutorial struct {
	RepoUrl       string `json:"repoUrl,omitempty"`
	Label         string `json:"label,omitempty"`
	CatalogueName string `json:"catalogue_name,omitempty"`
	DocType       string `json:"docusaurus,omitempty"`
	Img           string `json:"img,omitempty"`
	Desc          string `json:"desc,omitempty"`
	Branch        string `json:"branch,omitempty"`
	DocPath       string `json:"docPath,omitempty"`
	StartPage     string `json:"startPage,omitempty"`
	CommitHash    string `json:"commitHash,omitempty"`
	VideoCategory string `json:"videoCategory,omitempty"`
	Sort          string `json:"sort,omitempty"`
	Url           string `json:"url,omitempty"`
}
