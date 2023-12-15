package util

import "math"

func CosineSimilarity(matrix1, matrix2 [][]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for i := 0; i < len(matrix1); i++ {
		for j := 0; j < len(matrix1[0]); j++ {
			dotProduct += matrix1[i][j] * matrix2[i][j]
			magnitude1 += math.Pow(matrix1[i][j], 2)
			magnitude2 += math.Pow(matrix2[i][j], 2)
		}
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0 // Handle zero magnitude case
	}

	return dotProduct / (magnitude1 * magnitude2)
}
