package parser

import "encoding/xml"

// Repomd stores parsed repomd.xml
// Ref: https://blog.packagecloud.io/eng/2015/07/20/yum-repository-internals/
type Repomd struct {
	XMLName    xml.Name      `xml:"repomd"`
	RepomdData []RepomdDatas `xml:"data"`
}

// RepomdDatas stores each element under data
type RepomdDatas struct {
	Type         string `xml:"type,attr"`
	LocationData LocationData
}

// LocationData stores the relative url
type LocationData struct {
	XMLName xml.Name `xml:"location"`
	Href    string   `xml:"href,attr"`
}

// PackageListURI returns the relative path of the gzipped primary.xml.gz
func PackageListURI(data []byte) string {
	var repomd Repomd
	xml.Unmarshal(data, &repomd)

	var path string
	for i := range repomd.RepomdData {
		if repomd.RepomdData[i].Type == "primary" {
			path = repomd.RepomdData[i].LocationData.Href
			break
		}
	}
	return path
}
