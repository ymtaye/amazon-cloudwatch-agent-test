// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package multi_config

import (
	"log"
	"time"

	"github.com/aws/amazon-cloudwatch-agent-test/util/awsservice"
	"github.com/aws/amazon-cloudwatch-agent-test/util/common"
)

const (
	configOutputPath              = "/opt/aws/amazon-cloudwatch-agent/bin/config.json"
	namespace                     = "MultiConfigTest"
	numberofLinuxAppendDimensions = 1
)

var (
	expectedMetrics = []string{"mem_used_percent", "cpu_time_active_userdata", "disk_free"}
)

// Let the agent run for 2 minutes. This will give agent enough time to call server
const agentRuntime = 2 * time.Minute

func Validate() error {

	agentConfigurations := []string{"resources/linux/LinuxCpuOnlyConfig.json", "resources/linux/LinuxMemoryOnlyConfig.json", "resources/linux/LinuxDiskOnlyConfig.json"}

	AppendConfigs(agentConfigurations, configOutputPath)

	time.Sleep(agentRuntime)
	log.Printf("Agent has been running for : %s", agentRuntime.String())
	common.StopAgent()

	// test for cloud watch metrics
	dimensionFilter := awsservice.BuildDimensionFilterList(numberofLinuxAppendDimensions)
	for _, expectedMetric := range expectedMetrics {
		err := awsservice.ValidateMetric(expectedMetric, namespace, dimensionFilter)
		if err != nil {
			log.Printf("CloudWatch Agent append config not working : %v", err)
			return err
		}
	}
	return nil
}
