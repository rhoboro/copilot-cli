// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cloudwatch provides a client to make API requests to Amazon CloudWatch Service.
package cloudwatch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
)

const (
	resourceQueryType      = "TAG_FILTERS_1_0"
	cloudwatchResourceType = "AWS::CloudWatch::Alarm"
	compositeAlarmType     = "Composite"
	metricAlarmType        = "Metric"
)

type cwClient interface {
	DescribeAlarms(input *cloudwatch.DescribeAlarmsInput) (*cloudwatch.DescribeAlarmsOutput, error)
}

type resourceGroupClient interface {
	SearchResources(input *resourcegroups.SearchResourcesInput) (*resourcegroups.SearchResourcesOutput, error)
}

// CloudWatch wraps an Amazon CloudWatch client.
type CloudWatch struct {
	cwClient
	resourceGroupClient
}

// AlarmStatus contains CloudWatch alarm status.
type AlarmStatus struct {
	Arn          string
	Name         string
	Reason       string
	Status       string
	Type         string
	UpdatedTimes int64
}

// New returns a CloudWatch struct configured against the input session.
func New(s *session.Session) *CloudWatch {
	return &CloudWatch{
		cwClient:            cloudwatch.New(s),
		resourceGroupClient: resourcegroups.New(s),
	}
}

// GetAlarmsWithTags returns all the CloudWatch alarms that have the resource tags.
func (cw *CloudWatch) GetAlarmsWithTags(tags map[string]string) ([]AlarmStatus, error) {
	var alarmNames []*string
	resourceResp := &resourcegroups.SearchResourcesOutput{}
	query, err := cw.searchResourceQuery(tags)
	if err != nil {
		return nil, fmt.Errorf("construct search resource query: %w", err)
	}
	for {
		resourceResp, err = cw.SearchResources(&resourcegroups.SearchResourcesInput{
			NextToken: resourceResp.NextToken,
			ResourceQuery: &resourcegroups.ResourceQuery{
				Type:  aws.String(resourceQueryType),
				Query: aws.String(string(query)),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("search CloudWatch alarm resources: %w", err)
		}
		for _, identifier := range resourceResp.ResourceIdentifiers {
			name, err := cw.getAlarmName(*identifier.ResourceArn)
			if err != nil {
				return nil, err
			}
			alarmNames = append(alarmNames, name)
		}
		if resourceResp.NextToken == nil {
			break
		}
	}
	var alarmStatus []AlarmStatus
	alarmResp := &cloudwatch.DescribeAlarmsOutput{}
	for {
		alarmResp, err = cw.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
			AlarmNames: alarmNames,
			NextToken:  alarmResp.NextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe CloudWatch alarms: %w", err)
		}
		alarmStatus = append(alarmStatus, cw.compositeAlarmsStatus(alarmResp.CompositeAlarms)...)
		alarmStatus = append(alarmStatus, cw.metricAlarmsStatus(alarmResp.MetricAlarms)...)
		if alarmResp.NextToken == nil {
			break
		}
	}
	return alarmStatus, nil
}

func (cw *CloudWatch) searchResourceQuery(tags map[string]string) ([]byte, error) {
	type keyVal struct {
		Key    string
		Values []string
	}
	type query struct {
		ResourceTypeFilters []string
		TagFilters          []keyVal
	}
	var keyVals []keyVal
	for k, v := range tags {
		keyVals = append(keyVals, keyVal{
			Key:    k,
			Values: []string{v},
		})
	}
	queryStruct := query{
		ResourceTypeFilters: []string{cloudwatchResourceType},
		TagFilters:          keyVals,
	}
	return json.Marshal(queryStruct)
}

// getAlarmName gets the alarm name given a specific alarm ARN.
// For example: arn:aws:cloudwatch:us-west-2:1234567890:alarm:SDc-ReadCapacityUnitsLimit-BasicAlarm
// returns SDc-ReadCapacityUnitsLimit-BasicAlarm
func (cw *CloudWatch) getAlarmName(alarmArn string) (*string, error) {
	resp, err := arn.Parse(alarmArn)
	if err != nil {
		return nil, fmt.Errorf("parse alarm ARN %s: %w", alarmArn, err)
	}
	alarmNameList := strings.Split(resp.Resource, ":")
	if len(alarmNameList) != 2 {
		return nil, fmt.Errorf("cannot parse alarm ARN resource %s", resp.Resource)
	}
	return aws.String(alarmNameList[1]), nil
}

func (cw *CloudWatch) compositeAlarmsStatus(alarms []*cloudwatch.CompositeAlarm) []AlarmStatus {
	var alarmStatusList []AlarmStatus
	for _, alarm := range alarms {
		alarmStatusList = append(alarmStatusList, AlarmStatus{
			Arn:          aws.StringValue(alarm.AlarmArn),
			Name:         aws.StringValue(alarm.AlarmName),
			Reason:       aws.StringValue(alarm.StateReason),
			Status:       aws.StringValue(alarm.StateValue),
			Type:         compositeAlarmType,
			UpdatedTimes: alarm.StateUpdatedTimestamp.Unix(),
		})
	}
	return alarmStatusList
}

func (cw *CloudWatch) metricAlarmsStatus(alarms []*cloudwatch.MetricAlarm) []AlarmStatus {
	var alarmStatusList []AlarmStatus
	for _, alarm := range alarms {
		alarmStatusList = append(alarmStatusList, AlarmStatus{
			Arn:          aws.StringValue(alarm.AlarmArn),
			Name:         aws.StringValue(alarm.AlarmName),
			Reason:       aws.StringValue(alarm.StateReason),
			Status:       aws.StringValue(alarm.StateValue),
			Type:         metricAlarmType,
			UpdatedTimes: alarm.StateUpdatedTimestamp.Unix(),
		})
	}
	return alarmStatusList
}