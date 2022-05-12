package message

import (
	"time"
)

// DeliverySetting is used to configure registered delivery
// for short messages.
type DeliverySetting uint8

// Supported delivery settings.
const (
	NoDeliveryReceipt      DeliverySetting = 0x00
	FinalDeliveryReceipt   DeliverySetting = 0x01
	FailureDeliveryReceipt DeliverySetting = 0x02
)

type ShortMessage struct {
	Type       ShortMessageType
	Src        string
	Dst        string
	DstList    []string // List of destination addreses for submit multi
	DLs        []string //List if destribution list for submit multi
	Text       []byte
	DataCoding DataCoding
	Validity   time.Duration
	Register   DeliverySetting

	// Other fields, normally optional.
	TLVFields            map[Tag][]byte
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	DestAddrTON          uint8
	DestAddrNPI          uint8
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ReplaceIfPresentFlag uint8
	SMDefaultMsgID       uint8
	NumberDests          uint8
}

type ShortMessageResp struct {
	Type      ShortMessageType
	Status    Status
	MessageID string
}

type ShortMessageType uint8
type Status uint8

const (
	SubmitSM ShortMessageType = 0x00
)

type DataCoding uint8

// Supported text codecs.
const (
	DefaultType DataCoding = 0x00 // SMSC Default Alphabet
	//	IA5Type       DataCoding = 0x01 // IA5 (CCITT T.50)/ASCII (ANSI X3.4)
	//	BinaryType    DataCoding = 0x02 // Octet unspecified (8-bit binary)
	Latin1Type DataCoding = 0x03 // Latin 1 (ISO-8859-1)
	//	Binary2Type   DataCoding = 0x04 // Octet unspecified (8-bit binary)
	//	JISType       DataCoding = 0x05 // JIS (X 0208-1990)
	ISO88595Type DataCoding = 0x06 // Cyrillic (ISO-8859-5)
	//	ISO88598Type  DataCoding = 0x07 // Latin/Hebrew (ISO-8859-8)
	UCS2Type DataCoding = 0x08 // UCS2 (ISO/IEC-10646)
	//	PictogramType DataCoding = 0x09 // Pictogram Encoding
	//	ISO2022JPType DataCoding = 0x0A // ISO-2022-JP (Music Codes)
	//	EXTJISType    DataCoding = 0x0D // Extended Kanji JIS (X 0212-1990)
	//	KSC5601Type   DataCoding = 0x0E // KS C 5601
)

// Tag is the tag of a Tag-Length-Value (TLV) field.
type Tag uint16

// Common Tag-Length-Value (TLV) tags.
const (
	TagDestAddrSubunit          Tag = 0x0005
	TagDestNetworkType          Tag = 0x0006
	TagDestBearerType           Tag = 0x0007
	TagDestTelematicsID         Tag = 0x0008
	TagSourceAddrSubunit        Tag = 0x000D
	TagSourceNetworkType        Tag = 0x000E
	TagSourceBearerType         Tag = 0x000F
	TagSourceTelematicsID       Tag = 0x0010
	TagQosTimeToLive            Tag = 0x0017
	TagPayloadType              Tag = 0x0019
	TagAdditionalStatusInfoText Tag = 0x001D
	TagReceiptedMessageID       Tag = 0x001E
	TagMsMsgWaitFacilities      Tag = 0x0030
	TagPrivacyIndicator         Tag = 0x0201
	TagSourceSubaddress         Tag = 0x0202
	TagDestSubaddress           Tag = 0x0203
	TagUserMessageReference     Tag = 0x0204
	TagUserResponseCode         Tag = 0x0205
	TagSourcePort               Tag = 0x020A
	TagDestinationPort          Tag = 0x020B
	TagSarMsgRefNum             Tag = 0x020C
	TagLanguageIndicator        Tag = 0x020D
	TagSarTotalSegments         Tag = 0x020E
	TagSarSegmentSeqnum         Tag = 0x020F
	TagCallbackNumPresInd       Tag = 0x0302
	TagCallbackNumAtag          Tag = 0x0303
	TagNumberOfMessages         Tag = 0x0304
	TagCallbackNum              Tag = 0x0381
	TagDpfResult                Tag = 0x0420
	TagSetDpf                   Tag = 0x0421
	TagMsAvailabilityStatus     Tag = 0x0422
	TagNetworkErrorCode         Tag = 0x0423
	TagMessagePayload           Tag = 0x0424
	TagDeliveryFailureReason    Tag = 0x0425
	TagMoreMessagesToSend       Tag = 0x0426
	TagMessageStateOption       Tag = 0x0427
	TagUssdServiceOp            Tag = 0x0501
	TagDisplayTime              Tag = 0x1201
	TagSmsSignal                Tag = 0x1203
	TagMsValidity               Tag = 0x1204
	TagAlertOnMessageDelivery   Tag = 0x130C
	TagItsReplyType             Tag = 0x1380
	TagItsSessionInfo           Tag = 0x1383
)

// These are status codes used on application level
const (
	Status_OK                                Status = 0x00000000
	Status_InvalidPriorityFlag               Status = 0x00000006
	Status_InvalidRegistredDeliveryFlag      Status = 0x00000007
	Status_SystemError                       Status = 0x00000008
	Status_InvalidSourceAddress              Status = 0x0000000a
	Status_InvalidDestinationAddress         Status = 0x0000000b
	Status_InvalidMessageId                  Status = 0x0000000c
	Status_CancelSmFailed                    Status = 0x00000011
	Status_ReplaceSmFailed                   Status = 0x00000013
	Status_MessageQueueFull                  Status = 0x00000014
	Status_InvalidServiceType                Status = 0x00000015
	Status_InvalidNumberOfDestinations       Status = 0x00000033
	Status_InvalidDistributionListName       Status = 0x00000034
	Status_InvalidDestinationFlag            Status = 0x00000040
	Status_InvalidSubmitWithReplaceRequest   Status = 0x00000042
	Status_InvalidEsmClassFieldData          Status = 0x00000043
	Status_CannotSubmitToDistList            Status = 0x00000044
	Status_SubmitSmFailed                    Status = 0x00000045
	Status_InvalidSourceAddressTon           Status = 0x00000048
	Status_InvalidSourceAddressNpi           Status = 0x00000049
	Status_InvalidSDestinationAddressTon     Status = 0x00000050
	Status_InvalidSDestinationAddressNpi     Status = 0x00000051
	Status_InvalidSystemType                 Status = 0x00000053
	Status_InvalidReplaceIfPresentFlag       Status = 0x00000054
	Status_InvalidNumberOfMessages           Status = 0x00000055
	Status_ThrottlingError                   Status = 0x00000058
	Status_InvalidScheduledDeliveryTime      Status = 0x00000061
	Status_InvalidMessageValidityTime        Status = 0x00000062
	Status_PredfinedMessageInvalidOrNotFound Status = 0x00000063
	Status_EsmeReceiverTempAppError          Status = 0x00000064
	Status_EsmeReceiverPermanentAppError     Status = 0x00000065
	Status_EsmeReceiverRejectMessageError    Status = 0x00000066
	Status_QuerySmFailed                     Status = 0x00000067
	Status_OptionalParameterNotAllowed       Status = 0x000000c1
	Status_InvalidParameterLength            Status = 0x000000c2
	Status_ExpectedOptionalParameterMissing  Status = 0x000000c3
	Status_InvalidOptionalParameterValue     Status = 0x000000c4
	Status_UnknownError                      Status = 0x000000ff
)
