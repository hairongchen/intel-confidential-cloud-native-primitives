/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package eventlog

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/eventlog/proto"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	UDS_PATH = "unix:/run/ccnp/uds/eventlog.sock"
	TYPE_TDX = "TDX"
	TYPE_TPM = "TPM"
)

type CCEventLogEntry struct {
	regIdx  uint32
	evtType uint32
	evtSize uint32
	algId   uint16
	event   []uint8
	digest  []uint8
}

type TDEventLogSpecIdHeader struct {
	Address     uint64
	Length      int
	HeaderData  []byte
	Rtmr        uint32
	Etype       uint32
	DigestCount uint32
	DigestSizes map[uint16]uint16
}

type TDEventLog struct {
	Rtmr        uint32
	Etype       uint32
	DigestCount uint32
	Digests     []string
	Data        []byte
	Event       []byte
	Length      int
	EventSize   uint32
	AlgorithmId uint16
}

type TDEventLogs struct {
	Header    TDEventLogSpecIdHeader
	EventLogs []TDEventLog
}

type GetPlatformEventlogOptions struct {
	eventlogCategory pb.CATEGORY
	startPosition    int32
	count            int32
}

func WithEventlogCategory(eventlogCategory pb.CATEGORY) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.eventlogCategory = eventlogCategory
	}
}

func WithStartPosition(startPosition int32) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.startPosition = startPosition
	}
}

func WithCount(count int32) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.count = count
	}
}

func isEventlogCategoryValid(eventlogCategory pb.CATEGORY) bool {
	return eventlogCategory == pb.CATEGORY_TDX_EVENTLOG || eventlogCategory == pb.CATEGORY_TPM_EVENTLOG
}

func getRawEventlogs(response pb.GetEventlogReply) ([]byte, error) {
	path := response.EventlogDataLoc
	if path == "" {
		log.Fatalf("[getRawEventlogs] Failed to get eventlog from server")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("[getRawEventlogs] Error reading data from  %v: %v", path, err)
	}

	return data, nil
}

func parseTdxEventlog(rawEventlog []byte) ([]CCEventLogEntry, error) {
	var jsonEventlog = TDEventLogs{}
	err := json.Unmarshal(rawEventlog, &jsonEventlog)
	if err != nil {
		log.Fatalf("[parseEventlog] Error unmarshal raw eventlog: %v", err)
	}

	rawEventLogList := jsonEventlog.EventLogs
	var parsedEventLogList []CCEventLogEntry
	for i := 0; i < len(rawEventLogList); i++ {
		rawEventlog := rawEventLogList[i]
		eventLog := CCEventLogEntry{}

		if rawEventlog.DigestCount < 1 {
			continue
		}

		eventLog.regIdx = rawEventlog.Rtmr
		eventLog.evtType = rawEventlog.Etype
		eventLog.evtSize = rawEventlog.EventSize
		eventLog.algId = rawEventlog.AlgorithmId
		eventLog.event = rawEventlog.Event
		eventLog.digest = []uint8(rawEventlog.Digests[rawEventlog.DigestCount-1])
		parsedEventLogList = append(parsedEventLogList, eventLog)

	}

	return parsedEventLogList, nil
}

func getPlatformEventlog(opts ...func(*GetPlatformEventlogOptions)) ([]CCEventLogEntry, error) {

	input := GetPlatformEventlogOptions{eventlogCategory: pb.CATEGORY_TDX_EVENTLOG, startPosition: 0, count: 0}
	for _, opt := range opts {
		opt(&input)
	}

	if !isEventlogCategoryValid(input.eventlogCategory) {
		log.Fatalf("[getPlatformEventlog] Invalid eventlogCategory specified")
	}

	if input.eventlogCategory == pb.CATEGORY_TPM_EVENTLOG {
		log.Fatalf("[getPlatformEventlog] TPM to be supported later")
	}

	if input.startPosition < 0 {
		log.Fatalf("[getPlatformEventlog] Invalid startPosition specified")
	}

	if input.count <= 0 {
		log.Fatalf("[getPlatformEventlog] Invalid count specified")
	}

	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[getPlatformEventlog] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewEventlogClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.getEventlog(ctx, &pb.GetEventlogRequest{
		EventlogLevel:    pb.LEVEL_PAAS,
		EventlogCategory: input.eventlogCategory,
		StartPosition:    input.startPosition,
		Count:            input.count,
	})
	if err != nil {
		log.Fatalf("[getPlatformEventlog] fail to get Platform Eventlog: %v", err)
	}

	switch input.eventlogCategory {
	case pb.CATEGORY_TDX_EVENTLOG:
		rawEventlog, err := getRawEventlogs(response)
		if err != nil {
			log.Fatalf("[getPlatformEventlog] fail to get raw eventlog: %v", err)
		}

		return parseTdxEventlog(rawEventlog)

	case pb.CATEGORY_TPM_EVENTLOG:
		return nil, pkgerrors.New("[getPlatformEventlog] vTPM to be supported later")
	default:
		log.Fatalf("[getPlatformEventlog] unknown TEE enviroment!")
	}

	return nil, nil
}
