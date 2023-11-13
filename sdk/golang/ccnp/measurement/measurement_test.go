package measurement

import (
	"encoding/base64"
	"testing"
)

const (
	VALID_REPORT_ENCODED         = ""
	EXPECTED_REPORT_DATA_ENCODED = "XUccU3O9poJXiX53jNGj1w2v4WVAw8TKDyWm8Y0xgJ2khEMyCSCiWfO/sYMEn5xoC8ES2VzXwmKRv9NVu3YnUA=="
	CATEGORY_UNKNOWN             = 3
	TDX_RTMR_INDEX_UNKNOWN       = 4
	EXPECTED_TDX_REPORT_LEN      = 584
	TEE_TYPE_TDX                 = 129
	TDX_TCB_SVN_LENGTH           = 16
	TDX_MRSEAM_LENGTH            = 48
	TDX_MRSINGERSEAM_LENGTH      = 48
	TDX_SEAM_ATTRIBUTES_LENGTH   = 8
	TDX_TD_ATTRIBUTES_LENGTH     = 8
	TDX_XFAM_LENGTH              = 8
	TDX_MRTD_LENGTH              = 48
	TDX_MRCONFIGID_LENGTH        = 48
	TDX_MROWNER_LENGTH           = 48
	TDX_MROWNERCONFIG_LENGTH     = 48
	TDX_RTMR_LENGTH              = 192
	TDX_REPORT_DATA_LENGTH       = 64
)

func parseTDXReportAndEvaluate(report []byte) nil {
	tdreport := parseTDXReport(report)
	if len(tdreport.ReportData) != TDX_REPORT_DATA_LENGTH {
		t.Fatalf("[TestGetPlatformMeasurement] wrong TDReport size, retrieved: %v, expected: %v", len(tdreport.ReportData), TDX_REPORT_DATA_LENGTH)
	}
}

func TestGetPlatformMeasurementTDReport(t *testing.T) {
	reportData, errDecode := base64.StdEncoding.DecodeString(EXPECTED_REPORT_DATA_ENCODED)
	if errDecode != nil {
		t.Fatalf("[TestGetPlatformMeasurement] decode report data error: %v", errDecode)
	}

	//test get TEE Report
	ret, err := GetPlatformMeasurement()
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurement] get Platform Measurement error: %v", err)
	}
	//TODO: now only TDX report is supported, if other TEEs added, need to check TEE type first
	measurement := ret.Measurement
	if len(measurement) != EXPECTED_TDX_REPORT_LEN {
		t.Fatalf("[TestGetPlatformMeasurement] wrong TDReport size, retrieved: %v, expected: %v", len(measurement), EXPECTED_TDX_REPORT_LEN)
	}

	parseTDXReportAndEvaluate(measurement)
	/*
	   ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT))
	   ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(""))
	   ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(reportData))

	   //test call with undefined report category
	   ret, err := GetPlatformMeasurement(WithMeasurementType(CATEGORY_UNKNOWN))

	   //test call with undefined rtmr index
	   ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR), WithRegisterIndex(TDX_RTMR_INDEX_UNKNOWN))
	*/
}

/*
func TestGetPlatformMeasurementRTMR(t *testing.T) {
	//test get TDX RTMR
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR))
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR), WithRegisterIndex(1))

	//test call with undefined report category
	ret, err := GetPlatformMeasurement(WithMeasurementType(CATEGORY_UNKNOWN))

	//test call with undefined rtmr index
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR), WithRegisterIndex(TDX_RTMR_INDEX_UNKNOWN))

}

func TestGetPlatformMeasurementWrongParameters(t *testing.T) {
	//test call with undefined report category
	ret, err := GetPlatformMeasurement(WithMeasurementType(CATEGORY_UNKNOWN))

	//test call with undefined rtmr index
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR), WithRegisterIndex(TDX_RTMR_INDEX_UNKNOWN))

}
*/
