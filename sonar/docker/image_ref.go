package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrImageName = errors.New("Error: Invalid image name.")
	ErrImageTag  = errors.New("Error: Image tag doesn't exist.")
)

type ImageRef struct {
	Namespace string `json:"namesapce"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	ShowTag   bool
}

func (this *ImageRef) String() string {

	if this.ShowTag {
		return fmt.Sprintf("%s/%s:%s", this.Namespace, this.Name, this.Tag)
	}

	return fmt.Sprintf("%s/%s", this.Namespace, this.Name)
}

/* Returns true if the image tag exists for an image. */
func (this *ImageRef) Valid() (bool, error) {

	reqURL := "https://hub.docker.com/v2/repositories/" + this.Namespace + "/" + this.Name + "/tags/" + this.Tag

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return false, err
	}

	resp, err := SendRequest(req, "", "")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else {
		return false, ErrImageTag
	}
}

// NewImageRef creates a new reference to a Docker image.
// Namespace and tag can be empty strings in order to use Docker defaults of 'library' and 'latest'.
func NewImageRef(namespace, name, tag string) (*ImageRef, error) {

	if name == "" {
		return nil, errors.New("The image name must be specified.")
	}

	if namespace == "" {
		namespace = "library"
	}

	if tag == "" {
		tag = "latest"
	}

	return &ImageRef{namespace, name, tag, true}, nil
}

// Convert a string into a full image reference (imageRef).
func ParseImageRef(image string) (*ImageRef, error) {

	var namespace, name, tag string

	// namespace processing
	switch nsIndex := strings.Index(image, "/"); nsIndex {
	case -1:
		// namespace not specified
		name = image
	case 0:
	case len(image) - 1:
		// invalid location
		return nil, ErrImageName
	default:
		namespace = image[0:nsIndex]
		name = image[nsIndex+1:]
	}

	// tag processing
	switch tagIndex := strings.Index(name, ":"); tagIndex {
	case -1:
		//tag not specified
	case 0:
	case len(name) - 1:
		// invalid location
		return nil, ErrImageName
	default:
		tag = name[tagIndex+1 : len(name)]
		name = name[:tagIndex]
	}

	return NewImageRef(namespace, name, tag)
}

func ImagePulls(image string) (uint, error) {

	req, err := http.NewRequest("GET", "https://hub.docker.com/v2/repositories/"+image+"/", nil)
	if err != nil {
		return 0, err
	}

	resp, err := SendRequest(req, "", "")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result2 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result2)

	return uint(result2["pull_count"].(float64)), nil
}

func ImageStars(image string) (uint, error) {

	req, err := http.NewRequest("GET", "https://hub.docker.com/v2/repositories/"+image+"/", nil)
	if err != nil {
		return 0, err
	}

	resp, err := SendRequest(req, "", "")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result2 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result2)

	return uint(result2["star_count"].(float64)), nil
}

func ImageList(namespace string) ([]string, error) {

	req, err := http.NewRequest("GET", "https://hub.docker.com/v2/repositories/"+namespace+"/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := SendRequest(req, "", "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result2 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result2)

	var images []string

	for _, v := range result2["results"].([]interface{}) {
		images = append(images, v.(map[string]interface{})["name"].(string))
	}

	return images, nil
}
