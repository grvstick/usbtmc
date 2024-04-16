// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

const reservedField = 0x00

const headerSize = 12

type msgIDtype uint8

// The following msgID values are found in Table 2 under the MACRO column of
// the USBTMC Specificiation 1.0, April 14, 2003. The end of line comment shows
// the MACRO names as given in the USBTMC specification.
// The trigger msgID comes from Table 1 -- USB488 defined msgID values of the
// USBTMC-USB488 Specification 1.0, April 14, 2003.
const (
	MsgIdDevDepMsgOut            msgIDtype = 1   // DEV_DEP_MSG_OUT
	MsgIdRequestDevDepMsgIn      msgIDtype = 2   // REQUEST_DEV_DEP_MSG_IN
	MsgIdDevDepMsgIn             msgIDtype = 2   // DEV_DEP_MSG_IN
	MsgIdVendorSpecificOut       msgIDtype = 126 // VENDOR_SPECIFIC_OUT
	MsgIdRequestVendorSpecificIn msgIDtype = 127 // REQUEST_VENDOR_SPECIFIC_IN
	MsgIdVendorSpecificIn        msgIDtype = 127 // VENDOR_SPECIFIC_IN
	MsgIdTrigger                 msgIDtype = 128 // TRIGGER
)

type bRequest uint8

// The USBMTC bRequest constants come from Table 15 -- USBTMC bRequest values in
// the USBTMC Specificiation 1.0, April 14, 2003.
const (
	initiateAbortBulkOut    bRequest = 1  // INITIATE_ABORT_BULK_OUT
	checkAbortBulkOutStatus bRequest = 2  // CHECK_ABORT_BULK_OUT_STATUS
	initiateAbortBulkIn     bRequest = 3  // INITIATE_ABORT_BULK_IN
	checkAbortBulkInStatus  bRequest = 4  // CHECK_ABORT_BULK_IN_STATUS
	initiateClear           bRequest = 5  // INITIATE_CLEAR
	checkClearStatus        bRequest = 6  // CHECK_CLEAR_STATUS
	getCapabilities         bRequest = 7  // GET_CAPABILITIES
	indicatorPulse          bRequest = 64 // INDICATOR_PULSE
)

// The USB488 bRequest constants come from Table 9 -- USB488 defined bRequest
// values in the USBTMC-USB488 Specification 1.0, April 14, 2003
const (
	readStatusByte bRequest = 128 // READ_STATUS_BYTE
	renControl     bRequest = 160 // REN_CONTROL
	goToLocal      bRequest = 161 // GO_TO_LOCAL
	localLockout   bRequest = 162 // LOCAL_LOCKOUT
)

var requestDescription = map[bRequest]string{
	initiateAbortBulkOut:    "Aborts a Bulk-OUT transfer.",
	checkAbortBulkOutStatus: "Returns the status of the previously sent initiateAbortBulkOut request.",
	initiateAbortBulkIn:     "Aborts a Bulk-IN transfer.",
	checkAbortBulkInStatus:  "Returns the status of the previously sent initiateAbortBulkIn request.",
	initiateClear:           "Clears all previously sent pending and unprocessed Bulk-OUT USBTMC message content and clears all pending Bulk-IN transfers from the USBTMC interface.",
	checkClearStatus:        "Returns the status of the previously sent initiateClear request.",
	getCapabilities:         "Returns attributes and capabilities of the USBTMC interface.",
	indicatorPulse:          "A mechanism to turn on an activity indicator for identification purposes. The device indicates whether or not it supports this request in the getCapabilities response packet.",
	readStatusByte:          "Returns the IEEE 488 Status Byte.",
	renControl:              "Mechanism to enable or disable local controls on a device.",
	goToLocal:               "Mechanism to enable local controls on a device.",
	localLockout:            "Mechanism to disable local controls on a device.",
}

func (req bRequest) String() string {
	return requestDescription[req]
}

// type status byte

// // The status constant values come from Table 16 -- USBTMC_status values USBTMC
// // Specificiation 1.0, April 14, 2003, and from Table 10 -- USB488 defined
// // USBTMC_status values in the USBTMC-USB488 Specification 1.0, April 14, 2003
// const (
// 	statusSuccess               status = 0x01 // STATUS_SUCCESS
// 	statusPending               status = 0x02 // STATUS_PENDING
// 	statusInterruptInBusy       status = 0x20 // STATUS_INTERRUPT_IN_BUSY
// 	statusFailed                status = 0x80 // STATUS_FAILED
// 	statusTransferNotInProgress status = 0x81 // STATUS_TRANSFER_NOT_IN_PROGRESS
// 	statusSplitNotInProgress    status = 0x82 // STATUS_SPLIT_NOT_IN_PROGRESS
// 	statusSplitInProgress       status = 0x83 // STATUS_SPLIT_IN_PROGRESS
// )
