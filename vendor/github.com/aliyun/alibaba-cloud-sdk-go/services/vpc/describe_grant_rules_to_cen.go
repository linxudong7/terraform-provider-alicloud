package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeGrantRulesToCen invokes the vpc.DescribeGrantRulesToCen API synchronously
// api document: https://help.aliyun.com/api/vpc/describegrantrulestocen.html
func (client *Client) DescribeGrantRulesToCen(request *DescribeGrantRulesToCenRequest) (response *DescribeGrantRulesToCenResponse, err error) {
	response = CreateDescribeGrantRulesToCenResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeGrantRulesToCenWithChan invokes the vpc.DescribeGrantRulesToCen API asynchronously
// api document: https://help.aliyun.com/api/vpc/describegrantrulestocen.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeGrantRulesToCenWithChan(request *DescribeGrantRulesToCenRequest) (<-chan *DescribeGrantRulesToCenResponse, <-chan error) {
	responseChan := make(chan *DescribeGrantRulesToCenResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeGrantRulesToCen(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeGrantRulesToCenWithCallback invokes the vpc.DescribeGrantRulesToCen API asynchronously
// api document: https://help.aliyun.com/api/vpc/describegrantrulestocen.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeGrantRulesToCenWithCallback(request *DescribeGrantRulesToCenRequest, callback func(response *DescribeGrantRulesToCenResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeGrantRulesToCenResponse
		var err error
		defer close(result)
		response, err = client.DescribeGrantRulesToCen(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeGrantRulesToCenRequest is the request struct for api DescribeGrantRulesToCen
type DescribeGrantRulesToCenRequest struct {
	*requests.RpcRequest
	ResourceGroupId      string           `position:"Query" name:"ResourceGroupId"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	InstanceId           string           `position:"Query" name:"InstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	InstanceType         string           `position:"Query" name:"InstanceType"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeGrantRulesToCenResponse is the response struct for api DescribeGrantRulesToCen
type DescribeGrantRulesToCenResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
	PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
	PageSize      int           `json:"PageSize" xml:"PageSize"`
	CenGrantRules CenGrantRules `json:"CenGrantRules" xml:"CenGrantRules"`
}

// CreateDescribeGrantRulesToCenRequest creates a request to invoke DescribeGrantRulesToCen API
func CreateDescribeGrantRulesToCenRequest() (request *DescribeGrantRulesToCenRequest) {
	request = &DescribeGrantRulesToCenRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DescribeGrantRulesToCen", "vpc", "openAPI")
	return
}

// CreateDescribeGrantRulesToCenResponse creates a response to parse from DescribeGrantRulesToCen response
func CreateDescribeGrantRulesToCenResponse() (response *DescribeGrantRulesToCenResponse) {
	response = &DescribeGrantRulesToCenResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
