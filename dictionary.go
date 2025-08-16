package memorable_ids

/**
 * Dictionary of words for memorable ID generation
 *
 * Contains collections of English words categorized by part of speech.
 * Used to generate human-readable, memorable identifiers.
 *
 * @author Aris Ripandi
 * @license MIT
 */

// Adjectives contains English adjectives (78 total)
// Descriptive words that modify nouns
var Adjectives = []string{
	"cute", "dapper", "large", "small", "long", "short", "thick", "narrow",
	"deep", "flat", "whole", "low", "high", "near", "far", "fast", "quick",
	"slow", "early", "late", "bright", "dark", "cloudy", "warm", "cool",
	"cold", "windy", "noisy", "loud", "quiet", "dry", "clear", "hard",
	"soft", "heavy", "light", "strong", "weak", "tidy", "clean", "dirty",
	"empty", "full", "close", "thirsty", "hungry", "fat", "old", "fresh",
	"dead", "healthy", "sweet", "sour", "bitter", "salty", "good", "bad",
	"great", "important", "useful", "expensive", "cheap", "free", "difficult",
	"able", "rich", "afraid", "brave", "fine", "sad", "proud", "comfortable",
	"happy", "clever", "interesting", "famous", "exciting", "funny", "kind",
	"polite", "fair", "busy", "lazy", "lucky", "careful", "safe", "dangerous",
}

// Nouns contains English nouns - animals and common objects (68 total)
// Concrete things, animals, and objects
var Nouns = []string{
	"rabbit", "badger", "fox", "chicken", "bat", "deer", "snake", "hare",
	"hedgehog", "platypus", "mole", "mouse", "otter", "rat", "squirrel",
	"stoat", "weasel", "crow", "dove", "duck", "goose", "hawk", "heron",
	"kingfisher", "owl", "peacock", "pheasant", "pigeon", "robin", "rook",
	"sparrow", "starling", "swan", "ant", "bee", "butterfly", "dragonfly",
	"fly", "moth", "spider", "pike", "salmon", "trout", "frog", "newt",
	"toad", "crab", "lobster", "clam", "cockle", "mussel", "oyster", "snail",
	"cow", "dog", "donkey", "goat", "horse", "pig", "sheep", "ferret",
	"gerbil", "guinea-pig", "parrot", "book", "table", "chair", "lamp",
	"phone", "computer", "window", "door",
}

// Verbs contains English verbs - present tense (40 total)
// Action words in present tense form
var Verbs = []string{
	"sing", "play", "knit", "flounder", "dance", "listen", "run", "talk",
	"cuddle", "sit", "kiss", "hug", "whimper", "hide", "fight", "whisper",
	"cry", "snuggle", "walk", "drive", "loiter", "feel", "jump", "hop",
	"go", "marry", "engage", "sleep", "eat", "drink", "read", "write",
	"swim", "fly", "climb", "build", "create", "explore", "discover", "learn",
}

// Adverbs contains English adverbs (27 total)
// Words that modify verbs, adjectives, or other adverbs
var Adverbs = []string{
	"jovially", "merrily", "cordially", "carefully", "correctly", "eagerly",
	"easily", "fast", "loudly", "patiently", "quickly", "quietly", "slowly",
	"gently", "firmly", "softly", "boldly", "bravely", "calmly", "clearly",
	"closely", "deeply", "directly", "exactly", "fairly", "freely", "fully",
}

// Prepositions contains English prepositions (26 total)
// Words that show relationships between other words
var Prepositions = []string{
	"in", "on", "at", "by", "for", "with", "from", "to", "of", "about",
	"under", "over", "through", "between", "among", "during", "before",
	"after", "above", "below", "beside", "behind", "beyond", "within",
	"without", "across",
}

// DictionaryStats contains dictionary statistics for combination calculations
type DictionaryStats struct {
	Adjectives   int
	Nouns        int
	Verbs        int
	Adverbs      int
	Prepositions int
}

// GetDictionaryStats returns the statistics of all word collections
func GetDictionaryStats() DictionaryStats {
	return DictionaryStats{
		Adjectives:   len(Adjectives),
		Nouns:        len(Nouns),
		Verbs:        len(Verbs),
		Adverbs:      len(Adverbs),
		Prepositions: len(Prepositions),
	}
}

// Dictionary contains all word collections grouped by type
type Dictionary struct {
	Adjectives   []string
	Nouns        []string
	Verbs        []string
	Adverbs      []string
	Prepositions []string
	Stats        DictionaryStats
}

// GetDictionary returns the complete dictionary with all word collections
func GetDictionary() Dictionary {
	return Dictionary{
		Adjectives:   Adjectives,
		Nouns:        Nouns,
		Verbs:        Verbs,
		Adverbs:      Adverbs,
		Prepositions: Prepositions,
		Stats:        GetDictionaryStats(),
	}
}
