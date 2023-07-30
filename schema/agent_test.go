package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToolInput(t *testing.T) {
	t.Run("TestToolInput_GetString", func(t *testing.T) {
		// Test getting string from a plain string ToolInput.
		plainInput := NewToolInputFromString("plain string input")
		str, err := plainInput.GetString()
		require.NoError(t, err)
		require.Equal(t, "plain string input", str)

		// Test getting string from a structured ToolInput.
		structuredInput := NewToolInputFromArguments(`{"__arg1": "structured input"}`)
		_, err = structuredInput.GetString()
		require.Error(t, err)
		require.EqualError(t, err, "cannot return string for structured input")
	})

	t.Run("TestToolInput_Unmarshal", func(t *testing.T) {
		// Test unmarshaling a plain string ToolInput into a string.
		plainInput := NewToolInputFromString("plain string input")
		var str string
		err := plainInput.Unmarshal(&str)
		require.NoError(t, err)
		require.Equal(t, "plain string input", str)

		// Test unmarshaling a structured ToolInput into a string.
		structuredInput := NewToolInputFromArguments(`{"__arg1": "structured input"}`)
		var strFromStructured string
		err = structuredInput.Unmarshal(&strFromStructured)
		require.NoError(t, err)
		require.Equal(t, "structured input", strFromStructured)

		// Test unmarshaling a structured ToolInput into a custom struct.
		structuredInput2 := NewToolInputFromArguments(`{"key": "value"}`)
		var customStruct struct{ Key string }
		err = structuredInput2.Unmarshal(&customStruct)
		require.NoError(t, err)
		require.Equal(t, "value", customStruct.Key)
	})

	t.Run("TestToolInput_String", func(t *testing.T) {
		// Test String method of ToolInput.
		input := NewToolInputFromString("test input")
		require.Equal(t, "test input", input.String())
	})
}
