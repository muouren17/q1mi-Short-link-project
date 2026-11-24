package urltool

import (
	"errors"
	"net/url"
	"path"
)

func GetBasePath(targeturl string) (string, error) {
	myURL, err := url.Parse(targeturl)
	if err != nil {
		return "", err
	}
	if len(myURL.Host) == 0 {
		return "", errors.New("no host in targetURL")
	}

	basePath := path.Base(myURL.Path)
	return basePath, nil
}
