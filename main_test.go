package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/esiqveland/prototest/protos/org/noteroo"
	"github.com/esiqveland/prototest/protos/org/noteroo/kafka"
	"github.com/esiqveland/prototest/protos/org/noteroo/kafka/account"
	v1 "github.com/esiqveland/prototest/protos/org/noteroo/openapi/v1"
)

func TestServiceOption(t *testing.T) {
	opts, err := GetServiceOption(v1.File_org_noteroo_openapi_v1_noteroo_proto)
	require.NoError(t, err)
	require.EqualValues(t, "api.noteroo.org", opts.Service.Host)
	require.EqualValues(t, 443, opts.Service.Port)
}

func TestServiceOptionV2(t *testing.T) {
	opts, err := GetServiceOptionV2(v1.AccountService_ServiceDesc)
	require.NoError(t, err)
	require.EqualValues(t, "api.noteroo.org", opts.Service.Host)
	require.EqualValues(t, 443, opts.Service.Port)
}

func TestKafkaOption(t *testing.T) {
	opts, err := GetKafkaOption(&account.AccountEvent{})
	require.NoError(t, err)
	require.Equal(t, "account.events", opts.Topic)
}

func GetKafkaOption(protoMsg proto.Message) (*kafka.KafkaOptions, error) {
	descriptor := protoMsg.ProtoReflect().Descriptor()
	mOpts := descriptor.Options()

	if !proto.HasExtension(mOpts, kafka.E_Kafka) {
		return nil, fmt.Errorf("unable to find KafkaOptions on type: %v", protoMsg)
	}

	kExtOpts := proto.GetExtension(mOpts, kafka.E_Kafka)
	kOpts, ok := kExtOpts.(*kafka.KafkaOptions)
	if kOpts == nil || !ok {
		return nil, fmt.Errorf("unable to find KafkaOptions on type: %v", protoMsg)
	} else {
		return kOpts, nil
	}
}

// this uses the included compact representation of the original proto file to reflect and discover the option
func GetServiceOption(desc protoreflect.FileDescriptor) (*noteroo.NoterooServiceOptions, error) {
	svcs := desc.Services()
	num := svcs.Len()
	for i := 0; i < num; i++ {
		svc := svcs.Get(i)
		mOpts := svc.Options()

		if proto.HasExtension(mOpts, noteroo.E_Service) {
			kExtOpts := proto.GetExtension(mOpts, noteroo.E_Service)

			kOpts, ok := kExtOpts.(*noteroo.NoterooServiceOptions)
			if ok {
				return kOpts, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find options on type: %v", desc.Name())

}

// GetServiceOptionV2 uses the global protoregistry where all descriptors and types are registered on init (?)
// to discover the service options.
func GetServiceOptionV2(sd grpc.ServiceDesc) (*noteroo.NoterooServiceOptions, error) {
	des, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(sd.ServiceName))
	if err != nil {
		return nil, fmt.Errorf("unable to find options on servicename: %v: %w", sd.ServiceName, err)
	}
	desc, ok := des.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("unable to resolve options on servicename: %v: %w", sd.ServiceName, err)
	}

	mOpts := desc.Options()
	kExtOpts := proto.GetExtension(mOpts, noteroo.E_Service)

	kOpts, ok := kExtOpts.(*noteroo.NoterooServiceOptions)
	if kOpts == nil || !ok {
		return nil, fmt.Errorf("unable to find options on type: %v", desc.Name())
	} else {
		return kOpts, nil
	}
}
