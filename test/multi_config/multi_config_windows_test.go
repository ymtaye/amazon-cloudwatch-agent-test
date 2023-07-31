// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build windows
// +build windows

package multi_config

import (
	"log"
	"time"

	"github.com/aws/amazon-cloudwatch-agent-test/util/awsservice"
	"github.com/aws/amazon-cloudwatch-agent-test/util/common"
)

const (
	configOutputPath                = "C:\\ProgramData\\Amazon\\AmazonCloudWatchAgent\\config.json"
	namespace                       = "MultiConfigWindowsTest"
	agentRuntime                    = 2 * time.Minute
	numberofWindowsAppendDimensions = 1
)

var (
	expectedMetrics = []string{"% Committed Bytes In Use", "% InterruptTime", "% Disk Time"}
)

func Validate() error {
	agentConfigurations := []string{"resources/windows/WindowsCompleteConfig.json", "resources/windows/WindowsMemoryOnlyConfig.json"}

	AppendConfigs(agentConfigurations, configOutputPath)

	time.Sleep(agentRuntime)
	log.Printf("Agent has been running for : %s", agentRuntime.String())
	err := common.StopAgent()
	if err != nil {
		log.Printf("Stopping agent failed: %v", err)
		return err
	}

	dimensionFilter := awsservice.BuildDimensionFilterList(numberofWindowsAppendDimensions)
	for _, expectedMetric := range expectedMetrics {
		err = awsservice.ValidateMetric(expectedMetric, namespace, dimensionFilter)
	}
	if err != nil {
		log.Printf("CloudWatch Agent append config not working : %v", err)
		return err
	}
	return nil
}
