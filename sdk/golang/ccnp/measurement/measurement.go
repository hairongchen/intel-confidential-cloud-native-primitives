/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package measurement

import (
	"context"
	"log"
	"time"

	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/measurement/proto"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	UDS_PATH        = "unix:/run/ccnp/uds/measurement.sock"
	TYPE_TDX        = "TDX"
	TYPE_TPM        = "TPM"
	TYPE_TEE_REPORT = pb.CATEGORY_TEE_REPORT //Get TEE report
	TYPE_TDX_RTMR   = pb.CATEGORY_TDX_RTMR   //Get TDX RTMR measurement (of a specific register)
	TYPE_TPM_PCR    = pb.CATEGORY_TPM        //Get TPM PCR measurement (of a specific register)
)

func GetPlatformMeasurement(measurement_type pb.CATEGORY, report_data string, register_index int) (string, error) {
	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewMeasurementClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if measurement_type > 2 || measurement_type < 0 {
		log.Fatalf("[GetPlatformMeasurement] Invalid measurement type specified")
	}

	if report_data != "" && len(report_data) > 64 {
		log.Fatalf("[GetPlatformMeasurement] Invalid report data specified")
	}

	if register_index < 0 || register_index > 16 {
		log.Fatalf("[GetPlatformMeasurement] Invalid report data specified")
	}

	response, err := client.GetMeasurement(ctx, &pb.GetMeasurementRequest{
		MeasurementType:     pb.TYPE_PAAS,
		MeasurementCategory: measurement_type,
		ReportData:          report_data,
		RegisterIndex:       register_index,
	})
	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] fail to get Platform Measurement: %v", err)
	}

	return response.Measurement, nil
}

func GetContainerMeasurement() (interface{}, error) {
	// TODO: add Container Measurement support later
	return nil, pkgerrors.New("Container Measurement support to be implemented later.")
}
