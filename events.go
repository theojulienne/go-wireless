package wireless

import (
	"encoding/csv"
	"log"
	"strings"
)

// This file contains components from github.com/brlbil/wpaclient
//
// Copyright (c) 2017 Birol Bilgin
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// nd/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// UT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

const (
	// CtrlReq as defined in wpactrl/wpa_ctrl.h:19
	CtrlReq = "CTRL-REQ-"
	// CtrlRsp as defined in wpactrl/wpa_ctrl.h:22
	CtrlRsp = "CTRL-RSP-"
	// EventConnected as defined in wpactrl/wpa_ctrl.h:26
	EventConnected = "CTRL-EVENT-CONNECTED"
	// EventDisconnected as defined in wpactrl/wpa_ctrl.h:28
	EventDisconnected = "CTRL-EVENT-DISCONNECTED"
	// EventAssocReject as defined in wpactrl/wpa_ctrl.h:30
	EventAssocReject = "CTRL-EVENT-ASSOC-REJECT"
	// EventAuthReject as defined in wpactrl/wpa_ctrl.h:32
	EventAuthReject = "CTRL-EVENT-AUTH-REJECT"
	// EventTerminating as defined in wpactrl/wpa_ctrl.h:34
	EventTerminating = "CTRL-EVENT-TERMINATING"
	// EventPasswordChanged as defined in wpactrl/wpa_ctrl.h:36
	EventPasswordChanged = "CTRL-EVENT-PASSWORD-CHANGED"
	// EventEapNotification as defined in wpactrl/wpa_ctrl.h:38
	EventEapNotification = "CTRL-EVENT-EAP-NOTIFICATION"
	// EventEapStarted as defined in wpactrl/wpa_ctrl.h:40
	EventEapStarted = "CTRL-EVENT-EAP-STARTED"
	// EventEapProposedMethod as defined in wpactrl/wpa_ctrl.h:42
	EventEapProposedMethod = "CTRL-EVENT-EAP-PROPOSED-METHOD"
	// EventEapMethod as defined in wpactrl/wpa_ctrl.h:44
	EventEapMethod = "CTRL-EVENT-EAP-METHOD"
	// EventEapPeerCert as defined in wpactrl/wpa_ctrl.h:46
	EventEapPeerCert = "CTRL-EVENT-EAP-PEER-CERT"
	// EventEapPeerAlt as defined in wpactrl/wpa_ctrl.h:48
	EventEapPeerAlt = "CTRL-EVENT-EAP-PEER-ALT"
	// EventEapTLSCertError as defined in wpactrl/wpa_ctrl.h:50
	EventEapTLSCertError = "CTRL-EVENT-EAP-TLS-CERT-ERROR"
	// EventEapStatus as defined in wpactrl/wpa_ctrl.h:52
	EventEapStatus = "CTRL-EVENT-EAP-STATUS"
	// EventEapRetransmit as defined in wpactrl/wpa_ctrl.h:54
	EventEapRetransmit = "CTRL-EVENT-EAP-RETRANSMIT"
	// EventEapRetransmit2 as defined in wpactrl/wpa_ctrl.h:55
	EventEapRetransmit2 = "CTRL-EVENT-EAP-RETRANSMIT2"
	// EventEapSuccess as defined in wpactrl/wpa_ctrl.h:57
	EventEapSuccess = "CTRL-EVENT-EAP-SUCCESS"
	// EventEapSuccess2 as defined in wpactrl/wpa_ctrl.h:58
	EventEapSuccess2 = "CTRL-EVENT-EAP-SUCCESS2"
	// EventEapFailure as defined in wpactrl/wpa_ctrl.h:60
	EventEapFailure = "CTRL-EVENT-EAP-FAILURE"
	// EventEapFailure2 as defined in wpactrl/wpa_ctrl.h:61
	EventEapFailure2 = "CTRL-EVENT-EAP-FAILURE2"
	// EventEapTimeoutFailure as defined in wpactrl/wpa_ctrl.h:63
	EventEapTimeoutFailure = "CTRL-EVENT-EAP-TIMEOUT-FAILURE"
	// EventEapTimeoutFailure2 as defined in wpactrl/wpa_ctrl.h:64
	EventEapTimeoutFailure2 = "CTRL-EVENT-EAP-TIMEOUT-FAILURE2"
	// EventTempDisabled as defined in wpactrl/wpa_ctrl.h:66
	EventTempDisabled = "CTRL-EVENT-SSID-TEMP-DISABLED"
	// EventReenabled as defined in wpactrl/wpa_ctrl.h:68
	EventReenabled = "CTRL-EVENT-SSID-REENABLED"
	// EventScanStarted as defined in wpactrl/wpa_ctrl.h:70
	EventScanStarted = "CTRL-EVENT-SCAN-STARTED"
	// EventScanResults as defined in wpactrl/wpa_ctrl.h:72
	EventScanResults = "CTRL-EVENT-SCAN-RESULTS"
	// EventScanFailed as defined in wpactrl/wpa_ctrl.h:74
	EventScanFailed = "CTRL-EVENT-SCAN-FAILED"
	// EventStateChange as defined in wpactrl/wpa_ctrl.h:76
	EventStateChange = "CTRL-EVENT-STATE-CHANGE"
	// EventBssAdded as defined in wpactrl/wpa_ctrl.h:78
	EventBssAdded = "CTRL-EVENT-BSS-ADDED"
	// EventBssRemoved as defined in wpactrl/wpa_ctrl.h:80
	EventBssRemoved = "CTRL-EVENT-BSS-REMOVED"
	// EventNetworkNotFound as defined in wpactrl/wpa_ctrl.h:82
	EventNetworkNotFound = "CTRL-EVENT-NETWORK-NOT-FOUND"
	// EventSignalChange as defined in wpactrl/wpa_ctrl.h:84
	EventSignalChange = "CTRL-EVENT-SIGNAL-CHANGE"
	// EventBeaconLoss as defined in wpactrl/wpa_ctrl.h:86
	EventBeaconLoss = "CTRL-EVENT-BEACON-LOSS"
	// EventRegdomChange as defined in wpactrl/wpa_ctrl.h:88
	EventRegdomChange = "CTRL-EVENT-REGDOM-CHANGE"
	// EventChannelSwitch as defined in wpactrl/wpa_ctrl.h:90
	EventChannelSwitch = "CTRL-EVENT-CHANNEL-SWITCH"
	// EventSubnetStatusUpdate as defined in wpactrl/wpa_ctrl.h:103
	EventSubnetStatusUpdate = "CTRL-EVENT-SUBNET-STATUS-UPDATE"
	// IbssRsnCompleted as defined in wpactrl/wpa_ctrl.h:106
	IbssRsnCompleted = "IBSS-RSN-COMPLETED"
	// EventFreqConflict as defined in wpactrl/wpa_ctrl.h:113
	EventFreqConflict = "CTRL-EVENT-FREQ-CONFLICT"
	// EventAvoidFreq as defined in wpactrl/wpa_ctrl.h:115
	EventAvoidFreq = "CTRL-EVENT-AVOID-FREQ"
	// WpsEventOverlap as defined in wpactrl/wpa_ctrl.h:117
	WpsEventOverlap = "WPS-OVERLAP-DETECTED"
	// WpsEventApAvailablePbc as defined in wpactrl/wpa_ctrl.h:119
	WpsEventApAvailablePbc = "WPS-AP-AVAILABLE-PBC"
	// WpsEventApAvailableAuth as defined in wpactrl/wpa_ctrl.h:121
	WpsEventApAvailableAuth = "WPS-AP-AVAILABLE-AUTH"
	// WpsEventApAvailablePin as defined in wpactrl/wpa_ctrl.h:124
	WpsEventApAvailablePin = "WPS-AP-AVAILABLE-PIN"
	// WpsEventApAvailable as defined in wpactrl/wpa_ctrl.h:126
	WpsEventApAvailable = "WPS-AP-AVAILABLE"
	// WpsEventCredReceived as defined in wpactrl/wpa_ctrl.h:128
	WpsEventCredReceived = "WPS-CRED-RECEIVED"
	// WpsEventM2d as defined in wpactrl/wpa_ctrl.h:130
	WpsEventM2d = "WPS-M2D"
	// WpsEventFail as defined in wpactrl/wpa_ctrl.h:132
	WpsEventFail = "WPS-FAIL"
	// WpsEventSuccess as defined in wpactrl/wpa_ctrl.h:134
	WpsEventSuccess = "WPS-SUCCESS"
	// WpsEventTimeout as defined in wpactrl/wpa_ctrl.h:136
	WpsEventTimeout = "WPS-TIMEOUT"
	// WpsEventActive as defined in wpactrl/wpa_ctrl.h:138
	WpsEventActive = "WPS-PBC-ACTIVE"
	// WpsEventDisable as defined in wpactrl/wpa_ctrl.h:140
	WpsEventDisable = "WPS-PBC-DISABLE"
	// WpsEventEnrolleeSeen as defined in wpactrl/wpa_ctrl.h:142
	WpsEventEnrolleeSeen = "WPS-ENROLLEE-SEEN"
	// WpsEventOpenNetwork as defined in wpactrl/wpa_ctrl.h:144
	WpsEventOpenNetwork = "WPS-OPEN-NETWORK"
	// WpsEventErApAdd as defined in wpactrl/wpa_ctrl.h:147
	WpsEventErApAdd = "WPS-ER-AP-ADD"
	// WpsEventErApRemove as defined in wpactrl/wpa_ctrl.h:148
	WpsEventErApRemove = "WPS-ER-AP-REMOVE"
	// WpsEventErEnrolleeAdd as defined in wpactrl/wpa_ctrl.h:149
	WpsEventErEnrolleeAdd = "WPS-ER-ENROLLEE-ADD"
	// WpsEventErEnrolleeRemove as defined in wpactrl/wpa_ctrl.h:150
	WpsEventErEnrolleeRemove = "WPS-ER-ENROLLEE-REMOVE"
	// WpsEventErApSettings as defined in wpactrl/wpa_ctrl.h:151
	WpsEventErApSettings = "WPS-ER-AP-SETTINGS"
	// WpsEventErSetSelReg as defined in wpactrl/wpa_ctrl.h:152
	WpsEventErSetSelReg = "WPS-ER-AP-SET-SEL-REG"
	// DppEventAuthSuccess as defined in wpactrl/wpa_ctrl.h:155
	DppEventAuthSuccess = "DPP-AUTH-SUCCESS"
	// DppEventNotCompatible as defined in wpactrl/wpa_ctrl.h:156
	DppEventNotCompatible = "DPP-NOT-COMPATIBLE"
	// DppEventResponsePending as defined in wpactrl/wpa_ctrl.h:157
	DppEventResponsePending = "DPP-RESPONSE-PENDING"
	// DppEventScanPeerQrCode as defined in wpactrl/wpa_ctrl.h:158
	DppEventScanPeerQrCode = "DPP-SCAN-PEER-QR-CODE"
	// DppEventConfReceived as defined in wpactrl/wpa_ctrl.h:159
	DppEventConfReceived = "DPP-CONF-RECEIVED"
	// DppEventConfSent as defined in wpactrl/wpa_ctrl.h:160
	DppEventConfSent = "DPP-CONF-SENT"
	// DppEventConfFailed as defined in wpactrl/wpa_ctrl.h:161
	DppEventConfFailed = "DPP-CONF-FAILED"
	// DppEventConfobjSsid as defined in wpactrl/wpa_ctrl.h:162
	DppEventConfobjSsid = "DPP-CONFOBJ-SSID"
	// DppEventConfobjPass as defined in wpactrl/wpa_ctrl.h:163
	DppEventConfobjPass = "DPP-CONFOBJ-PASS"
	// DppEventConfobjPsk as defined in wpactrl/wpa_ctrl.h:164
	DppEventConfobjPsk = "DPP-CONFOBJ-PSK"
	// DppEventConnector as defined in wpactrl/wpa_ctrl.h:165
	DppEventConnector = "DPP-CONNECTOR"
	// DppEventCSignKey as defined in wpactrl/wpa_ctrl.h:166
	DppEventCSignKey = "DPP-C-SIGN-KEY"
	// DppEventNetAccessKey as defined in wpactrl/wpa_ctrl.h:167
	DppEventNetAccessKey = "DPP-NET-ACCESS-KEY"
	// DppEventMissingConnector as defined in wpactrl/wpa_ctrl.h:168
	DppEventMissingConnector = "DPP-MISSING-CONNECTOR"
	// DppEventNetworkID as defined in wpactrl/wpa_ctrl.h:169
	DppEventNetworkID = "DPP-NETWORK-ID"
	// DppEventRx as defined in wpactrl/wpa_ctrl.h:170
	DppEventRx = "DPP-RX"
	// DppEventTx as defined in wpactrl/wpa_ctrl.h:171
	DppEventTx = "DPP-TX"
	// DppEventTxStatus as defined in wpactrl/wpa_ctrl.h:172
	DppEventTxStatus = "DPP-TX-STATUS"
	// DppEventFail as defined in wpactrl/wpa_ctrl.h:173
	DppEventFail = "DPP-FAIL"
	// MeshGroupStarted as defined in wpactrl/wpa_ctrl.h:176
	MeshGroupStarted = "MESH-GROUP-STARTED"
	// MeshGroupRemoved as defined in wpactrl/wpa_ctrl.h:177
	MeshGroupRemoved = "MESH-GROUP-REMOVED"
	// MeshPeerConnected as defined in wpactrl/wpa_ctrl.h:178
	MeshPeerConnected = "MESH-PEER-CONNECTED"
	// MeshPeerDisconnected as defined in wpactrl/wpa_ctrl.h:179
	MeshPeerDisconnected = "MESH-PEER-DISCONNECTED"
	// MeshSaeAuthFailure as defined in wpactrl/wpa_ctrl.h:181
	MeshSaeAuthFailure = "MESH-SAE-AUTH-FAILURE"
	// MeshSaeAuthBlocked as defined in wpactrl/wpa_ctrl.h:182
	MeshSaeAuthBlocked = "MESH-SAE-AUTH-BLOCKED"
	// P2pEventDeviceFound as defined in wpactrl/wpa_ctrl.h:190
	P2pEventDeviceFound = "P2P-DEVICE-FOUND"
	// P2pEventDeviceLost as defined in wpactrl/wpa_ctrl.h:193
	P2pEventDeviceLost = "P2P-DEVICE-LOST"
	// P2pEventGoNegRequest as defined in wpactrl/wpa_ctrl.h:197
	P2pEventGoNegRequest = "P2P-GO-NEG-REQUEST"
	// P2pEventGoNegSuccess as defined in wpactrl/wpa_ctrl.h:198
	P2pEventGoNegSuccess = "P2P-GO-NEG-SUCCESS"
	// P2pEventGoNegFailure as defined in wpactrl/wpa_ctrl.h:199
	P2pEventGoNegFailure = "P2P-GO-NEG-FAILURE"
	// P2pEventGroupFormationSuccess as defined in wpactrl/wpa_ctrl.h:200
	P2pEventGroupFormationSuccess = "P2P-GROUP-FORMATION-SUCCESS"
	// P2pEventGroupFormationFailure as defined in wpactrl/wpa_ctrl.h:201
	P2pEventGroupFormationFailure = "P2P-GROUP-FORMATION-FAILURE"
	// P2pEventGroupStarted as defined in wpactrl/wpa_ctrl.h:202
	P2pEventGroupStarted = "P2P-GROUP-STARTED"
	// P2pEventGroupRemoved as defined in wpactrl/wpa_ctrl.h:203
	P2pEventGroupRemoved = "P2P-GROUP-REMOVED"
	// P2pEventCrossConnectEnable as defined in wpactrl/wpa_ctrl.h:204
	P2pEventCrossConnectEnable = "P2P-CROSS-CONNECT-ENABLE"
	// P2pEventCrossConnectDisable as defined in wpactrl/wpa_ctrl.h:205
	P2pEventCrossConnectDisable = "P2P-CROSS-CONNECT-DISABLE"
	// P2pEventProvDiscShowPin as defined in wpactrl/wpa_ctrl.h:207
	P2pEventProvDiscShowPin = "P2P-PROV-DISC-SHOW-PIN"
	// P2pEventProvDiscEnterPin as defined in wpactrl/wpa_ctrl.h:209
	P2pEventProvDiscEnterPin = "P2P-PROV-DISC-ENTER-PIN"
	// P2pEventProvDiscPbcReq as defined in wpactrl/wpa_ctrl.h:211
	P2pEventProvDiscPbcReq = "P2P-PROV-DISC-PBC-REQ"
	// P2pEventProvDiscPbcResp as defined in wpactrl/wpa_ctrl.h:213
	P2pEventProvDiscPbcResp = "P2P-PROV-DISC-PBC-RESP"
	// P2pEventProvDiscFailure as defined in wpactrl/wpa_ctrl.h:215
	P2pEventProvDiscFailure = "P2P-PROV-DISC-FAILURE"
	// P2pEventServDiscReq as defined in wpactrl/wpa_ctrl.h:217
	P2pEventServDiscReq = "P2P-SERV-DISC-REQ"
	// P2pEventServDiscResp as defined in wpactrl/wpa_ctrl.h:219
	P2pEventServDiscResp = "P2P-SERV-DISC-RESP"
	// P2pEventServAspResp as defined in wpactrl/wpa_ctrl.h:220
	P2pEventServAspResp = "P2P-SERV-ASP-RESP"
	// P2pEventInvitationReceived as defined in wpactrl/wpa_ctrl.h:221
	P2pEventInvitationReceived = "P2P-INVITATION-RECEIVED"
	// P2pEventInvitationResult as defined in wpactrl/wpa_ctrl.h:222
	P2pEventInvitationResult = "P2P-INVITATION-RESULT"
	// P2pEventInvitationAccepted as defined in wpactrl/wpa_ctrl.h:223
	P2pEventInvitationAccepted = "P2P-INVITATION-ACCEPTED"
	// P2pEventFindStopped as defined in wpactrl/wpa_ctrl.h:224
	P2pEventFindStopped = "P2P-FIND-STOPPED"
	// P2pEventPersistentPskFail as defined in wpactrl/wpa_ctrl.h:225
	P2pEventPersistentPskFail = "P2P-PERSISTENT-PSK-FAIL id="
	// P2pEventPresenceResponse as defined in wpactrl/wpa_ctrl.h:226
	P2pEventPresenceResponse = "P2P-PRESENCE-RESPONSE"
	// P2pEventNfcBothGo as defined in wpactrl/wpa_ctrl.h:227
	P2pEventNfcBothGo = "P2P-NFC-BOTH-GO"
	// P2pEventNfcPeerClient as defined in wpactrl/wpa_ctrl.h:228
	P2pEventNfcPeerClient = "P2P-NFC-PEER-CLIENT"
	// P2pEventNfcWhileClient as defined in wpactrl/wpa_ctrl.h:229
	P2pEventNfcWhileClient = "P2P-NFC-WHILE-CLIENT"
	// P2pEventFallbackToGoNeg as defined in wpactrl/wpa_ctrl.h:230
	P2pEventFallbackToGoNeg = "P2P-FALLBACK-TO-GO-NEG"
	// P2pEventFallbackToGoNegEnabled as defined in wpactrl/wpa_ctrl.h:231
	P2pEventFallbackToGoNegEnabled = "P2P-FALLBACK-TO-GO-NEG-ENABLED"
	// EssDisassocImminent as defined in wpactrl/wpa_ctrl.h:234
	EssDisassocImminent = "ESS-DISASSOC-IMMINENT"
	// P2pEventRemoveAndReformGroup as defined in wpactrl/wpa_ctrl.h:235
	P2pEventRemoveAndReformGroup = "P2P-REMOVE-AND-REFORM-GROUP"
	// P2pEventP2psProvisionStart as defined in wpactrl/wpa_ctrl.h:237
	P2pEventP2psProvisionStart = "P2PS-PROV-START"
	// P2pEventP2psProvisionDone as defined in wpactrl/wpa_ctrl.h:238
	P2pEventP2psProvisionDone = "P2PS-PROV-DONE"
	// InterworkingAp as defined in wpactrl/wpa_ctrl.h:240
	InterworkingAp = "INTERWORKING-AP"
	// InterworkingBlacklisted as defined in wpactrl/wpa_ctrl.h:241
	InterworkingBlacklisted = "INTERWORKING-BLACKLISTED"
	// InterworkingNoMatch as defined in wpactrl/wpa_ctrl.h:242
	InterworkingNoMatch = "INTERWORKING-NO-MATCH"
	// InterworkingAlreadyConnected as defined in wpactrl/wpa_ctrl.h:243
	InterworkingAlreadyConnected = "INTERWORKING-ALREADY-CONNECTED"
	// InterworkingSelected as defined in wpactrl/wpa_ctrl.h:244
	InterworkingSelected = "INTERWORKING-SELECTED"
	// CredAdded as defined in wpactrl/wpa_ctrl.h:247
	CredAdded = "CRED-ADDED"
	// CredModified as defined in wpactrl/wpa_ctrl.h:249
	CredModified = "CRED-MODIFIED"
	// CredRemoved as defined in wpactrl/wpa_ctrl.h:251
	CredRemoved = "CRED-REMOVED"
	// GasResponseInfo as defined in wpactrl/wpa_ctrl.h:253
	GasResponseInfo = "GAS-RESPONSE-INFO"
	// GasQueryStart as defined in wpactrl/wpa_ctrl.h:255
	GasQueryStart = "GAS-QUERY-START"
	// GasQueryDone as defined in wpactrl/wpa_ctrl.h:257
	GasQueryDone = "GAS-QUERY-DONE"
	// RxAnqp as defined in wpactrl/wpa_ctrl.h:262
	RxAnqp = "RX-ANQP"
	// RxHs20Anqp as defined in wpactrl/wpa_ctrl.h:263
	RxHs20Anqp = "RX-HS20-ANQP"
	// RxHs20AnqpIcon as defined in wpactrl/wpa_ctrl.h:264
	RxHs20AnqpIcon = "RX-HS20-ANQP-ICON"
	// RxHs20Icon as defined in wpactrl/wpa_ctrl.h:265
	RxHs20Icon = "RX-HS20-ICON"
	// RxMboAnqp as defined in wpactrl/wpa_ctrl.h:266
	RxMboAnqp = "RX-MBO-ANQP"
	// Hs20SubscriptionRemediation as defined in wpactrl/wpa_ctrl.h:268
	Hs20SubscriptionRemediation = "HS20-SUBSCRIPTION-REMEDIATION"
	// Hs20DeauthImminentNotice as defined in wpactrl/wpa_ctrl.h:269
	Hs20DeauthImminentNotice = "HS20-DEAUTH-IMMINENT-NOTICE"
	// RrmEventNeighborRepRxed as defined in wpactrl/wpa_ctrl.h:274
	RrmEventNeighborRepRxed = "RRM-NEIGHBOR-REP-RECEIVED"
	// RrmEventNeighborRepFailed as defined in wpactrl/wpa_ctrl.h:275
	RrmEventNeighborRepFailed = "RRM-NEIGHBOR-REP-REQUEST-FAILED"
	// WpsEventPinNeeded as defined in wpactrl/wpa_ctrl.h:278
	WpsEventPinNeeded = "WPS-PIN-NEEDED"
	// WpsEventNewApSettings as defined in wpactrl/wpa_ctrl.h:279
	WpsEventNewApSettings = "WPS-NEW-AP-SETTINGS"
	// WpsEventRegSuccess as defined in wpactrl/wpa_ctrl.h:280
	WpsEventRegSuccess = "WPS-REG-SUCCESS"
	// WpsEventApSetupLocked as defined in wpactrl/wpa_ctrl.h:281
	WpsEventApSetupLocked = "WPS-AP-SETUP-LOCKED"
	// WpsEventApSetupUnlocked as defined in wpactrl/wpa_ctrl.h:282
	WpsEventApSetupUnlocked = "WPS-AP-SETUP-UNLOCKED"
	// WpsEventApPinEnabled as defined in wpactrl/wpa_ctrl.h:283
	WpsEventApPinEnabled = "WPS-AP-PIN-ENABLED"
	// WpsEventApPinDisabled as defined in wpactrl/wpa_ctrl.h:284
	WpsEventApPinDisabled = "WPS-AP-PIN-DISABLED"
	// ApStaConnected as defined in wpactrl/wpa_ctrl.h:285
	ApStaConnected = "AP-STA-CONNECTED"
	// ApStaDisconnected as defined in wpactrl/wpa_ctrl.h:286
	ApStaDisconnected = "AP-STA-DISCONNECTED"
	// ApStaPossiblePskMismatch as defined in wpactrl/wpa_ctrl.h:287
	ApStaPossiblePskMismatch = "AP-STA-POSSIBLE-PSK-MISMATCH"
	// ApStaPollOk as defined in wpactrl/wpa_ctrl.h:288
	ApStaPollOk = "AP-STA-POLL-OK"
	// ApRejectedMaxSta as defined in wpactrl/wpa_ctrl.h:290
	ApRejectedMaxSta = "AP-REJECTED-MAX-STA"
	// ApRejectedBlockedSta as defined in wpactrl/wpa_ctrl.h:291
	ApRejectedBlockedSta = "AP-REJECTED-BLOCKED-STA"
	// ApEventEnabled as defined in wpactrl/wpa_ctrl.h:293
	ApEventEnabled = "AP-ENABLED"
	// ApEventDisabled as defined in wpactrl/wpa_ctrl.h:294
	ApEventDisabled = "AP-DISABLED"
	// InterfaceEnabled as defined in wpactrl/wpa_ctrl.h:296
	InterfaceEnabled = "INTERFACE-ENABLED"
	// InterfaceDisabled as defined in wpactrl/wpa_ctrl.h:297
	InterfaceDisabled = "INTERFACE-DISABLED"
	// AcsEventStarted as defined in wpactrl/wpa_ctrl.h:299
	AcsEventStarted = "ACS-STARTED"
	// AcsEventCompleted as defined in wpactrl/wpa_ctrl.h:300
	AcsEventCompleted = "ACS-COMPLETED"
	// AcsEventFailed as defined in wpactrl/wpa_ctrl.h:301
	AcsEventFailed = "ACS-FAILED"
	// DfsEventRadarDetected as defined in wpactrl/wpa_ctrl.h:303
	DfsEventRadarDetected = "DFS-RADAR-DETECTED"
	// DfsEventNewChannel as defined in wpactrl/wpa_ctrl.h:304
	DfsEventNewChannel = "DFS-NEW-CHANNEL"
	// DfsEventCacStart as defined in wpactrl/wpa_ctrl.h:305
	DfsEventCacStart = "DFS-CAC-START"
	// DfsEventCacCompleted as defined in wpactrl/wpa_ctrl.h:306
	DfsEventCacCompleted = "DFS-CAC-COMPLETED"
	// DfsEventNopFinished as defined in wpactrl/wpa_ctrl.h:307
	DfsEventNopFinished = "DFS-NOP-FINISHED"
	// DfsEventPreCacExpired as defined in wpactrl/wpa_ctrl.h:308
	DfsEventPreCacExpired = "DFS-PRE-CAC-EXPIRED"
	// ApCsaFinished as defined in wpactrl/wpa_ctrl.h:310
	ApCsaFinished = "AP-CSA-FINISHED"
	// P2pEventListenOffloadStop as defined in wpactrl/wpa_ctrl.h:312
	P2pEventListenOffloadStop = "P2P-LISTEN-OFFLOAD-STOPPED"
	// P2pListenOffloadStopReason as defined in wpactrl/wpa_ctrl.h:313
	P2pListenOffloadStopReason = "P2P-LISTEN-OFFLOAD-STOP-REASON"
	// BssTmResp as defined in wpactrl/wpa_ctrl.h:316
	BssTmResp = "BSS-TM-RESP"
	// MboCellPreference as defined in wpactrl/wpa_ctrl.h:319
	MboCellPreference = "MBO-CELL-PREFERENCE"
	// MboTransitionReason as defined in wpactrl/wpa_ctrl.h:322
	MboTransitionReason = "MBO-TRANSITION-REASON"
	// PmksaCacheAdded as defined in wpactrl/wpa_ctrl.h:330
	PmksaCacheAdded = "PMKSA-CACHE-ADDED"
	// PmksaCacheRemoved as defined in wpactrl/wpa_ctrl.h:332
	PmksaCacheRemoved = "PMKSA-CACHE-REMOVED"
	// FilsHlpRx as defined in wpactrl/wpa_ctrl.h:336
	FilsHlpRx = "FILS-HLP-RX"
	// BssMaskAll as defined in wpactrl/wpa_ctrl.h:340
	BssMaskAll = 0xFFFDFFFF
	// CtrlIfacePort as defined in wpactrl/wpa_ctrl.h:541
	CtrlIfacePort = 9877
	// CtrlIfacePortLimit as defined in wpactrl/wpa_ctrl.h:542
	CtrlIfacePortLimit = 50
	// GlobalCtrlIfacePort as defined in wpactrl/wpa_ctrl.h:543
	GlobalCtrlIfacePort = 9878
	// GlobalCtrlIfacePortLimit as defined in wpactrl/wpa_ctrl.h:544
	GlobalCtrlIfacePortLimit = 20
)

// Event is an event that happens in the WPA supplicant
type Event struct {
	Name      string
	Arguments map[string]string
}

// NewEventFromMsg will create a new event from the given message
func NewEventFromMsg(msg string) (Event, error) {
	// control event, sent to the channel
	reader := csv.NewReader(strings.NewReader(msg))
	reader.Comma = ' '
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = false
	parts, err := reader.Read()
	if err != nil {
		log.Println("Error during parsing:", err)
	}
	if len(parts) == 0 {
		return Event{}, ErrInvalidEvent
	}

	event := Event{Name: parts[0][3:], Arguments: make(map[string]string)}
	for _, record := range parts[1:] {
		if strings.Contains(record, "=") {
			nvs := strings.SplitN(record, "=", 2)
			event.Arguments[nvs[0]] = nvs[1]
		}
	}
	return event, nil
}
