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

import { adjectives, adverbs, dictionaryStats, nouns, prepositions, verbs } from './dictionary'

/**
 * Type definition for suffix generator function
 */
export type SuffixGenerator = () => string | null | undefined

/**
 * Configuration options for ID generation
 */
export interface GenerateOptions {
  /** Number of word components (1-5, default: 2) */
  components?: number
  /** Suffix generator function (default: null) */
  suffix?: SuffixGenerator | null
  /** Separator between parts (default: "-") */
  separator?: string
}

/**
 * Parsed ID components structure
 */
export interface ParsedId {
  /** Array of word components */
  components: string[]
  /** Suffix part if detected, null otherwise */
  suffix: string | null
}

/**
 * Collision scenario analysis
 */
export interface CollisionScenario {
  /** Number of IDs in scenario */
  ids: number
  /** Collision probability (0-1) */
  probability: number
  /** Formatted percentage string */
  percentage: string
}

/**
 * Collision analysis result
 */
export interface CollisionAnalysis {
  /** Total possible combinations */
  totalCombinations: number
  /** Array of collision scenarios */
  scenarios: CollisionScenario[]
}

/**
 * Generate a memorable ID
 *
 * @param options - Configuration options
 * @returns Generated memorable ID
 *
 * @example
 * ```typescript
 * // Default: 2 components, no suffix
 * generate() // "cute-rabbit"
 *
 * // 3 components
 * generate({ components: 3 }) // "large-fox-swim"
 *
 * // With numeric suffix
 * generate({
 *   components: 2,
 *   suffix: suffixGenerators.number
 * }) // "quick-mouse-042"
 *
 * // Custom separator
 * generate({
 *   components: 2,
 *   separator: "_"
 * }) // "warm_duck"
 * ```
 */
export function generate(options: GenerateOptions = {}): string {
  const { components = 2, suffix = null, separator = '-' } = options

  if (components < 1 || components > 5) {
    throw new Error('Components must be between 1 and 5')
  }

  function random(max: number): number {
    return Math.floor(Math.random() * max)
  }

  function randomItem<T>(array: readonly T[]): T {
    return array[random(array.length)]
  }

  const parts: string[] = []
  const componentGenerators: Array<() => string> = [
    () => randomItem(adjectives), // 0: adjective
    () => randomItem(nouns), // 1: noun
    () => randomItem(verbs), // 2: verb
    () => randomItem(adverbs), // 3: adverb
    () => randomItem(prepositions), // 4: preposition
  ]

  // Generate requested number of components
  for (let i = 0; i < components; i++) {
    parts.push(componentGenerators[i]())
  }

  // Add suffix if provided
  if (suffix && typeof suffix === 'function') {
    const suffixValue = suffix()
    if (suffixValue !== null && suffixValue !== undefined) {
      parts.push(suffixValue)
    }
  }

  return parts.join(separator)
}

/**
 * Default suffix generator - random 3-digit number
 *
 * @returns Random number suffix (000-999)
 *
 * @example
 * ```typescript
 * defaultSuffix() // "042"
 * defaultSuffix() // "789"
 * ```
 */
export function defaultSuffix(): string {
  return Math.floor(Math.random() * 1000)
    .toString()
    .padStart(3, '0')
}

/**
 * Parse a memorable ID back to its components
 *
 * @param id - The memorable ID to parse
 * @param separator - Separator used (default: "-")
 * @returns Parsed components with structure
 *
 * @example
 * ```typescript
 * parse("cute-rabbit-042")
 * // { components: ["cute", "rabbit"], suffix: "042" }
 *
 * parse("large-fox-swim")
 * // { components: ["large", "fox", "swim"], suffix: null }
 * ```
 */
export function parse(id: string, separator = '-'): ParsedId {
  const parts = id.split(separator)
  const result: ParsedId = {
    components: [],
    suffix: null,
  }

  // Last part is likely suffix if it's numeric
  const lastPart = parts[parts.length - 1]
  if (/^\d+$/.test(lastPart)) {
    result.suffix = lastPart
    result.components = parts.slice(0, -1)
  } else {
    result.components = parts
  }

  return result
}

/**
 * Calculate total possible combinations for given configuration
 *
 * @param components - Number of word components (1-5)
 * @param suffixRange - Range of suffix values (default: 1 for no suffix)
 * @returns Total possible unique combinations
 *
 * @example
 * ```typescript
 * calculateCombinations(2) // 5,304 (2 components, no suffix)
 * calculateCombinations(2, 1000) // 5,304,000 (2 components + 3-digit suffix)
 * calculateCombinations(3) // 212,160 (3 components, no suffix)
 * ```
 */
export function calculateCombinations(components = 2, suffixRange = 1): number {
  const componentSizes = [
    dictionaryStats.adjectives, // 78 adjectives
    dictionaryStats.nouns, // 68 nouns
    dictionaryStats.verbs, // 40 verbs
    dictionaryStats.adverbs, // 27 adverbs
    dictionaryStats.prepositions, // 26 prepositions
  ]

  let total = 1
  for (let i = 0; i < components; i++) {
    total *= componentSizes[i]
  }

  return total * suffixRange
}

/**
 * Calculate collision probability using Birthday Paradox
 *
 * @param totalCombinations - Total possible combinations
 * @param generatedIds - Number of IDs to generate
 * @returns Collision probability (0-1)
 *
 * @example
 * ```typescript
 * // For 2 components (5,304 total), generating 100 IDs
 * calculateCollisionProbability(5304, 100) // ~0.0093 (0.93%)
 *
 * // For 3 components (212,160 total), generating 10,000 IDs
 * calculateCollisionProbability(212160, 10000) // ~0.00235 (0.235%)
 * ```
 */
export function calculateCollisionProbability(
  totalCombinations: number,
  generatedIds: number
): number {
  if (generatedIds >= totalCombinations) return 1
  if (generatedIds <= 1) return 0

  // Birthday paradox approximation: 1 - e^(-n²/2N)
  const exponent = -(generatedIds * generatedIds) / (2 * totalCombinations)
  return 1 - Math.exp(exponent)
}

/**
 * Get collision analysis for different ID generation scenarios
 *
 * @param components - Number of components
 * @param suffixRange - Suffix range (1 for no suffix)
 * @returns Analysis with total combinations and collision probabilities
 *
 * @example
 * ```typescript
 * getCollisionAnalysis(2)
 * // {
 * //   totalCombinations: 5304,
 * //   scenarios: [
 * //     { ids: 100, probability: 0.0093, percentage: "0.93%" },
 * //     { ids: 500, probability: 0.218, percentage: "21.8%" },
 * //     ...
 * //   ]
 * // }
 * ```
 */
export function getCollisionAnalysis(components = 2, suffixRange = 1): CollisionAnalysis {
  const total = calculateCombinations(components, suffixRange)
  const testSizes = [50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000]

  return {
    totalCombinations: total,
    scenarios: testSizes
      .filter((size) => size < total * 0.8) // Only show realistic scenarios
      .map((size) => ({
        ids: size,
        probability: calculateCollisionProbability(total, size),
        percentage: `${(calculateCollisionProbability(total, size) * 100).toFixed(2)}%`,
      })),
  }
}

/**
 * Collection of predefined suffix generators
 */
export const suffixGenerators = {
  /**
   * Random 3-digit number (000-999)
   * Adds 1,000x multiplier to total combinations
   */
  number: defaultSuffix,

  /**
   * Random 4-digit number (0000-9999)
   * Adds 10,000x multiplier to total combinations
   */
  number4: (): string =>
    Math.floor(Math.random() * 10000)
      .toString()
      .padStart(4, '0'),

  /**
   * Random 2-digit hex (00-ff)
   * Adds 256x multiplier to total combinations
   */
  hex: (): string =>
    Math.floor(Math.random() * 256)
      .toString(16)
      .padStart(2, '0'),

  /**
   * Last 4 digits of current timestamp
   * Adds ~10,000x multiplier (time-based, not truly random)
   */
  timestamp: (): string => Date.now().toString().slice(-4),

  /**
   * Random lowercase letter (a-z)
   * Adds 26x multiplier to total combinations
   */
  letter: (): string => String.fromCharCode(97 + Math.floor(Math.random() * 26)),
} as const

/**
 * Default export with all main functions
 */
const memorableId = {
  generate,
  parse,
  calculateCombinations,
  calculateCollisionProbability,
  getCollisionAnalysis,
  suffixGenerators,
  defaultSuffix,
} as const

export default memorableId

/**
 * Re-export dictionary for external use
 */
export { dictionary, dictionaryStats } from './dictionary'
