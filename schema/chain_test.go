package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainValues(t *testing.T) {
	t.Run("TestGetString", func(t *testing.T) {
		cv := ChainValues{
			"name": "John",
			"age":  30,
		}

		// Test getting a valid string value
		name, err := cv.GetString("name")
		require.NoError(t, err)
		require.Equal(t, "John", name)

		// Test getting an invalid string value
		_, err = cv.GetString("age")
		require.ErrorIs(t, err, ErrInputValuesWrongType)

		// Test getting a non-existent key
		_, err = cv.GetString("address")
		require.ErrorIs(t, err, ErrInvalidInputValues)
	})

	t.Run("TestGetDocuments", func(t *testing.T) {
		doc1 := Document{
			PageContent: "Document 1",
		}
		doc2 := Document{
			PageContent: "Document 2",
		}

		cv := ChainValues{
			"docs": []Document{doc1, doc2},
		}

		// Test getting valid documents
		docs, err := cv.GetDocuments("docs")
		require.NoError(t, err)
		require.Len(t, docs, 2)

		// Test getting an invalid document value (not a slice of documents)
		cv["docs"] = "invalid"
		_, err = cv.GetDocuments("docs")
		require.ErrorIs(t, err, ErrInputValuesWrongType)

		// Test getting a non-existent key
		_, err = cv.GetDocuments("files")
		require.ErrorIs(t, err, ErrInvalidInputValues)

		// Test getting an empty slice of documents
		cv["docs"] = []Document{}
		_, err = cv.GetDocuments("docs")
		require.ErrorIs(t, err, ErrInvalidInputValues)
	})

	t.Run("TestClone", func(t *testing.T) {
		// Create a sample ChainValues map
		cv := ChainValues{
			"key1": "value1",
			"key2": "value2",
			"key3": 123,
		}

		// Call Clone to create a shallow copy
		clone := cv.Clone()

		// Assert that the cloned map is equal to the original map
		require.Equal(t, cv, clone)

		// Modify the cloned map
		clone["key1"] = "modified value"

		// Assert that the original map is not affected by the modification to the clone
		require.NotEqual(t, cv, clone)
	})

	t.Run("TestClone_Empty", func(t *testing.T) {
		// Create an empty ChainValues map
		cv := ChainValues{}

		// Call Clone to create a shallow copy
		clone := cv.Clone()

		// Assert that the cloned map is equal to the original map (both should be empty)
		require.Equal(t, cv, clone)

		// Modify the cloned map
		clone["key1"] = "value1"

		// Assert that the original map is not affected by the modification to the clone
		require.NotEqual(t, cv, clone)
	})
}
