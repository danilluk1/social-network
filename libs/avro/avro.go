package avro

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/linkedin/goavro/v2"
)

func Encode[T any](data *T, codec *goavro.Codec, schemaId int) ([]byte, error) {
	if codec == nil {
		return []byte{}, errors.New("codec cannot be empty")
	}

	value, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	native, _, err := codec.NativeFromTextual(value)
	if err != nil {
		return []byte{}, err
	}
	avroBytes, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return []byte{}, err
	}

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schemaId))
	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, avroBytes...)

	return recordValue, nil
}

func Decode[T any](data []byte, codec *goavro.Codec, decoded *T) error {
	if len(data) <= 5 {
		return errors.New("mailformed message")
	}

	native, _, err := codec.NativeFromBinary(data[5:])
	if err != nil {
		return err
	}
	value, err := codec.TextualFromNative(nil, native)
	if err != nil {
		return err
	}
	err = json.Unmarshal(value, &decoded)
	if err != nil {
		return err
	}

	return nil
}
