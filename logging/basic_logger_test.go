package logging

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicLogger_ManualEmit(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	logger := NewBasicLogger(WithWriter(buf))

	emitter := logger.Emitter()

	logger.Add("field1", "value1")
	logger.Add("field2", 2)
	logger.Infof("message 1: %s", "hello world")
	logger.Debugf("message 2: %s", "invisible message")

	require.NoError(t, emitter())

	// we expect the same structure
	output := map[string]interface{}{}

	require.NoError(t, json.Unmarshal(buf.Bytes(), &output))

	require.EqualValues(t, "value1", output["field1"])
	require.EqualValues(t, 2, output["field2"])

	require.Len(t, output["lines"], 1)
	require.Equal(t, "message 1: hello world", output["lines"].([]interface{})[0].(map[string]interface{})["message"])
}
