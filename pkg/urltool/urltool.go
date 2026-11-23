package urltool

import (
	"net/url"
	"path"
)

func GetBasePath(targeturl string) (string, error) {
	myURL, err := url.Parse(targeturl)
	if err != nil {
		return "", err
	}

	basePath := path.Base(myURL.Path)
	return basePath, nil
}
