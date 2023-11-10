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
	UDS_PATH = "unix:/run/ccnp/uds/measurement.sock"
)

type GetPlatformMeasurementOptions struct {
	measurement_type pb.CATEGORY
	report_data      string
	register_index   int32
}

func checkMeasurementType(measurement_type pb.CATEGORY) bool {
	return measurement_type == pb.CATEGORY_TEE_REPORT || measurement_type == pb.CATEGORY_TDX_RTMR || measurement_type == pb.CATEGORY_TPM
}

func WithMeasurementType(measurement_type pb.CATEGORY) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.measurement_type = measurement_type
	}
}

func WithReportData(report_data string) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.report_data = report_data
	}
}

func WithRegisterIndex(register_index int32) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.register_index = register_index
	}
}

func GetPlatformMeasurement(opts ...func(*GetPlatformMeasurementOptions)) (string, error) {
	//check parameters

	input := GetPlatformMeasurementOptions{measurement_type: pb.CATEGORY_TEE_REPORT, report_data: "", register_index: 0}
	for _, opt := range opts {
		opt(&input)
	}

	if !checkMeasurementType(input.measurement_type) {
		log.Fatalf("[GetPlatformMeasurement] Invalid measurement_type specified")
	}

	if len(report_data) > 64 {
		log.Fatalf("[GetPlatformMeasurement] Invalid report_data specified")
	}

	if register_index < 0 || register_index > 16 {
		log.Fatalf("[GetPlatformMeasurement] Invalid register_index specified")
	}

	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewMeasurementClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetMeasurement(ctx, &pb.GetMeasurementRequest{
		MeasurementType:     pb.TYPE_PAAS,
		MeasurementCategory: input.measurement_type,
		ReportData:          input.report_data,
		RegisterIndex:       input.register_index,
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
