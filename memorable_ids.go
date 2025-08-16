package memorable_ids

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * Memorable ID Generator
 *
 * A flexible library for generating human-readable, memorable identifiers.
 * Uses combinations of adjectives, nouns, verbs, adverbs, and prepositions
 * with optional numeric/custom suffixes.
 *
 * @author Aris Ripandi
 * @license MIT
 */

// SuffixGenerator is a function type for generating suffixes
type SuffixGenerator func() *string

// GenerateOptions contains configuration options for ID generation
type GenerateOptions struct {
	// Components is the number of word components (1-5, default: 2)
	Components int
	// Suffix is the suffix generator function (default: nil)
	Suffix SuffixGenerator
	// Separator between parts (default: "-")
	Separator string
}

// ParsedID represents parsed ID components structure
type ParsedID struct {
	// Components is the array of word components
	Components []string
	// Suffix is the suffix part if detected, nil otherwise
	Suffix *string
}

// CollisionScenario represents collision scenario analysis
type CollisionScenario struct {
	// IDs is the number of IDs in scenario
	IDs int
	// Probability is the collision probability (0-1)
	Probability float64
	// Percentage is the formatted percentage string
	Percentage string
}

// CollisionAnalysis represents collision analysis result
type CollisionAnalysis struct {
	// TotalCombinations is the total possible combinations
	TotalCombinations int
	// Scenarios is the array of collision scenarios
	Scenarios []CollisionScenario
}

// Generate creates a memorable ID
//
// Example usage:
//
//	// Default: 2 components, no suffix
//	Generate(GenerateOptions{}) // "cute-rabbit"
//
//	// 3 components
//	Generate(GenerateOptions{Components: 3}) // "large-fox-swim"
//
//	// With numeric suffix
//	Generate(GenerateOptions{
//	  Components: 2,
//	  Suffix: SuffixGenerators.Number,
//	}) // "quick-mouse-042"
//
//	// Custom separator
//	Generate(GenerateOptions{
//	  Components: 2,
//	  Separator: "_",
//	}) // "warm_duck"
func Generate(options GenerateOptions) (string, error) {
	// Set defaults
	if options.Components == 0 {
		options.Components = 2
	}
	if options.Separator == "" {
		options.Separator = "-"
	}

	// Validate components range (after setting defaults)
	if options.Components < 1 || options.Components > 5 {
		return "", errors.New("components must be between 1 and 5")
	}

	var parts []string
	componentGenerators := []func() string{
		func() string { return randomItem(Adjectives) },   // 0: adjective
		func() string { return randomItem(Nouns) },        // 1: noun
		func() string { return randomItem(Verbs) },        // 2: verb
		func() string { return randomItem(Adverbs) },      // 3: adverb
		func() string { return randomItem(Prepositions) }, // 4: preposition
	}

	// Generate requested number of components
	for i := 0; i < options.Components; i++ {
		parts = append(parts, componentGenerators[i]())
	}

	// Add suffix if provided
	if options.Suffix != nil {
		suffixValue := options.Suffix()
		if suffixValue != nil {
			parts = append(parts, *suffixValue)
		}
	}

	return strings.Join(parts, options.Separator), nil
}

// randomItem returns a random item from a string slice
func randomItem(items []string) string {
	return items[rand.Intn(len(items))]
}

// DefaultSuffix generates a random 3-digit number suffix
//
// Example:
//
//	DefaultSuffix() // "042"
//	DefaultSuffix() // "789"
func DefaultSuffix() *string {
	suffix := fmt.Sprintf("%03d", rand.Intn(1000))
	return &suffix
}

// Parse parses a memorable ID back to its components
//
// Example:
//
//	Parse("cute-rabbit-042", "-")
//	// ParsedID{Components: ["cute", "rabbit"], Suffix: "042"}
//
//	Parse("large-fox-swim", "-")
//	// ParsedID{Components: ["large", "fox", "swim"], Suffix: nil}
func Parse(id string, separator string) ParsedID {
	if separator == "" {
		separator = "-"
	}

	parts := strings.Split(id, separator)
	result := ParsedID{
		Components: make([]string, 0),
		Suffix:     nil,
	}

	// Last part is likely suffix if it's numeric
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		matched, _ := regexp.MatchString(`^\d+$`, lastPart)
		if matched {
			result.Suffix = &lastPart
			result.Components = parts[:len(parts)-1]
		} else {
			result.Components = parts
		}
	}

	return result
}

// CalculateCombinations calculates total possible combinations for given configuration
//
// Example:
//
//	CalculateCombinations(2, 1)    // 5,304 (2 components, no suffix)
//	CalculateCombinations(2, 1000) // 5,304,000 (2 components + 3-digit suffix)
//	CalculateCombinations(3, 1)    // 212,160 (3 components, no suffix)
func CalculateCombinations(components int, suffixRange int) int {
	if components < 1 || components > 5 {
		return 0
	}
	if suffixRange < 1 {
		suffixRange = 1
	}

	stats := GetDictionaryStats()
	componentSizes := []int{
		stats.Adjectives,   // 78 adjectives
		stats.Nouns,        // 68 nouns
		stats.Verbs,        // 40 verbs
		stats.Adverbs,      // 27 adverbs
		stats.Prepositions, // 26 prepositions
	}

	total := 1
	for i := 0; i < components; i++ {
		total *= componentSizes[i]
	}

	return total * suffixRange
}

// CalculateCollisionProbability calculates collision probability using Birthday Paradox
//
// Example:
//
//	// For 2 components (5,304 total), generating 100 IDs
//	CalculateCollisionProbability(5304, 100) // ~0.0093 (0.93%)
//
//	// For 3 components (212,160 total), generating 10,000 IDs
//	CalculateCollisionProbability(212160, 10000) // ~0.00235 (0.235%)
func CalculateCollisionProbability(totalCombinations int, generatedIDs int) float64 {
	if generatedIDs >= totalCombinations {
		return 1.0
	}
	if generatedIDs <= 1 {
		return 0.0
	}

	// Birthday paradox approximation: 1 - e^(-n²/2N)
	exponent := -float64(generatedIDs*generatedIDs) / (2.0 * float64(totalCombinations))
	return 1.0 - math.Exp(exponent)
}

// GetCollisionAnalysis gets collision analysis for different ID generation scenarios
//
// Example:
//
//	GetCollisionAnalysis(2, 1)
//	// CollisionAnalysis{
//	//   TotalCombinations: 5304,
//	//   Scenarios: [
//	//     {IDs: 100, Probability: 0.0093, Percentage: "0.93%"},
//	//     {IDs: 500, Probability: 0.218, Percentage: "21.8%"},
//	//     ...
//	//   ]
//	// }
func GetCollisionAnalysis(components int, suffixRange int) CollisionAnalysis {
	if suffixRange < 1 {
		suffixRange = 1
	}

	total := CalculateCombinations(components, suffixRange)
	testSizes := []int{50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000}

	var scenarios []CollisionScenario
	threshold := int(float64(total) * 0.8) // Only show realistic scenarios

	for _, size := range testSizes {
		if size < threshold {
			probability := CalculateCollisionProbability(total, size)
			scenarios = append(scenarios, CollisionScenario{
				IDs:         size,
				Probability: probability,
				Percentage:  fmt.Sprintf("%.2f%%", probability*100),
			})
		}
	}

	return CollisionAnalysis{
		TotalCombinations: total,
		Scenarios:         scenarios,
	}
}

// SuffixGeneratorCollection contains predefined suffix generators
type SuffixGeneratorCollection struct {
	// Number generates random 3-digit number (000-999)
	// Adds 1,000x multiplier to total combinations
	Number func() *string

	// Number4 generates random 4-digit number (0000-9999)
	// Adds 10,000x multiplier to total combinations
	Number4 func() *string

	// Hex generates random 2-digit hex (00-ff)
	// Adds 256x multiplier to total combinations
	Hex func() *string

	// Timestamp generates last 4 digits of current timestamp
	// Adds ~10,000x multiplier (time-based, not truly random)
	Timestamp func() *string

	// Letter generates random lowercase letter (a-z)
	// Adds 26x multiplier to total combinations
	Letter func() *string
}

// SuffixGenerators contains collection of predefined suffix generators
var SuffixGenerators = SuffixGeneratorCollection{
	Number: DefaultSuffix,

	Number4: func() *string {
		suffix := fmt.Sprintf("%04d", rand.Intn(10000))
		return &suffix
	},

	Hex: func() *string {
		suffix := fmt.Sprintf("%02x", rand.Intn(256))
		return &suffix
	},

	Timestamp: func() *string {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		if len(timestamp) >= 4 {
			suffix := timestamp[len(timestamp)-4:]
			return &suffix
		}
		suffix := fmt.Sprintf("%04d", rand.Intn(10000))
		return &suffix
	},

	Letter: func() *string {
		suffix := string(rune('a' + rand.Intn(26)))
		return &suffix
	},
}
