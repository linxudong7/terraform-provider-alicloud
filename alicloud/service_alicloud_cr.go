package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CrService struct {
	client *connectivity.AliyunClient
}

type crCreateNamespaceRequestPayload struct {
	Namespace struct {
		Namespace string `json:"Namespace"`
	} `json:"Namespace"`
}

type crUpdateNamespaceRequestPayload struct {
	Namespace struct {
		AutoCreate        bool   `json:"AutoCreate"`
		DefaultVisibility string `json:"DefaultVisibility"`
	} `json:"Namespace"`
}

type crDescribeNamespaceResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace struct {
			Namespace         string `json:"namespace"`
			AuthorizeType     string `json:"authorizeType"`
			DefaultVisibility string `json:"defaultVisibility"`
			AutoCreate        bool   `json:"autoCreate"`
			NamespaceStatus   string `json:"namespaceStatus"`
		} `json:"namespace"`
	} `json:"data"`
}

type crDescribeNamespaceListResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace []struct {
			Namespace       string `json:"namespace"`
			AuthorizeType   string `json:"authorizeType"`
			NamespaceStatus string `json:"namespaceStatus"`
		} `json:"namespaces"`
	} `json:"data"`
}

const (
	RepoTypePublic  = "PUBLIC"
	RepoTypePrivate = "PRIVATE"
)

type crCreateRepoRequestPayload struct {
	Repo struct {
		RepoNamespace string `json:"RepoNamespace"`
		RepoName      string `json:"RepoName"`
		Summary       string `json:"Summary"`
		Detail        string `json:"Detail"`
		RepoType      string `json:"RepoType"`
	} `json:"Repo"`
}

type crUpdateRepoRequestPayload struct {
	Repo struct {
		Summary  string `json:"Summary"`
		Detail   string `json:"Detail"`
		RepoType string `json:"RepoType"`
	} `json:"Repo"`
}

type crDescribeRepoResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repo struct {
			Summary        string `json:"summary"`
			Detail         string `json:"detail"`
			RepoNamespace  string `json:"repoNamespace"`
			RepoName       string `json:"repoName"`
			RepoType       string `json:"repoType"`
			RepoDomainList struct {
				Public   string `json:"public"`
				Internal string `json:"internal"`
				Vpc      string `json:"vpc"`
			}
		} `json:"repo"`
	} `json:"data"`
}

type crDescribeReposResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repos    []crRepo `json:"repos"`
		Total    int      `json:"total"`
		PageSize int      `json:"pageSize"`
		Page     int      `json:"page"`
	} `json:"data"`
}

type crRepo struct {
	Summary        string `json:"summary"`
	RepoNamespace  string `json:"repoNamespace"`
	RepoName       string `json:"repoName"`
	RepoType       string `json:"repoType"`
	RegionId       string `json:"regionId"`
	RepoDomainList struct {
		Public   string `json:"public"`
		Internal string `json:"internal"`
		Vpc      string `json:"vpc"`
	} `json:"repoDomainList"`
}

type crDescribeRepoTagsResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Tags     []crTag `json:"tags"`
		Total    int     `json:"total"`
		PageSize int     `json:"pageSize"`
		Page     int     `json:"page"`
	} `json:"data"`
}

type crTag struct {
	ImageId     string `json:"imageId"`
	Digest      string `json:"digest"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	ImageUpdate int    `json:"imageUpdate"`
	ImageCreate int    `json:"imageCreate"`
	ImageSize   int    `json:"imageSize"`
}

func (c *CrService) DescribeCrNamespace(id string) (*cr.GetNamespaceResponse, error) {
	request := cr.CreateGetNamespaceRequest()
	request.Namespace = id

	var response *cr.GetNamespaceResponse

	var err error
	raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.GetNamespace(request)
	})
	if err != nil {
		if IsExceptedError(err, ErrorNamespaceNotExist) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*cr.GetNamespaceResponse)

	return response, nil
}

func (c *CrService) WaitForCRNamespace(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := c.DescribeCrNamespace(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		var response crDescribeNamespaceResponse
		err = json.Unmarshal(object.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Data.Namespace.Namespace == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, response.Data.Namespace.Namespace, id, ProviderERROR)
		}
	}
}

func (c *CrService) DescribeCrRepo(id string) (*cr.GetRepoResponse, error) {
	sli := strings.Split(id, SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]

	request := cr.CreateGetRepoRequest()
	request.RepoNamespace = repoNamespace
	request.RepoName = repoName

	raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.GetRepo(request)
	})
	response, _ := raw.(*cr.GetRepoResponse)
	if err != nil {
		if IsExceptedError(err, ErrorRepoNotExist) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return response, nil
}

func (c *CrService) WaitForCrRepo(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := c.DescribeCrRepo(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		var response crDescribeRepoResponse
		err = json.Unmarshal(object.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		respId := response.Data.Repo.RepoNamespace + SLASH_SEPARATED + response.Data.Repo.RepoName
		if respId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, respId, id, ProviderERROR)
		}
	}
}
