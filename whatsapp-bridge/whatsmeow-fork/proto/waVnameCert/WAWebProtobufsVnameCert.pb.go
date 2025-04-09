// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.21.12
// source: waVnameCert/WAWebProtobufsVnameCert.proto

package waVnameCert

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "embed"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BizAccountLinkInfo_AccountType int32

const (
	BizAccountLinkInfo_ENTERPRISE BizAccountLinkInfo_AccountType = 0
)

// Enum value maps for BizAccountLinkInfo_AccountType.
var (
	BizAccountLinkInfo_AccountType_name = map[int32]string{
		0: "ENTERPRISE",
	}
	BizAccountLinkInfo_AccountType_value = map[string]int32{
		"ENTERPRISE": 0,
	}
)

func (x BizAccountLinkInfo_AccountType) Enum() *BizAccountLinkInfo_AccountType {
	p := new(BizAccountLinkInfo_AccountType)
	*p = x
	return p
}

func (x BizAccountLinkInfo_AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BizAccountLinkInfo_AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[0].Descriptor()
}

func (BizAccountLinkInfo_AccountType) Type() protoreflect.EnumType {
	return &file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[0]
}

func (x BizAccountLinkInfo_AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BizAccountLinkInfo_AccountType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BizAccountLinkInfo_AccountType(num)
	return nil
}

// Deprecated: Use BizAccountLinkInfo_AccountType.Descriptor instead.
func (BizAccountLinkInfo_AccountType) EnumDescriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{0, 0}
}

type BizAccountLinkInfo_HostStorageType int32

const (
	BizAccountLinkInfo_ON_PREMISE BizAccountLinkInfo_HostStorageType = 0
	BizAccountLinkInfo_FACEBOOK   BizAccountLinkInfo_HostStorageType = 1
)

// Enum value maps for BizAccountLinkInfo_HostStorageType.
var (
	BizAccountLinkInfo_HostStorageType_name = map[int32]string{
		0: "ON_PREMISE",
		1: "FACEBOOK",
	}
	BizAccountLinkInfo_HostStorageType_value = map[string]int32{
		"ON_PREMISE": 0,
		"FACEBOOK":   1,
	}
)

func (x BizAccountLinkInfo_HostStorageType) Enum() *BizAccountLinkInfo_HostStorageType {
	p := new(BizAccountLinkInfo_HostStorageType)
	*p = x
	return p
}

func (x BizAccountLinkInfo_HostStorageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BizAccountLinkInfo_HostStorageType) Descriptor() protoreflect.EnumDescriptor {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[1].Descriptor()
}

func (BizAccountLinkInfo_HostStorageType) Type() protoreflect.EnumType {
	return &file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[1]
}

func (x BizAccountLinkInfo_HostStorageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BizAccountLinkInfo_HostStorageType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BizAccountLinkInfo_HostStorageType(num)
	return nil
}

// Deprecated: Use BizAccountLinkInfo_HostStorageType.Descriptor instead.
func (BizAccountLinkInfo_HostStorageType) EnumDescriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{0, 1}
}

type BizIdentityInfo_ActualActorsType int32

const (
	BizIdentityInfo_SELF BizIdentityInfo_ActualActorsType = 0
	BizIdentityInfo_BSP  BizIdentityInfo_ActualActorsType = 1
)

// Enum value maps for BizIdentityInfo_ActualActorsType.
var (
	BizIdentityInfo_ActualActorsType_name = map[int32]string{
		0: "SELF",
		1: "BSP",
	}
	BizIdentityInfo_ActualActorsType_value = map[string]int32{
		"SELF": 0,
		"BSP":  1,
	}
)

func (x BizIdentityInfo_ActualActorsType) Enum() *BizIdentityInfo_ActualActorsType {
	p := new(BizIdentityInfo_ActualActorsType)
	*p = x
	return p
}

func (x BizIdentityInfo_ActualActorsType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BizIdentityInfo_ActualActorsType) Descriptor() protoreflect.EnumDescriptor {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[2].Descriptor()
}

func (BizIdentityInfo_ActualActorsType) Type() protoreflect.EnumType {
	return &file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[2]
}

func (x BizIdentityInfo_ActualActorsType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BizIdentityInfo_ActualActorsType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BizIdentityInfo_ActualActorsType(num)
	return nil
}

// Deprecated: Use BizIdentityInfo_ActualActorsType.Descriptor instead.
func (BizIdentityInfo_ActualActorsType) EnumDescriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{1, 0}
}

type BizIdentityInfo_HostStorageType int32

const (
	BizIdentityInfo_ON_PREMISE BizIdentityInfo_HostStorageType = 0
	BizIdentityInfo_FACEBOOK   BizIdentityInfo_HostStorageType = 1
)

// Enum value maps for BizIdentityInfo_HostStorageType.
var (
	BizIdentityInfo_HostStorageType_name = map[int32]string{
		0: "ON_PREMISE",
		1: "FACEBOOK",
	}
	BizIdentityInfo_HostStorageType_value = map[string]int32{
		"ON_PREMISE": 0,
		"FACEBOOK":   1,
	}
)

func (x BizIdentityInfo_HostStorageType) Enum() *BizIdentityInfo_HostStorageType {
	p := new(BizIdentityInfo_HostStorageType)
	*p = x
	return p
}

func (x BizIdentityInfo_HostStorageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BizIdentityInfo_HostStorageType) Descriptor() protoreflect.EnumDescriptor {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[3].Descriptor()
}

func (BizIdentityInfo_HostStorageType) Type() protoreflect.EnumType {
	return &file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[3]
}

func (x BizIdentityInfo_HostStorageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BizIdentityInfo_HostStorageType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BizIdentityInfo_HostStorageType(num)
	return nil
}

// Deprecated: Use BizIdentityInfo_HostStorageType.Descriptor instead.
func (BizIdentityInfo_HostStorageType) EnumDescriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{1, 1}
}

type BizIdentityInfo_VerifiedLevelValue int32

const (
	BizIdentityInfo_UNKNOWN BizIdentityInfo_VerifiedLevelValue = 0
	BizIdentityInfo_LOW     BizIdentityInfo_VerifiedLevelValue = 1
	BizIdentityInfo_HIGH    BizIdentityInfo_VerifiedLevelValue = 2
)

// Enum value maps for BizIdentityInfo_VerifiedLevelValue.
var (
	BizIdentityInfo_VerifiedLevelValue_name = map[int32]string{
		0: "UNKNOWN",
		1: "LOW",
		2: "HIGH",
	}
	BizIdentityInfo_VerifiedLevelValue_value = map[string]int32{
		"UNKNOWN": 0,
		"LOW":     1,
		"HIGH":    2,
	}
)

func (x BizIdentityInfo_VerifiedLevelValue) Enum() *BizIdentityInfo_VerifiedLevelValue {
	p := new(BizIdentityInfo_VerifiedLevelValue)
	*p = x
	return p
}

func (x BizIdentityInfo_VerifiedLevelValue) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BizIdentityInfo_VerifiedLevelValue) Descriptor() protoreflect.EnumDescriptor {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[4].Descriptor()
}

func (BizIdentityInfo_VerifiedLevelValue) Type() protoreflect.EnumType {
	return &file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes[4]
}

func (x BizIdentityInfo_VerifiedLevelValue) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BizIdentityInfo_VerifiedLevelValue) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BizIdentityInfo_VerifiedLevelValue(num)
	return nil
}

// Deprecated: Use BizIdentityInfo_VerifiedLevelValue.Descriptor instead.
func (BizIdentityInfo_VerifiedLevelValue) EnumDescriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{1, 2}
}

type BizAccountLinkInfo struct {
	state               protoimpl.MessageState              `protogen:"open.v1"`
	WhatsappBizAcctFbid *uint64                             `protobuf:"varint,1,opt,name=whatsappBizAcctFbid" json:"whatsappBizAcctFbid,omitempty"`
	WhatsappAcctNumber  *string                             `protobuf:"bytes,2,opt,name=whatsappAcctNumber" json:"whatsappAcctNumber,omitempty"`
	IssueTime           *uint64                             `protobuf:"varint,3,opt,name=issueTime" json:"issueTime,omitempty"`
	HostStorage         *BizAccountLinkInfo_HostStorageType `protobuf:"varint,4,opt,name=hostStorage,enum=WAWebProtobufsVnameCert.BizAccountLinkInfo_HostStorageType" json:"hostStorage,omitempty"`
	AccountType         *BizAccountLinkInfo_AccountType     `protobuf:"varint,5,opt,name=accountType,enum=WAWebProtobufsVnameCert.BizAccountLinkInfo_AccountType" json:"accountType,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *BizAccountLinkInfo) Reset() {
	*x = BizAccountLinkInfo{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BizAccountLinkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BizAccountLinkInfo) ProtoMessage() {}

func (x *BizAccountLinkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BizAccountLinkInfo.ProtoReflect.Descriptor instead.
func (*BizAccountLinkInfo) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{0}
}

func (x *BizAccountLinkInfo) GetWhatsappBizAcctFbid() uint64 {
	if x != nil && x.WhatsappBizAcctFbid != nil {
		return *x.WhatsappBizAcctFbid
	}
	return 0
}

func (x *BizAccountLinkInfo) GetWhatsappAcctNumber() string {
	if x != nil && x.WhatsappAcctNumber != nil {
		return *x.WhatsappAcctNumber
	}
	return ""
}

func (x *BizAccountLinkInfo) GetIssueTime() uint64 {
	if x != nil && x.IssueTime != nil {
		return *x.IssueTime
	}
	return 0
}

func (x *BizAccountLinkInfo) GetHostStorage() BizAccountLinkInfo_HostStorageType {
	if x != nil && x.HostStorage != nil {
		return *x.HostStorage
	}
	return BizAccountLinkInfo_ON_PREMISE
}

func (x *BizAccountLinkInfo) GetAccountType() BizAccountLinkInfo_AccountType {
	if x != nil && x.AccountType != nil {
		return *x.AccountType
	}
	return BizAccountLinkInfo_ENTERPRISE
}

type BizIdentityInfo struct {
	state           protoimpl.MessageState              `protogen:"open.v1"`
	Vlevel          *BizIdentityInfo_VerifiedLevelValue `protobuf:"varint,1,opt,name=vlevel,enum=WAWebProtobufsVnameCert.BizIdentityInfo_VerifiedLevelValue" json:"vlevel,omitempty"`
	VnameCert       *VerifiedNameCertificate            `protobuf:"bytes,2,opt,name=vnameCert" json:"vnameCert,omitempty"`
	Signed          *bool                               `protobuf:"varint,3,opt,name=signed" json:"signed,omitempty"`
	Revoked         *bool                               `protobuf:"varint,4,opt,name=revoked" json:"revoked,omitempty"`
	HostStorage     *BizIdentityInfo_HostStorageType    `protobuf:"varint,5,opt,name=hostStorage,enum=WAWebProtobufsVnameCert.BizIdentityInfo_HostStorageType" json:"hostStorage,omitempty"`
	ActualActors    *BizIdentityInfo_ActualActorsType   `protobuf:"varint,6,opt,name=actualActors,enum=WAWebProtobufsVnameCert.BizIdentityInfo_ActualActorsType" json:"actualActors,omitempty"`
	PrivacyModeTS   *uint64                             `protobuf:"varint,7,opt,name=privacyModeTS" json:"privacyModeTS,omitempty"`
	FeatureControls *uint64                             `protobuf:"varint,8,opt,name=featureControls" json:"featureControls,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *BizIdentityInfo) Reset() {
	*x = BizIdentityInfo{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BizIdentityInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BizIdentityInfo) ProtoMessage() {}

func (x *BizIdentityInfo) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BizIdentityInfo.ProtoReflect.Descriptor instead.
func (*BizIdentityInfo) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{1}
}

func (x *BizIdentityInfo) GetVlevel() BizIdentityInfo_VerifiedLevelValue {
	if x != nil && x.Vlevel != nil {
		return *x.Vlevel
	}
	return BizIdentityInfo_UNKNOWN
}

func (x *BizIdentityInfo) GetVnameCert() *VerifiedNameCertificate {
	if x != nil {
		return x.VnameCert
	}
	return nil
}

func (x *BizIdentityInfo) GetSigned() bool {
	if x != nil && x.Signed != nil {
		return *x.Signed
	}
	return false
}

func (x *BizIdentityInfo) GetRevoked() bool {
	if x != nil && x.Revoked != nil {
		return *x.Revoked
	}
	return false
}

func (x *BizIdentityInfo) GetHostStorage() BizIdentityInfo_HostStorageType {
	if x != nil && x.HostStorage != nil {
		return *x.HostStorage
	}
	return BizIdentityInfo_ON_PREMISE
}

func (x *BizIdentityInfo) GetActualActors() BizIdentityInfo_ActualActorsType {
	if x != nil && x.ActualActors != nil {
		return *x.ActualActors
	}
	return BizIdentityInfo_SELF
}

func (x *BizIdentityInfo) GetPrivacyModeTS() uint64 {
	if x != nil && x.PrivacyModeTS != nil {
		return *x.PrivacyModeTS
	}
	return 0
}

func (x *BizIdentityInfo) GetFeatureControls() uint64 {
	if x != nil && x.FeatureControls != nil {
		return *x.FeatureControls
	}
	return 0
}

type LocalizedName struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lg            *string                `protobuf:"bytes,1,opt,name=lg" json:"lg,omitempty"`
	Lc            *string                `protobuf:"bytes,2,opt,name=lc" json:"lc,omitempty"`
	VerifiedName  *string                `protobuf:"bytes,3,opt,name=verifiedName" json:"verifiedName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocalizedName) Reset() {
	*x = LocalizedName{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocalizedName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalizedName) ProtoMessage() {}

func (x *LocalizedName) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalizedName.ProtoReflect.Descriptor instead.
func (*LocalizedName) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{2}
}

func (x *LocalizedName) GetLg() string {
	if x != nil && x.Lg != nil {
		return *x.Lg
	}
	return ""
}

func (x *LocalizedName) GetLc() string {
	if x != nil && x.Lc != nil {
		return *x.Lc
	}
	return ""
}

func (x *LocalizedName) GetVerifiedName() string {
	if x != nil && x.VerifiedName != nil {
		return *x.VerifiedName
	}
	return ""
}

type VerifiedNameCertificate struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Details         []byte                 `protobuf:"bytes,1,opt,name=details" json:"details,omitempty"`
	Signature       []byte                 `protobuf:"bytes,2,opt,name=signature" json:"signature,omitempty"`
	ServerSignature []byte                 `protobuf:"bytes,3,opt,name=serverSignature" json:"serverSignature,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *VerifiedNameCertificate) Reset() {
	*x = VerifiedNameCertificate{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifiedNameCertificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifiedNameCertificate) ProtoMessage() {}

func (x *VerifiedNameCertificate) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifiedNameCertificate.ProtoReflect.Descriptor instead.
func (*VerifiedNameCertificate) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{3}
}

func (x *VerifiedNameCertificate) GetDetails() []byte {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *VerifiedNameCertificate) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *VerifiedNameCertificate) GetServerSignature() []byte {
	if x != nil {
		return x.ServerSignature
	}
	return nil
}

type BizAccountPayload struct {
	state           protoimpl.MessageState   `protogen:"open.v1"`
	VnameCert       *VerifiedNameCertificate `protobuf:"bytes,1,opt,name=vnameCert" json:"vnameCert,omitempty"`
	BizAcctLinkInfo []byte                   `protobuf:"bytes,2,opt,name=bizAcctLinkInfo" json:"bizAcctLinkInfo,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *BizAccountPayload) Reset() {
	*x = BizAccountPayload{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BizAccountPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BizAccountPayload) ProtoMessage() {}

func (x *BizAccountPayload) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BizAccountPayload.ProtoReflect.Descriptor instead.
func (*BizAccountPayload) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{4}
}

func (x *BizAccountPayload) GetVnameCert() *VerifiedNameCertificate {
	if x != nil {
		return x.VnameCert
	}
	return nil
}

func (x *BizAccountPayload) GetBizAcctLinkInfo() []byte {
	if x != nil {
		return x.BizAcctLinkInfo
	}
	return nil
}

type VerifiedNameCertificate_Details struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Serial         *uint64                `protobuf:"varint,1,opt,name=serial" json:"serial,omitempty"`
	Issuer         *string                `protobuf:"bytes,2,opt,name=issuer" json:"issuer,omitempty"`
	VerifiedName   *string                `protobuf:"bytes,4,opt,name=verifiedName" json:"verifiedName,omitempty"`
	LocalizedNames []*LocalizedName       `protobuf:"bytes,8,rep,name=localizedNames" json:"localizedNames,omitempty"`
	IssueTime      *uint64                `protobuf:"varint,10,opt,name=issueTime" json:"issueTime,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *VerifiedNameCertificate_Details) Reset() {
	*x = VerifiedNameCertificate_Details{}
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifiedNameCertificate_Details) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifiedNameCertificate_Details) ProtoMessage() {}

func (x *VerifiedNameCertificate_Details) ProtoReflect() protoreflect.Message {
	mi := &file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifiedNameCertificate_Details.ProtoReflect.Descriptor instead.
func (*VerifiedNameCertificate_Details) Descriptor() ([]byte, []int) {
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP(), []int{3, 0}
}

func (x *VerifiedNameCertificate_Details) GetSerial() uint64 {
	if x != nil && x.Serial != nil {
		return *x.Serial
	}
	return 0
}

func (x *VerifiedNameCertificate_Details) GetIssuer() string {
	if x != nil && x.Issuer != nil {
		return *x.Issuer
	}
	return ""
}

func (x *VerifiedNameCertificate_Details) GetVerifiedName() string {
	if x != nil && x.VerifiedName != nil {
		return *x.VerifiedName
	}
	return ""
}

func (x *VerifiedNameCertificate_Details) GetLocalizedNames() []*LocalizedName {
	if x != nil {
		return x.LocalizedNames
	}
	return nil
}

func (x *VerifiedNameCertificate_Details) GetIssueTime() uint64 {
	if x != nil && x.IssueTime != nil {
		return *x.IssueTime
	}
	return 0
}

var File_waVnameCert_WAWebProtobufsVnameCert_proto protoreflect.FileDescriptor

//go:embed WAWebProtobufsVnameCert.pb.raw
var file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDesc []byte

var (
	file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescOnce sync.Once
	file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescData = file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDesc
)

func file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescGZIP() []byte {
	file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescOnce.Do(func() {
		file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescData = protoimpl.X.CompressGZIP(file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescData)
	})
	return file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDescData
}

var file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_waVnameCert_WAWebProtobufsVnameCert_proto_goTypes = []any{
	(BizAccountLinkInfo_AccountType)(0),     // 0: WAWebProtobufsVnameCert.BizAccountLinkInfo.AccountType
	(BizAccountLinkInfo_HostStorageType)(0), // 1: WAWebProtobufsVnameCert.BizAccountLinkInfo.HostStorageType
	(BizIdentityInfo_ActualActorsType)(0),   // 2: WAWebProtobufsVnameCert.BizIdentityInfo.ActualActorsType
	(BizIdentityInfo_HostStorageType)(0),    // 3: WAWebProtobufsVnameCert.BizIdentityInfo.HostStorageType
	(BizIdentityInfo_VerifiedLevelValue)(0), // 4: WAWebProtobufsVnameCert.BizIdentityInfo.VerifiedLevelValue
	(*BizAccountLinkInfo)(nil),              // 5: WAWebProtobufsVnameCert.BizAccountLinkInfo
	(*BizIdentityInfo)(nil),                 // 6: WAWebProtobufsVnameCert.BizIdentityInfo
	(*LocalizedName)(nil),                   // 7: WAWebProtobufsVnameCert.LocalizedName
	(*VerifiedNameCertificate)(nil),         // 8: WAWebProtobufsVnameCert.VerifiedNameCertificate
	(*BizAccountPayload)(nil),               // 9: WAWebProtobufsVnameCert.BizAccountPayload
	(*VerifiedNameCertificate_Details)(nil), // 10: WAWebProtobufsVnameCert.VerifiedNameCertificate.Details
}
var file_waVnameCert_WAWebProtobufsVnameCert_proto_depIdxs = []int32{
	1, // 0: WAWebProtobufsVnameCert.BizAccountLinkInfo.hostStorage:type_name -> WAWebProtobufsVnameCert.BizAccountLinkInfo.HostStorageType
	0, // 1: WAWebProtobufsVnameCert.BizAccountLinkInfo.accountType:type_name -> WAWebProtobufsVnameCert.BizAccountLinkInfo.AccountType
	4, // 2: WAWebProtobufsVnameCert.BizIdentityInfo.vlevel:type_name -> WAWebProtobufsVnameCert.BizIdentityInfo.VerifiedLevelValue
	8, // 3: WAWebProtobufsVnameCert.BizIdentityInfo.vnameCert:type_name -> WAWebProtobufsVnameCert.VerifiedNameCertificate
	3, // 4: WAWebProtobufsVnameCert.BizIdentityInfo.hostStorage:type_name -> WAWebProtobufsVnameCert.BizIdentityInfo.HostStorageType
	2, // 5: WAWebProtobufsVnameCert.BizIdentityInfo.actualActors:type_name -> WAWebProtobufsVnameCert.BizIdentityInfo.ActualActorsType
	8, // 6: WAWebProtobufsVnameCert.BizAccountPayload.vnameCert:type_name -> WAWebProtobufsVnameCert.VerifiedNameCertificate
	7, // 7: WAWebProtobufsVnameCert.VerifiedNameCertificate.Details.localizedNames:type_name -> WAWebProtobufsVnameCert.LocalizedName
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_waVnameCert_WAWebProtobufsVnameCert_proto_init() }
func file_waVnameCert_WAWebProtobufsVnameCert_proto_init() {
	if File_waVnameCert_WAWebProtobufsVnameCert_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_waVnameCert_WAWebProtobufsVnameCert_proto_goTypes,
		DependencyIndexes: file_waVnameCert_WAWebProtobufsVnameCert_proto_depIdxs,
		EnumInfos:         file_waVnameCert_WAWebProtobufsVnameCert_proto_enumTypes,
		MessageInfos:      file_waVnameCert_WAWebProtobufsVnameCert_proto_msgTypes,
	}.Build()
	File_waVnameCert_WAWebProtobufsVnameCert_proto = out.File
	file_waVnameCert_WAWebProtobufsVnameCert_proto_rawDesc = nil
	file_waVnameCert_WAWebProtobufsVnameCert_proto_goTypes = nil
	file_waVnameCert_WAWebProtobufsVnameCert_proto_depIdxs = nil
}
