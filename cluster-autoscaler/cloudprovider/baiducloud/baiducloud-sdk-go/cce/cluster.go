/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	klog "k8s.io/klog/v2"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/baiducloud/baiducloud-sdk-go/bce"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/baiducloud/baiducloud-sdk-go/util"
)

const (
	// InstanceStatusRunning status
	InstanceStatusRunning string = "RUNNING"
	// InstanceStatusCreating status
	InstanceStatusCreating string = "CREATING"
	// InstanceStatusDeleting status
	InstanceStatusDeleting string = "DELETING"
	// InstanceStatusDeleted status
	InstanceStatusDeleted string = "DELETED"
	// InstanceStatusCreateFailed status
	InstanceStatusCreateFailed string = "CREATE_FAILED"
	// InstanceStatusError status
	InstanceStatusError string = "ERROR"
)

// NodeConfig is the config for node
type NodeConfig struct {
	InstanceType int    `json:"instanceType"`
	CPU          int    `json:"cpu,omitempty"`
	Memory       int    `json:"memory,omitempty"`
	GpuCount     int    `json:"gpuCount,omitempty"`
	GpuCard      string `json:"gpuCard,omitempty"`
	DiskSize     int    `json:"diskSize,omitempty"`
	GroupID      string `json:"groupID"`
}

// CceCluster defines cluster of cce
type CceCluster struct {
	ClusterUuid string     `json:"clusterUuid"`
	NodeConfig  NodeConfig `json:"nodeConfig"`
}

// CceGroup defines autoscaling group
type CceGroup struct {
	InstanceType     int    `json:"instanceType"`
	CPU              int    `json:"cpu,omitempty"`
	Memory           int    `json:"memory,omitempty"`
	GpuCount         int    `json:"gpuCount,omitempty"`
	GpuCard          string `json:"gpuCard,omitempty"`
	DiskSize         int    `json:"diskSize,omitempty"`
	EphemeralStorage int    `json:"ephemeralStorage,omitempty"`
	Tags             []Tag  `json:"tags"`
}

// Tag defines label
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"Value"`
}

// DescribeCluster describe the cluster
func (c *Client) DescribeCluster(clusterID string) (*CceCluster, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/cluster/"+clusterID, nil), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}

	var cceCluster CceCluster
	err = json.Unmarshal(bodyContent, &cceCluster)
	if err != nil {
		return nil, err
	}

	return &cceCluster, nil
}

// CceInstance define instance of cce
type CceInstance struct {
	InstanceId            string `json:"id"`
	InstanceName          string `json:"name"`
	Description           string `json:"desc"`
	Status                string `json:"status"`
	PaymentTiming         string `json:"paymentTiming"`
	CreationTime          string `json:"createTime"`
	ExpireTime            string `json:"expireTime"`
	PublicIP              string `json:"publicIp"`
	InternalIP            string `json:"internalIp"`
	CpuCount              int    `json:"cpuCount"`
	GpuCount              int    `json:"gpuCount"`
	MemoryCapacityInGB    int    `json:"memory"`
	LocalDiskSizeInGB     int    `json:"localDiskSizeInGB"`
	ImageId               string `json:"imageId"`
	NetworkCapacityInMbps int    `json:"networkCapacityInMbps"`
	PlacementPolicy       string `json:"placementPolicy"`
	ZoneName              string `json:"zoneName"`
	SubnetId              string `json:"subnetId"`
	VpcId                 string `json:"vpcId"`
}

// ListInstancesResponse define response of cce list
type ListInstancesResponse struct {
	Instances []CceInstance `json:"instanceList"`
}

// ListInstances gets all Instances of a cluster.
func (c *Client) ListInstances(clusterID string) ([]CceInstance, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}
	params := map[string]string{
		"clusterid": clusterID,
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/instance", params), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var insList ListInstancesResponse
	err = json.Unmarshal(bodyContent, &insList)

	if err != nil {
		return nil, err
	}

	return insList.Instances, nil
}

// CdsPreMountInfo define premount
type CdsPreMountInfo struct {
	MountPath string           `json:"mountPath"`
	CdsConfig []DiskSizeConfig `json:"cdsConfig"`
}

// DiskSizeConfig define distsize config
type DiskSizeConfig struct {
	Size         string `json:"size"`
	SnapshotID   string `json:"snapshotID"`
	SnapshotName string `json:"snapshotName"`
	VolumeType   string `json:"volumeType"`
	StorageType  string `json:"storageType"`
	LogicalZone  string `json:"logicalZone"`
}

// ScaleUpClusterArgs define  args
type ScaleUpClusterArgs struct {
	ClusterID       string          `json:"clusterUuid,omitempty"`
	CdsPreMountInfo CdsPreMountInfo `json:"cdsPreMountInfo,omitempty"`
	OrderContent    OrderContent    `json:"orderContent,omitempty"`
}

// ScaleUpClusterResponse define  args
type ScaleUpClusterResponse struct {
	ClusterID string   `json:"clusterUuid"`
	OrderID   []string `json:"orderId"`
}

// OrderContent define  bcc order content
type OrderContent struct {
	PaymentMethod []string    `json:"paymentMethod,omitempty"`
	Items         []OrderItem `json:"items,omitempty"`
}

// OrderItem define  bcc order content item
type OrderItem struct {
	Config        interface{} `json:"config,omitempty"`
	PaymentMethod []string    `json:"paymentMethod,omitempty"`
}

// BccOrderConfig define BCC order config
type BccOrderConfig struct {
	// ??????????????????????????????postpay
	ProductType string `json:"productType,omitempty"`
	Region      string `json:"region,omitempty"`
	LogicalZone string `json:"logicalZone,omitempty"`
	// BCC???????????????omitempty
	InstanceType int `json:"instanceType"`
	// ???????????????????????? ?????????????????????
	FpgaCard string `json:"fpgaCard,omitempty"`
	GpuCard  int    `json:"gpuCard,omitempty"`
	GpuCount int    `json:"gpuCount,omitempty"`

	CPU    int `json:"cpu,omitempty"`
	Memory int `json:"memory,omitempty"`
	// ??????????????? ubuntu1604
	ImageType string `json:"imageType,omitempty"`
	// ????????????
	OsType string `json:"osType,omitempty"`
	// ????????????
	OsVersion string `json:"osVersion,omitempty"`
	// ???????????????
	DiskSize int `json:"diskSize"`
	// ????????????
	EbsSize []int `json:"ebsSize,omitempty"`
	// ??????????????????EIP
	IfBuyEip int `json:"ifBuyEip,omitempty"`
	// eip??????
	EipName        string `json:"eipName,omitempty"`
	SubProductType string `json:"subProductType,omitempty"`
	// eip??????
	BandwidthInMbps int `json:"bandwidthInMbps,omitempty"`

	SubnetUuiD      string `json:"subnetUuid,omitempty"`      // ??????uuid
	SecurityGroupID string `json:"securityGroupId,omitempty"` // ?????????id

	AdminPass        string `json:"adminPass,omitempty"`
	AdminPassConfirm string `json:"adminPassConfirm,omitempty"`
	PurchaseLength   int    `json:"purchaseLength,omitempty"`
	// ?????????????????????
	PurchaseNum int `json:"purchaseNum,omitempty"`

	AutoRenewTimeUnit   string                `json:"autoRenewTimeUnit,omitempty"`
	AutoRenewTime       int64                 `json:"autoRenewTime,omitempty"`
	CreateEphemeralList []CreateEphemeralList `json:"createEphemeralList,omitempty"`
	// ?????????????????? ???????????? ??????????????????????????????
	AutoRenew bool `json:"autoRenew,omitempty"`
	// ??????id ??????????????? ?????????ubuntu1604
	ImageID           string `json:"imageId,omitempty"`
	OsName            string `json:"osName,omitempty"`
	SecurityGroupName string `json:"securityGroupName,omitempty"`
	// BCC
	ServiceType string `json:"serviceType,omitempty"`
	GroupID     string `json:"groupID"`
}

// CreateEphemeralList define storage
type CreateEphemeralList struct {
	// ?????????????????? ???????????????????????? ???????????? ?????????ssd
	StorageType string `json:"storageType,omitempty"`
	// ????????????
	SizeInGB int `json:"sizeInGB,omitempty"`
}

// CdsOrderConfig define CDS order config
type CdsOrderConfig struct {
	// ??????????????????????????????postpay
	ProductType string `json:"productType,omitempty"`
	// "zoneA"
	LogicalZone    string `json:"logicalZone,omitempty"`
	Region         string `json:"region,omitempty"`         // "bj"
	PurchaseNum    int    `json:"purchaseNum,omitempty"`    // 1
	PurchaseLength int    `json:"purchaseLength,omitempty"` // 1
	AutoRenewTime  int    `json:"autoRenewTime,omitempty"`  // 0
	// "month"
	AutoRenewTimeUnit string           `json:"autoRenewTimeUnit,omitempty"`
	CdsDiskSize       []DiskSizeConfig `json:"cdsDiskSize,omitempty"`
	// "CDS"
	ServiceType string `json:"serviceType,omitempty"`
}

// EipOrderConfig define CDS order config
type EipOrderConfig struct {
	// ??????????????????????????????postpay
	ProductType     string `json:"productType,omitempty"`
	BandwidthInMbps int    `json:"bandwidthInMbps,omitempty"` // 1000
	Region          string `json:"region,omitempty"`          // "bj"
	SubProductType  string `json:"subProductType,omitempty"`  // "netraffic",
	// EIP???????????????????????????BCC???????????????
	PurchaseNum       int    `json:"purchaseNum,omitempty"`
	PurchaseLength    int    `json:"purchaseLength,omitempty"`    // 1
	AutoRenewTime     int    `json:"autoRenewTime,omitempty"`     // 0
	AutoRenewTimeUnit string `json:"autoRenewTimeUnit,omitempty"` // "month",
	Name              string `json:"name,omitempty"`              // "kkk"
	ServiceType       string `json:"serviceType,omitempty"`       // "EIP"
}

// ScaleDownClusterArgs define  args
type ScaleDownClusterArgs struct {
	ClusterID string     `json:"clusterUuid"`
	AuthCode  string     `json:"authCode"`
	NodeInfos []NodeInfo `json:"nodeInfo"`
}

// NodeInfo defines instanceid
type NodeInfo struct {
	InstanceID string `json:"instanceId"`
}

// ScaleDownClusterResponse defines args
type ScaleDownClusterResponse struct {
	ClusterID string   `json:"clusterUuid"`
	OrderID   []string `json:"orderId"`
}

// ScaleUpClusterWithGroupIDArgs define the args of ScaleUpCluster's request
type ScaleUpClusterWithGroupIDArgs struct {
	GroupID   string `json:"groupId"`
	ClusterID string `json:"clusterId"`
	Num       int    `json:"num"`
}

// ScaleUpCluster scaleup a  cluster
func (c *Client) ScaleUpCluster(args *ScaleUpClusterArgs) (*ScaleUpClusterResponse, error) {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"clientToken": c.GenerateClientToken(),
			"scalingUp":   "",
		}
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	util.Debug("", fmt.Sprintf("ScaleUpCluster Post body: %s", string(postContent)))
	req, err := bce.NewRequest("POST", c.GetURL("v1/cluster", params), bytes.NewBuffer(postContent))
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}
	var scResp *ScaleUpClusterResponse
	err = json.Unmarshal(bodyContent, &scResp)

	if err != nil {
		return nil, err
	}
	return scResp, nil
}

// ScaleDownCluster scale down a  cluster
func (c *Client) ScaleDownCluster(args *ScaleDownClusterArgs) error {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"clientToken": c.GenerateClientToken(),
			"scalingDown": "",
		}
	}

	klog.Infof("ScaleDownCluster args: %v", params)
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/cluster", params), bytes.NewBuffer(postContent))
	klog.Infof("ScaleDownCluster req: %v", req)

	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	return err
}

// DescribeGroup returns the description of the group
func (c *Client) DescribeGroup(groupID string, clusterID string) (*CceGroup, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}
	if groupID == "" {
		return nil, fmt.Errorf("groupID should not be nil")
	}

	params := map[string]string{
		"clusterUuid": clusterID,
		"groupId":     groupID,
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/cluster/group", params), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var cceGroup CceGroup
	err = json.Unmarshal(bodyContent, &cceGroup)

	if err != nil {
		return nil, err
	}
	return &cceGroup, nil
}

// GetAsgNodes returns the group's nodes
func (c *Client) GetAsgNodes(groupID string, clusterID string) ([]CceInstance, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}

	if groupID == "" {
		return nil, fmt.Errorf("groupID should not be nil")
	}

	params := map[string]string{
		"clusterUuid": clusterID,
		"groupId":     groupID,
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/cluster/group/instances", params), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var insList ListInstancesResponse
	err = json.Unmarshal(bodyContent, &insList)

	if err != nil {
		return nil, err
	}

	return insList.Instances, nil
}

// ScaleUpClusterWithGroupID scales up cluster
func (c *Client) ScaleUpClusterWithGroupID(args *ScaleUpClusterWithGroupIDArgs) error {
	if args == nil || args.ClusterID == "" ||
		args.GroupID == "" || args.Num < 0 {
		return fmt.Errorf("ScaleUpClusterWithGroupIDArgs err")
	}

	var params map[string]string
	if args != nil {
		params = map[string]string{
			"groupId":     args.GroupID,
			"clusterUuid": args.ClusterID,
			"num":         strconv.Itoa(args.Num),
		}
	}

	req, err := bce.NewRequest("POST", c.GetURL("v1/cluster/group/scaling_up", params), nil)
	if err != nil {
		return err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return err
	}
	var scResp ScaleUpClusterResponse
	err = json.Unmarshal(bodyContent, &scResp)

	if err != nil {
		return err
	}
	return nil

}
