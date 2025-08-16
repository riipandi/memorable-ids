package memorable_ids

// Run these tests using:
// gotestsum --format short-verbose -- ./pkg/memorable-ids -v

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("should generate ID with default options (2 components)", func(t *testing.T) {
		id, err := Generate(GenerateOptions{})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 2, "Expected 2 parts")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
		assert.True(t, contains(Nouns, parts[1]), "Second part '%s' not found in nouns", parts[1])
	})

	t.Run("should generate ID with 1 component", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 1})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 1, "Expected 1 part")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
	})

	t.Run("should generate ID with 3 components", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 3})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 3, "Expected 3 parts")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
		assert.True(t, contains(Nouns, parts[1]), "Second part '%s' not found in nouns", parts[1])
		assert.True(t, contains(Verbs, parts[2]), "Third part '%s' not found in verbs", parts[2])
	})

	t.Run("should generate ID with 4 components", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 4})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 4, "Expected 4 parts")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
		assert.True(t, contains(Nouns, parts[1]), "Second part '%s' not found in nouns", parts[1])
		assert.True(t, contains(Verbs, parts[2]), "Third part '%s' not found in verbs", parts[2])
		assert.True(t, contains(Adverbs, parts[3]), "Fourth part '%s' not found in adverbs", parts[3])
	})

	t.Run("should generate ID with 5 components", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 5})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 5, "Expected 5 parts")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
		assert.True(t, contains(Nouns, parts[1]), "Second part '%s' not found in nouns", parts[1])
		assert.True(t, contains(Verbs, parts[2]), "Third part '%s' not found in verbs", parts[2])
		assert.True(t, contains(Adverbs, parts[3]), "Fourth part '%s' not found in adverbs", parts[3])
		assert.True(t, contains(Prepositions, parts[4]), "Fifth part '%s' not found in prepositions", parts[4])
	})

	t.Run("should use custom separator", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 2, Separator: "_"})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "_")
		assert.Len(t, parts, 2, "Expected 2 parts")

		assert.Contains(t, id, "_", "ID should contain underscore separator")

		assert.True(t, contains(Adjectives, parts[0]), "First part '%s' not found in adjectives", parts[0])
		assert.True(t, contains(Nouns, parts[1]), "Second part '%s' not found in nouns", parts[1])
	})

	t.Run("should add suffix when provided", func(t *testing.T) {
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     SuffixGenerators.Number,
		})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 3, "Expected 3 parts")

		// Check if last part is 3-digit number
		matched, _ := regexp.MatchString(`^\d{3}$`, parts[2])
		assert.True(t, matched, "Expected 3-digit number suffix, got '%s'", parts[2])
	})

	t.Run("should handle null suffix gracefully", func(t *testing.T) {
		nullSuffix := func() *string { return nil }
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     nullSuffix,
		})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 2, "Expected 2 parts (no suffix added)")
	})

	t.Run("should handle suffix that is not a function", func(t *testing.T) {
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     nil,
		})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 2, "Expected 2 parts (no suffix added)")
	})

	t.Run("should throw error for invalid component count", func(t *testing.T) {
		// Test cases that should return errors (excluding 0 since it defaults to 2)
		testCases := []int{-1, 6, 10}
		for _, components := range testCases {
			_, err := Generate(GenerateOptions{Components: components})
			assert.Error(t, err, "Expected error for components=%d", components)
			assert.Contains(t, err.Error(), "components must be between 1 and 5", "Expected specific error message")
		}

		// Test that 0 components defaults to 2 (no error)
		id, err := Generate(GenerateOptions{Components: 0})
		require.NoError(t, err, "Expected no error for components=0 (should default)")
		parts := strings.Split(id, "-")
		assert.Len(t, parts, 2, "Expected 2 parts for components=0 (default)")
	})

	t.Run("should generate different IDs on multiple calls", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id, err := Generate(GenerateOptions{})
			require.NoError(t, err, "Generate should not fail")
			ids[id] = true
		}

		// Should have high uniqueness (allowing for some collisions)
		assert.Greater(t, len(ids), 90, "Expected high uniqueness (>90), got %d unique IDs out of 100", len(ids))
	})
}

func TestParse(t *testing.T) {
	t.Run("should parse ID without suffix", func(t *testing.T) {
		result := Parse("cute-rabbit", "-")

		expectedComponents := []string{"cute", "rabbit"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		assert.Nil(t, result.Suffix, "Expected nil suffix")
	})

	t.Run("should parse ID with numeric suffix", func(t *testing.T) {
		result := Parse("cute-rabbit-042", "-")

		expectedComponents := []string{"cute", "rabbit"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "042", *result.Suffix, "Expected suffix '042'")
	})

	t.Run("should parse ID with non-numeric suffix as component", func(t *testing.T) {
		result := Parse("cute-rabbit-swim", "-")

		expectedComponents := []string{"cute", "rabbit", "swim"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		assert.Nil(t, result.Suffix, "Expected nil suffix")
	})

	t.Run("should handle custom separator", func(t *testing.T) {
		result := Parse("cute_rabbit_123", "_")

		expectedComponents := []string{"cute", "rabbit"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "123", *result.Suffix, "Expected suffix '123'")
	})

	t.Run("should handle single component", func(t *testing.T) {
		result := Parse("cute", "-")

		expectedComponents := []string{"cute"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		assert.Nil(t, result.Suffix, "Expected nil suffix")
	})

	t.Run("should handle single component with suffix", func(t *testing.T) {
		result := Parse("cute-123", "-")

		expectedComponents := []string{"cute"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "123", *result.Suffix, "Expected suffix '123'")
	})

	t.Run("should handle ID with only numeric part", func(t *testing.T) {
		result := Parse("123", "-")

		expectedComponents := []string{}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected empty components")
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "123", *result.Suffix, "Expected suffix '123'")
	})

	t.Run("should handle mixed numeric patterns", func(t *testing.T) {
		result := Parse("cute-123abc-456", "-")

		expectedComponents := []string{"cute", "123abc"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "456", *result.Suffix, "Expected suffix '456'")
	})
}

func TestSuffixGenerators(t *testing.T) {
	t.Run("number should generate 3-digit string", func(t *testing.T) {
		suffix := SuffixGenerators.Number()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^\d{3}$`, *suffix)
		assert.True(t, matched, "Expected 3-digit string, got '%s'", *suffix)
	})

	t.Run("number4 should generate 4-digit string", func(t *testing.T) {
		suffix := SuffixGenerators.Number4()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^\d{4}$`, *suffix)
		assert.True(t, matched, "Expected 4-digit string, got '%s'", *suffix)
	})

	t.Run("hex should generate 2-digit hex string", func(t *testing.T) {
		suffix := SuffixGenerators.Hex()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^[0-9a-f]{2}$`, *suffix)
		assert.True(t, matched, "Expected 2-digit hex string, got '%s'", *suffix)
	})

	t.Run("timestamp should generate 4-digit string", func(t *testing.T) {
		suffix := SuffixGenerators.Timestamp()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^\d{4}$`, *suffix)
		assert.True(t, matched, "Expected 4-digit string, got '%s'", *suffix)
	})

	t.Run("letter should generate single lowercase letter", func(t *testing.T) {
		suffix := SuffixGenerators.Letter()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^[a-z]$`, *suffix)
		assert.True(t, matched, "Expected single lowercase letter, got '%s'", *suffix)
	})

	t.Run("defaultSuffix should work same as number generator", func(t *testing.T) {
		suffix := DefaultSuffix()
		require.NotNil(t, suffix, "Expected non-nil suffix")

		matched, _ := regexp.MatchString(`^\d{3}$`, *suffix)
		assert.True(t, matched, "Expected 3-digit string, got '%s'", *suffix)
	})

	t.Run("all suffix generators should produce valid output", func(t *testing.T) {
		// Test edge cases for all generators
		for i := 0; i < 10; i++ {
			assert.NotNil(t, SuffixGenerators.Number(), "Number generator returned nil")
			assert.NotNil(t, SuffixGenerators.Number4(), "Number4 generator returned nil")
			assert.NotNil(t, SuffixGenerators.Hex(), "Hex generator returned nil")
			assert.NotNil(t, SuffixGenerators.Timestamp(), "Timestamp generator returned nil")
			assert.NotNil(t, SuffixGenerators.Letter(), "Letter generator returned nil")
		}
	})
}

func TestCalculateCombinations(t *testing.T) {
	t.Run("should calculate combinations for 1 component", func(t *testing.T) {
		combinations := CalculateCombinations(1, 1)
		expected := len(Adjectives)
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should calculate combinations for 2 components", func(t *testing.T) {
		combinations := CalculateCombinations(2, 1)
		expected := len(Adjectives) * len(Nouns)
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should calculate combinations for 3 components", func(t *testing.T) {
		combinations := CalculateCombinations(3, 1)
		expected := len(Adjectives) * len(Nouns) * len(Verbs)
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should calculate combinations for 4 components", func(t *testing.T) {
		combinations := CalculateCombinations(4, 1)
		expected := len(Adjectives) * len(Nouns) * len(Verbs) * len(Adverbs)
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should calculate combinations for 5 components", func(t *testing.T) {
		combinations := CalculateCombinations(5, 1)
		expected := len(Adjectives) * len(Nouns) * len(Verbs) * len(Adverbs) * len(Prepositions)
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should apply suffix multiplier", func(t *testing.T) {
		combinations := CalculateCombinations(2, 1000)
		expected := len(Adjectives) * len(Nouns) * 1000
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})
}

func TestCalculateCollisionProbability(t *testing.T) {
	t.Run("should return 0 for 1 or fewer IDs", func(t *testing.T) {
		assert.Equal(t, 0.0, CalculateCollisionProbability(1000, 0), "Expected 0 for 0 IDs")
		assert.Equal(t, 0.0, CalculateCollisionProbability(1000, 1), "Expected 0 for 1 ID")
		assert.Equal(t, 0.0, CalculateCollisionProbability(1000, -1), "Expected 0 for -1 IDs")
	})

	t.Run("should return 1 when IDs >= total combinations", func(t *testing.T) {
		assert.Equal(t, 1.0, CalculateCollisionProbability(100, 100), "Expected 1 for equal IDs and combinations")
		assert.Equal(t, 1.0, CalculateCollisionProbability(100, 150), "Expected 1 for more IDs than combinations")
	})

	t.Run("should return probability between 0 and 1 for normal cases", func(t *testing.T) {
		totalCombinations := len(Adjectives) * len(Nouns)
		probability := CalculateCollisionProbability(totalCombinations, 100)

		assert.GreaterOrEqual(t, probability, 0.0, "Probability should be >= 0")
		assert.LessOrEqual(t, probability, 1.0, "Probability should be <= 1")
		assert.Greater(t, probability, 0.0, "Probability should be > 0 for this case")
	})

	t.Run("should increase probability with more IDs", func(t *testing.T) {
		totalCombinations := len(Adjectives) * len(Nouns)
		prob1 := CalculateCollisionProbability(totalCombinations, 50)
		prob2 := CalculateCollisionProbability(totalCombinations, 100)
		prob3 := CalculateCollisionProbability(totalCombinations, 200)

		assert.Less(t, prob1, prob2, "Expected prob1 < prob2")
		assert.Less(t, prob2, prob3, "Expected prob2 < prob3")
	})

	t.Run("should handle edge case with very small total combinations", func(t *testing.T) {
		probability := CalculateCollisionProbability(2, 2)
		assert.Equal(t, 1.0, probability, "Expected 1 for edge case")
	})
}

func TestGetCollisionAnalysis(t *testing.T) {
	t.Run("should return analysis with total combinations", func(t *testing.T) {
		analysis := GetCollisionAnalysis(2, 1)
		expected := len(Adjectives) * len(Nouns)

		assert.Equal(t, expected, analysis.TotalCombinations, "Expected %d total combinations", expected)
	})

	t.Run("should return scenarios array", func(t *testing.T) {
		analysis := GetCollisionAnalysis(2, 1)

		assert.NotEmpty(t, analysis.Scenarios, "Expected non-empty scenarios array")
	})

	t.Run("should have valid scenario structure", func(t *testing.T) {
		analysis := GetCollisionAnalysis(2, 1)
		if len(analysis.Scenarios) == 0 {
			t.Skip("No scenarios to test")
		}

		scenario := analysis.Scenarios[0]
		assert.Greater(t, scenario.IDs, 0, "Expected positive IDs")
		assert.GreaterOrEqual(t, scenario.Probability, 0.0, "Expected probability >= 0")
		assert.LessOrEqual(t, scenario.Probability, 1.0, "Expected probability <= 1")
		assert.True(t, strings.HasSuffix(scenario.Percentage, "%"), "Expected percentage to end with %%")
	})

	t.Run("should filter out unrealistic scenarios", func(t *testing.T) {
		analysis := GetCollisionAnalysis(2, 1)

		// All scenarios should be less than 80% of total combinations
		threshold := float64(analysis.TotalCombinations) * 0.8
		for _, scenario := range analysis.Scenarios {
			assert.Less(t, float64(scenario.IDs), threshold, "Scenario with %d IDs should be filtered out", scenario.IDs)
		}
	})

	t.Run("should handle suffix range", func(t *testing.T) {
		analysis := GetCollisionAnalysis(2, 1000)
		expected := len(Adjectives) * len(Nouns) * 1000

		assert.Equal(t, expected, analysis.TotalCombinations, "Expected %d total combinations with suffix", expected)
	})

	t.Run("should handle all component counts", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			analysis := GetCollisionAnalysis(i, 1)
			assert.Greater(t, analysis.TotalCombinations, 0, "Expected positive total combinations for %d components", i)
		}
	})
}

func TestIntegration(t *testing.T) {
	digit3Regex := regexp.MustCompile(`^\d{3}$`)
	hexRegex := regexp.MustCompile(`^[0-9a-f]{2}$`)

	t.Run("should generate and parse ID correctly", func(t *testing.T) {
		id, err := Generate(GenerateOptions{Components: 3})
		require.NoError(t, err, "Generate should not fail")

		parsed := Parse(id, "-")
		assert.Len(t, parsed.Components, 3, "Expected 3 components")
		assert.Nil(t, parsed.Suffix, "Expected nil suffix")
	})

	t.Run("should generate and parse ID with suffix correctly", func(t *testing.T) {
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     SuffixGenerators.Number,
		})
		require.NoError(t, err, "Generate should not fail")

		parsed := Parse(id, "-")
		assert.Len(t, parsed.Components, 2, "Expected 2 components")
		require.NotNil(t, parsed.Suffix, "Expected non-nil suffix")
		assert.True(t, digit3Regex.MatchString(*parsed.Suffix), "Expected 3-digit suffix, got '%s'", *parsed.Suffix)
	})

	t.Run("should maintain consistency across multiple generations", func(t *testing.T) {
		options := GenerateOptions{
			Components: 3,
			Suffix:     SuffixGenerators.Hex,
			Separator:  "_",
		}

		for i := 0; i < 10; i++ {
			id, err := Generate(options)
			require.NoError(t, err, "Generate should not fail")

			parts := strings.Split(id, "_")
			assert.Len(t, parts, 4, "Expected 4 parts") // 3 components + 1 suffix

			assert.True(t, hexRegex.MatchString(parts[3]), "Expected hex suffix, got '%s'", parts[3])
		}
	})

	t.Run("should work with all suffix generators", func(t *testing.T) {
		generators := []func() *string{
			SuffixGenerators.Number,
			SuffixGenerators.Number4,
			SuffixGenerators.Hex,
			SuffixGenerators.Timestamp,
			SuffixGenerators.Letter,
		}

		for _, generator := range generators {
			id, err := Generate(GenerateOptions{
				Components: 2,
				Suffix:     generator,
			})
			require.NoError(t, err, "Generate should not fail")

			parts := strings.Split(id, "-")
			assert.Len(t, parts, 3, "Expected 3 parts") // 2 components + 1 suffix
		}
	})

	t.Run("should handle round trip with all component counts", func(t *testing.T) {
		for components := 1; components <= 5; components++ {
			id, err := Generate(GenerateOptions{Components: components})
			require.NoError(t, err, "Generate should not fail")

			parsed := Parse(id, "-")
			assert.Len(t, parsed.Components, components, "Expected %d components", components)
			assert.Nil(t, parsed.Suffix, "Expected nil suffix")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("should handle empty options object", func(t *testing.T) {
		id, err := Generate(GenerateOptions{})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 2, "Expected 2 parts (default)")
	})

	t.Run("should handle custom suffix returning empty string", func(t *testing.T) {
		emptySuffix := func() *string {
			s := ""
			return &s
		}
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     emptySuffix,
		})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 3, "Expected 3 parts") // empty string is still added
		assert.Equal(t, "", parts[2], "Expected empty suffix")
	})

	t.Run("should handle custom suffix returning whitespace", func(t *testing.T) {
		whitespaceSuffix := func() *string {
			s := "   "
			return &s
		}
		id, err := Generate(GenerateOptions{
			Components: 2,
			Suffix:     whitespaceSuffix,
		})
		require.NoError(t, err, "Generate should not fail")

		parts := strings.Split(id, "-")
		assert.Len(t, parts, 3, "Expected 3 parts")
		assert.Equal(t, "   ", parts[2], "Expected whitespace suffix")
	})

	t.Run("should parse empty string gracefully", func(t *testing.T) {
		result := Parse("", "-")

		expectedComponents := []string{""}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		assert.Nil(t, result.Suffix, "Expected nil suffix")
	})

	t.Run("should handle very large suffix ranges", func(t *testing.T) {
		combinations := CalculateCombinations(1, 1000000)
		expected := len(Adjectives) * 1000000
		assert.Equal(t, expected, combinations, "Expected %d combinations", expected)
	})

	t.Run("should handle parsing with different separators", func(t *testing.T) {
		separators := []string{"_", ".", "|", ":"}

		for _, sep := range separators {
			testID := "word1" + sep + "word2" + sep + "123"
			result := Parse(testID, sep)

			expectedComponents := []string{"word1", "word2"}
			assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
			require.NotNil(t, result.Suffix, "Expected non-nil suffix")
			assert.Equal(t, "123", *result.Suffix, "Expected suffix '123'")
		}
	})

	t.Run("should handle parsing IDs with no separators", func(t *testing.T) {
		result := Parse("singleword", "-")

		expectedComponents := []string{"singleword"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		assert.Nil(t, result.Suffix, "Expected nil suffix")
	})

	t.Run("should handle parsing numeric-only IDs", func(t *testing.T) {
		result := Parse("123-456-789", "-")

		expectedComponents := []string{"123", "456"}
		assert.True(t, sliceEqual(result.Components, expectedComponents), "Expected components %v, got %v", expectedComponents, result.Components)
		require.NotNil(t, result.Suffix, "Expected non-nil suffix")
		assert.Equal(t, "789", *result.Suffix, "Expected suffix '789'")
	})

	t.Run("should handle extreme collision probability scenarios", func(t *testing.T) {
		// Test with very large numbers
		prob1 := CalculateCollisionProbability(1000000, 1000)
		assert.GreaterOrEqual(t, prob1, 0.0, "Expected probability >= 0")
		assert.LessOrEqual(t, prob1, 1.0, "Expected probability <= 1")

		// Test with equal numbers
		prob2 := CalculateCollisionProbability(100, 100)
		assert.Equal(t, 1.0, prob2, "Expected probability 1")

		// Test with very small combinations
		prob3 := CalculateCollisionProbability(1, 2)
		assert.Equal(t, 1.0, prob3, "Expected probability 1")
	})

	t.Run("should handle boundary values for calculateCombinations", func(t *testing.T) {
		// Test with minimum values
		result := CalculateCombinations(1, 1)
		expected := len(Adjectives)
		assert.Equal(t, expected, result, "Expected %d combinations", expected)

		// Test with maximum components
		maxCombinations := CalculateCombinations(5, 1)
		assert.Greater(t, maxCombinations, 0, "Expected positive combinations")

		// Test with large suffix range
		largeCombinations := CalculateCombinations(1, 999999)
		expected = len(Adjectives) * 999999
		assert.Equal(t, expected, largeCombinations, "Expected %d combinations", expected)
	})
}

func TestDictionary(t *testing.T) {
	t.Run("should have correct dictionary stats", func(t *testing.T) {
		stats := GetDictionaryStats()

		assert.Equal(t, len(Adjectives), stats.Adjectives, "Expected %d adjectives", len(Adjectives))
		assert.Equal(t, len(Nouns), stats.Nouns, "Expected %d nouns", len(Nouns))
		assert.Equal(t, len(Verbs), stats.Verbs, "Expected %d verbs", len(Verbs))
		assert.Equal(t, len(Adverbs), stats.Adverbs, "Expected %d adverbs", len(Adverbs))
		assert.Equal(t, len(Prepositions), stats.Prepositions, "Expected %d prepositions", len(Prepositions))
	})

	t.Run("should return complete dictionary", func(t *testing.T) {
		dict := GetDictionary()

		assert.True(t, sliceEqual(dict.Adjectives, Adjectives), "Dictionary adjectives don't match")
		assert.True(t, sliceEqual(dict.Nouns, Nouns), "Dictionary nouns don't match")
		assert.True(t, sliceEqual(dict.Verbs, Verbs), "Dictionary verbs don't match")
		assert.True(t, sliceEqual(dict.Adverbs, Adverbs), "Dictionary adverbs don't match")
		assert.True(t, sliceEqual(dict.Prepositions, Prepositions), "Dictionary prepositions don't match")
	})

	t.Run("should validate all component ranges work correctly", func(t *testing.T) {
		// Test that each component position uses correct dictionary
		id1, err := Generate(GenerateOptions{Components: 1})
		require.NoError(t, err, "Generate should not fail")
		parts1 := strings.Split(id1, "-")
		assert.True(t, contains(Adjectives, parts1[0]), "First component '%s' not found in adjectives", parts1[0])

		id2, err := Generate(GenerateOptions{Components: 2})
		require.NoError(t, err, "Generate should not fail")
		parts2 := strings.Split(id2, "-")
		assert.True(t, contains(Adjectives, parts2[0]), "First component '%s' not found in adjectives", parts2[0])
		assert.True(t, contains(Nouns, parts2[1]), "Second component '%s' not found in nouns", parts2[1])

		id3, err := Generate(GenerateOptions{Components: 3})
		require.NoError(t, err, "Generate should not fail")
		parts3 := strings.Split(id3, "-")
		assert.True(t, contains(Adjectives, parts3[0]), "First component '%s' not found in adjectives", parts3[0])
		assert.True(t, contains(Nouns, parts3[1]), "Second component '%s' not found in nouns", parts3[1])
		assert.True(t, contains(Verbs, parts3[2]), "Third component '%s' not found in verbs", parts3[2])

		id4, err := Generate(GenerateOptions{Components: 4})
		require.NoError(t, err, "Generate should not fail")
		parts4 := strings.Split(id4, "-")
		assert.True(t, contains(Adjectives, parts4[0]), "First component '%s' not found in adjectives", parts4[0])
		assert.True(t, contains(Nouns, parts4[1]), "Second component '%s' not found in nouns", parts4[1])
		assert.True(t, contains(Verbs, parts4[2]), "Third component '%s' not found in verbs", parts4[2])
		assert.True(t, contains(Adverbs, parts4[3]), "Fourth component '%s' not found in adverbs", parts4[3])

		id5, err := Generate(GenerateOptions{Components: 5})
		require.NoError(t, err, "Generate should not fail")
		parts5 := strings.Split(id5, "-")
		assert.True(t, contains(Adjectives, parts5[0]), "First component '%s' not found in adjectives", parts5[0])
		assert.True(t, contains(Nouns, parts5[1]), "Second component '%s' not found in nouns", parts5[1])
		assert.True(t, contains(Verbs, parts5[2]), "Third component '%s' not found in verbs", parts5[2])
		assert.True(t, contains(Adverbs, parts5[3]), "Fourth component '%s' not found in adverbs", parts5[3])
		assert.True(t, contains(Prepositions, parts5[4]), "Fifth component '%s' not found in prepositions", parts5[4])
	})
}

// Helper functions

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// sliceEqual checks if two string slices are equal
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
