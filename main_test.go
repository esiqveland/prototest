package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/esiqveland/prototest/protos/org/noteroo/kafka"
	"github.com/esiqveland/prototest/protos/org/noteroo/kafka/noteroo"
)

func TestKafkaOption(t *testing.T) {
	opts, err := GetKafkaOption(&noteroo.AccountEvent{})
	require.NoError(t, err)
	require.Equal(t, "account.events", opts.Topic)
}

func GetKafkaOption(protoMsg proto.Message) (*kafka.KafkaOptions, error) {
	descriptor := protoMsg.ProtoReflect().Descriptor()
	mOpts := descriptor.Options()
	kExtOpts := proto.GetExtension(mOpts, kafka.E_Kafka)

	kOpts, ok := kExtOpts.(*kafka.KafkaOptions)
	if kOpts == nil || !ok {
		return nil, fmt.Errorf("unable to find KafkaOptions on type: %v", protoMsg)

	} else {
		return kOpts, nil
	}
}
