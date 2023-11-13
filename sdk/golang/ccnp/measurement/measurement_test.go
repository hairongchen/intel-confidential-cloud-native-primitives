package measurement

import (
	"testing"

	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/measurement/proto"
)

const (
	//EMPTY_REPORT_DATA_ENCODED    = "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="
	//EXPECTED_REPORT_DATA_ENCODED = "XUccU3O9poJXiX53jNGj1w2v4WVAw8TKDyWm8Y0xgJ2khEMyCSCiWfO/sYMEn5xoC8ES2VzXwmKRv9NVu3YnUA=="
	EXPECTED_REPORT_DATA       = "abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"
	CATEGORY_UNKNOWN           = 3
	TDX_RTMR_INDEX_UNKNOWN     = 4
	EXPECTED_TDX_REPORT_LEN    = 584
	TEE_TYPE_TDX               = 129
	TDX_TCB_SVN_LENGTH         = 16
	TDX_MRSEAM_LENGTH          = 48
	TDX_MRSEAMSINGER_LENGTH    = 48
	TDX_SEAM_ATTRIBUTES_LENGTH = 8
	TDX_TD_ATTRIBUTES_LENGTH   = 8
	TDX_XFAM_LENGTH            = 8
	TDX_MRTD_LENGTH            = 48
	TDX_MRCONFIGID_LENGTH      = 48
	TDX_MROWNER_LENGTH         = 48
	TDX_MROWNERCONFIG_LENGTH   = 48
	TDX_RTMR_LENGTH            = 48
	TDX_RTMRS_LENGTH           = 192
	TDX_REPORT_DATA_LENGTH     = 64
)

func parseTDXReportAndEvaluate(r TDReportInfo, withUserData bool, t *testing.T) {
	if len(r.TDReportRaw) != EXPECTED_TDX_REPORT_LEN {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport size, retrieved: %v, expected: %v", len(r.TDReportRaw), EXPECTED_TDX_REPORT_LEN)
	}

	tdreport := r.TDReport
	if len(tdreport.TeeTcbSvn) != TDX_TCB_SVN_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport TEE TCB SVN length, retrieved: %v, expected: %v", len(tdreport.TeeTcbSvn), TDX_TCB_SVN_LENGTH)
	}

	if len(tdreport.Mrseam) != TDX_MRSEAM_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrseam length, retrieved: %v, expected: %v", len(tdreport.Mrseam), TDX_MRSEAM_LENGTH)
	}

	if len(tdreport.Mrseamsigner) != TDX_MRSEAMSINGER_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrseamsigner length, retrieved: %v, expected: %v", len(tdreport.Mrseamsigner), TDX_MRSEAMSINGER_LENGTH)
	}

	if len(tdreport.SeamAttributes) != TDX_SEAM_ATTRIBUTES_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport SeamAttributes length, retrieved: %v, expected: %v", len(tdreport.SeamAttributes), TDX_SEAM_ATTRIBUTES_LENGTH)
	}

	if len(tdreport.TdAttributes) != TDX_TD_ATTRIBUTES_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport TdAttributes length, retrieved: %v, expected: %v", len(tdreport.TdAttributes), TDX_TD_ATTRIBUTES_LENGTH)
	}

	if len(tdreport.Xfam) != TDX_XFAM_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Xfam length, retrieved: %v, expected: %v", len(tdreport.Xfam), TDX_XFAM_LENGTH)
	}

	if len(tdreport.Mrtd) != TDX_MRTD_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrtd length, retrieved: %v, expected: %v", len(tdreport.Mrtd), TDX_MRTD_LENGTH)
	}

	if len(tdreport.Mrconfigid) != TDX_MRCONFIGID_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrconfigid length, retrieved: %v, expected: %v", len(tdreport.Mrconfigid), TDX_MRCONFIGID_LENGTH)
	}

	if len(tdreport.Mrowner) != TDX_MROWNER_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrowner length, retrieved: %v, expected: %v", len(tdreport.Mrowner), TDX_MROWNER_LENGTH)
	}

	if len(tdreport.Mrownerconfig) != TDX_MROWNERCONFIG_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrownerconfig length, retrieved: %v, expected: %v", len(tdreport.Mrownerconfig), TDX_MROWNERCONFIG_LENGTH)
	}

	if len(tdreport.Rtmrs) != TDX_RTMRS_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Rtmrs length, retrieved: %v, expected: %v", len(tdreport.Rtmrs), TDX_RTMRS_LENGTH)
	}

	if len(tdreport.ReportData) != TDX_REPORT_DATA_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Data size, retrieved: %v, expected: %v", len(tdreport.ReportData), TDX_REPORT_DATA_LENGTH)
	}

	if withUserData {
		if string(tdreport.ReportData[:]) != EXPECTED_REPORT_DATA {
			t.Fatalf("[parseTDXReportAndEvaluate], report data retrieve = %v, want %v",
				EXPECTED_REPORT_DATA, EXPECTED_REPORT_DATA)
		}
	} else {
		if string(tdreport.ReportData[:]) != "" {
			t.Fatalf("[parseTDXReportAndEvaluate], report data retrieve = %v, want empty string",
				EXPECTED_REPORT_DATA)
		}
	}
}

func parseTDXRtmrAndEvaluate(r TDXRtmrInfo, t *testing.T) {
	if len(r.TDXRtmrRaw) != TDX_RTMR_LENGTH {
		t.Fatalf("[parseTDXRtmrAndEvaluate] wrong RTMT size, retrieved: %v, expected: %v", len(r.TDXRtmrRaw), TDX_RTMR_LENGTH)
	}
}

func TestGetPlatformMeasurementTDReportDefault(t *testing.T) {
	ret, err := GetPlatformMeasurement()
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportDefault] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, false, t)
	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportDefault] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryOnly(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryOnly] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, false, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryOnly] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryAndEmptyReportData(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(""))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, false, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryAndReportData(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(EXPECTED_REPORT_DATA))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, true, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] unknown TEE enviroment!")
	}
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
