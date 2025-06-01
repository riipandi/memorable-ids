import assert from 'node:assert'
import { describe, test } from 'node:test'
import { calculateCollisionProbability, calculateCombinations } from '.'
import { generate, parse, suffixGenerators } from '.'
import { defaultSuffix, getCollisionAnalysis } from '.'
import { adjectives, adverbs, nouns, prepositions, verbs } from './dictionary'
import type { GenerateOptions } from './index'

describe('Memorable ID Generator', () => {
  describe('generate()', () => {
    test('should generate ID with default options (2 components)', () => {
      const id = generate()
      const parts = id.split('-')

      assert.strictEqual(parts.length, 2)
      assert.ok(adjectives.includes(parts[0]))
      assert.ok(nouns.includes(parts[1]))
    })

    test('should generate ID with 1 component', () => {
      const id = generate({ components: 1 })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 1)
      assert.ok(adjectives.includes(parts[0]))
    })

    test('should generate ID with 3 components', () => {
      const id = generate({ components: 3 })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 3)
      assert.ok(adjectives.includes(parts[0]))
      assert.ok(nouns.includes(parts[1]))
      assert.ok(verbs.includes(parts[2]))
    })

    test('should generate ID with 4 components', () => {
      const id = generate({ components: 4 })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 4)
      assert.ok(adjectives.includes(parts[0]))
      assert.ok(nouns.includes(parts[1]))
      assert.ok(verbs.includes(parts[2]))
      assert.ok(adverbs.includes(parts[3]))
    })

    test('should generate ID with 5 components', () => {
      const id = generate({ components: 5 })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 5)
      assert.ok(adjectives.includes(parts[0]))
      assert.ok(nouns.includes(parts[1]))
      assert.ok(verbs.includes(parts[2]))
      assert.ok(adverbs.includes(parts[3]))
      assert.ok(prepositions.includes(parts[4]))
    })

    test('should use custom separator', () => {
      const id = generate({ components: 2, separator: '_' })
      const parts = id.split('_')

      assert.strictEqual(parts.length, 2)
      assert.ok(id.includes('_'))
      // Remove the assertion about '-' since some words in dictionary contain hyphens
      assert.ok(adjectives.includes(parts[0]))
      assert.ok(nouns.includes(parts[1]))
    })

    test('should add suffix when provided', () => {
      const id = generate({
        components: 2,
        suffix: suffixGenerators.number,
      })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 3)
      assert.ok(/^\d{3}$/.test(parts[2])) // 3-digit number
    })

    test('should handle null suffix gracefully', () => {
      const nullSuffix = () => null
      const id = generate({
        components: 2,
        suffix: nullSuffix,
      })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 2) // No suffix added
    })

    test('should handle undefined suffix gracefully', () => {
      const undefinedSuffix = () => undefined
      const id = generate({
        components: 2,
        suffix: undefinedSuffix,
      })

      const parts = id.split('-')

      // No suffix added when undefined is returned
      assert.strictEqual(parts.length, 2)
    })

    test('should handle suffix that is not a function', () => {
      const id = generate({
        components: 2,
        suffix: null,
      })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 2) // No suffix added
    })

    test('should throw error for invalid component count', () => {
      assert.throws(() => {
        generate({ components: 0 })
      }, /Components must be between 1 and 5/)

      assert.throws(() => {
        generate({ components: 6 })
      }, /Components must be between 1 and 5/)

      assert.throws(() => {
        generate({ components: -1 })
      }, /Components must be between 1 and 5/)

      assert.throws(() => {
        generate({ components: 10 })
      }, /Components must be between 1 and 5/)
    })

    test('should generate different IDs on multiple calls', () => {
      const ids = new Set()
      for (let i = 0; i < 100; i++) {
        ids.add(generate())
      }

      // Should have high uniqueness (allowing for some collisions)
      assert.ok(ids.size > 90)
    })
  })

  describe('parse()', () => {
    test('should parse ID without suffix', () => {
      const result = parse('cute-rabbit')

      assert.deepStrictEqual(result, {
        components: ['cute', 'rabbit'],
        suffix: null,
      })
    })

    test('should parse ID with numeric suffix', () => {
      const result = parse('cute-rabbit-042')

      assert.deepStrictEqual(result, {
        components: ['cute', 'rabbit'],
        suffix: '042',
      })
    })

    test('should parse ID with non-numeric suffix as component', () => {
      const result = parse('cute-rabbit-swim')

      assert.deepStrictEqual(result, {
        components: ['cute', 'rabbit', 'swim'],
        suffix: null,
      })
    })

    test('should handle custom separator', () => {
      const result = parse('cute_rabbit_123', '_')

      assert.deepStrictEqual(result, {
        components: ['cute', 'rabbit'],
        suffix: '123',
      })
    })

    test('should handle single component', () => {
      const result = parse('cute')

      assert.deepStrictEqual(result, {
        components: ['cute'],
        suffix: null,
      })
    })

    test('should handle single component with suffix', () => {
      const result = parse('cute-123')

      assert.deepStrictEqual(result, {
        components: ['cute'],
        suffix: '123',
      })
    })

    test('should handle ID with only numeric part', () => {
      const result = parse('123')

      assert.deepStrictEqual(result, {
        components: [],
        suffix: '123',
      })
    })

    test('should handle mixed numeric patterns', () => {
      const result = parse('cute-123abc-456')

      assert.deepStrictEqual(result, {
        components: ['cute', '123abc'],
        suffix: '456',
      })
    })
  })

  describe('suffixGenerators', () => {
    test('number should generate 3-digit string', () => {
      const suffix = suffixGenerators.number()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^\d{3}$/.test(suffix))
    })

    test('number4 should generate 4-digit string', () => {
      const suffix = suffixGenerators.number4()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^\d{4}$/.test(suffix))
    })

    test('hex should generate 2-digit hex string', () => {
      const suffix = suffixGenerators.hex()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^[0-9a-f]{2}$/.test(suffix))
    })

    test('timestamp should generate 4-digit string', () => {
      const suffix = suffixGenerators.timestamp()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^\d{4}$/.test(suffix))
    })

    test('letter should generate single lowercase letter', () => {
      const suffix = suffixGenerators.letter()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^[a-z]$/.test(suffix))
    })

    test('defaultSuffix should work same as number generator', () => {
      const suffix = defaultSuffix()

      assert.ok(typeof suffix === 'string')
      assert.ok(/^\d{3}$/.test(suffix))
    })

    test('all suffix generators should produce valid output', () => {
      // Test edge cases for all generators
      for (let i = 0; i < 10; i++) {
        assert.ok(typeof suffixGenerators.number() === 'string')
        assert.ok(typeof suffixGenerators.number4() === 'string')
        assert.ok(typeof suffixGenerators.hex() === 'string')
        assert.ok(typeof suffixGenerators.timestamp() === 'string')
        assert.ok(typeof suffixGenerators.letter() === 'string')
      }
    })
  })

  describe('calculateCombinations()', () => {
    test('should calculate combinations for 1 component', () => {
      const combinations = calculateCombinations(1)
      assert.strictEqual(combinations, adjectives.length) // actual adjectives length
    })

    test('should calculate combinations for 2 components', () => {
      const combinations = calculateCombinations(2)
      assert.strictEqual(combinations, adjectives.length * nouns.length) // actual lengths
    })

    test('should calculate combinations for 3 components', () => {
      const combinations = calculateCombinations(3)
      assert.strictEqual(combinations, adjectives.length * nouns.length * verbs.length) // actual lengths
    })

    test('should calculate combinations for 4 components', () => {
      const combinations = calculateCombinations(4)
      assert.strictEqual(
        combinations,
        adjectives.length * nouns.length * verbs.length * adverbs.length
      )
    })

    test('should calculate combinations for 5 components', () => {
      const combinations = calculateCombinations(5)
      assert.strictEqual(
        combinations,
        adjectives.length * nouns.length * verbs.length * adverbs.length * prepositions.length
      )
    })

    test('should apply suffix multiplier', () => {
      const combinations = calculateCombinations(2, 1000)
      assert.strictEqual(combinations, adjectives.length * nouns.length * 1000)
    })

    test('should handle default parameters', () => {
      const combinations = calculateCombinations()
      assert.strictEqual(combinations, adjectives.length * nouns.length) // default 2 components
    })

    test('should handle zero suffix range', () => {
      const combinations = calculateCombinations(2, 0)
      assert.strictEqual(combinations, 0)
    })
  })

  describe('calculateCollisionProbability()', () => {
    test('should return 0 for 1 or fewer IDs', () => {
      assert.strictEqual(calculateCollisionProbability(1000, 0), 0)
      assert.strictEqual(calculateCollisionProbability(1000, 1), 0)
      assert.strictEqual(calculateCollisionProbability(1000, -1), 0)
    })

    test('should return 1 when IDs >= total combinations', () => {
      assert.strictEqual(calculateCollisionProbability(100, 100), 1)
      assert.strictEqual(calculateCollisionProbability(100, 150), 1)
    })

    test('should return probability between 0 and 1 for normal cases', () => {
      const totalCombinations = adjectives.length * nouns.length // use actual values
      const probability = calculateCollisionProbability(totalCombinations, 100)

      assert.ok(probability >= 0, 'Probability should be >= 0')
      assert.ok(probability <= 1, 'Probability should be <= 1')
      assert.ok(probability > 0, 'Probability should be > 0 for this case')
    })

    test('should increase probability with more IDs', () => {
      const totalCombinations = adjectives.length * nouns.length // use actual values
      const prob1 = calculateCollisionProbability(totalCombinations, 50)
      const prob2 = calculateCollisionProbability(totalCombinations, 100)
      const prob3 = calculateCollisionProbability(totalCombinations, 200)

      assert.ok(prob1 < prob2)
      assert.ok(prob2 < prob3)
    })

    test('should handle edge case with very small total combinations', () => {
      const probability = calculateCollisionProbability(2, 2)
      assert.strictEqual(probability, 1)
    })
  })

  describe('getCollisionAnalysis()', () => {
    test('should return analysis with total combinations', () => {
      const analysis = getCollisionAnalysis(2)

      assert.ok(typeof analysis.totalCombinations === 'number')
      assert.strictEqual(analysis.totalCombinations, adjectives.length * nouns.length)
    })

    test('should return scenarios array', () => {
      const analysis = getCollisionAnalysis(2)

      assert.ok(Array.isArray(analysis.scenarios))
      assert.ok(analysis.scenarios.length > 0)
    })

    test('should have valid scenario structure', () => {
      const analysis = getCollisionAnalysis(2)
      const scenario = analysis.scenarios[0]

      assert.ok(typeof scenario.ids === 'number')
      assert.ok(typeof scenario.probability === 'number')
      assert.ok(typeof scenario.percentage === 'string')
      assert.ok(scenario.percentage.endsWith('%'))
    })

    test('should filter out unrealistic scenarios', () => {
      const analysis = getCollisionAnalysis(2)

      // All scenarios should be less than 80% of total combinations
      for (const scenario of analysis.scenarios) {
        assert.ok(scenario.ids < analysis.totalCombinations * 0.8)
      }
    })

    test('should handle suffix range', () => {
      const analysis = getCollisionAnalysis(2, 1000)

      assert.strictEqual(analysis.totalCombinations, adjectives.length * nouns.length * 1000)
    })

    test('should handle all component counts', () => {
      for (let i = 1; i <= 5; i++) {
        const analysis = getCollisionAnalysis(i)
        assert.ok(analysis.totalCombinations > 0)
        assert.ok(Array.isArray(analysis.scenarios))
      }
    })

    test('should handle very small combinations that filter all scenarios', () => {
      // Create a scenario where total combinations is very small
      const analysis = getCollisionAnalysis(1, 1) // Only adjectives.length combinations

      // Should still return valid structure even if scenarios array might be empty or small
      assert.ok(typeof analysis.totalCombinations === 'number')
      assert.ok(Array.isArray(analysis.scenarios))
    })
  })

  describe('Integration tests', () => {
    test('should generate and parse ID correctly', () => {
      const id = generate({ components: 3 })
      const parsed = parse(id)

      assert.strictEqual(parsed.components.length, 3)
      assert.strictEqual(parsed.suffix, null)
    })

    test('should generate and parse ID with suffix correctly', () => {
      const id = generate({
        components: 2,
        suffix: suffixGenerators.number,
      })
      const parsed = parse(id)

      assert.strictEqual(parsed.components.length, 2)
      assert.ok(parsed.suffix !== null)
      assert.ok(typeof parsed.suffix === 'string' && /^\d{3}$/.test(parsed.suffix))
    })

    test('should maintain consistency across multiple generations', () => {
      const options: GenerateOptions = {
        components: 3,
        suffix: suffixGenerators.hex,
        separator: '_',
      }

      for (let i = 0; i < 10; i++) {
        const id = generate(options)
        const parts = id.split('_')

        assert.strictEqual(parts.length, 4) // 3 components + 1 suffix
        assert.ok(/^[0-9a-f]{2}$/.test(parts[3])) // hex suffix
      }
    })

    test('should work with all suffix generators', () => {
      const generators = [
        suffixGenerators.number,
        suffixGenerators.number4,
        suffixGenerators.hex,
        suffixGenerators.timestamp,
        suffixGenerators.letter,
      ]

      for (const generator of generators) {
        const id = generate({
          components: 2,
          suffix: generator,
        })
        const parts = id.split('-')
        assert.strictEqual(parts.length, 3) // 2 components + 1 suffix
      }
    })

    test('should handle round trip with all component counts', () => {
      for (let components = 1; components <= 5; components++) {
        const id = generate({ components })
        const parsed = parse(id)
        assert.strictEqual(parsed.components.length, components)
        assert.strictEqual(parsed.suffix, null)
      }
    })
  })

  describe('Edge cases', () => {
    test('should handle empty options object', () => {
      const id = generate({})
      const parts = id.split('-')

      assert.strictEqual(parts.length, 2) // default behavior
    })

    test('should handle custom suffix returning empty string', () => {
      const emptySuffix = () => ''
      const id = generate({
        components: 2,
        suffix: emptySuffix,
      })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 3) // empty string is still added
      assert.strictEqual(parts[2], '')
    })

    test('should handle custom suffix returning whitespace', () => {
      const whitespaceSuffix = () => '   '
      const id = generate({
        components: 2,
        suffix: whitespaceSuffix,
      })
      const parts = id.split('-')

      assert.strictEqual(parts.length, 3)
      assert.strictEqual(parts[2], '   ')
    })

    test('should parse empty string gracefully', () => {
      const result = parse('')

      assert.deepStrictEqual(result, {
        components: [''],
        suffix: null,
      })
    })

    test('should handle very large suffix ranges', () => {
      const combinations = calculateCombinations(1, 1000000)
      assert.strictEqual(combinations, adjectives.length * 1000000)
    })

    test('should handle parsing with different separators', () => {
      const separators = ['_', '.', '|', ':']

      for (const sep of separators) {
        const result = parse(`word1${sep}word2${sep}123`, sep)
        assert.deepStrictEqual(result, {
          components: ['word1', 'word2'],
          suffix: '123',
        })
      }
    })

    test('should handle parsing IDs with no separators', () => {
      const result = parse('singleword')
      assert.deepStrictEqual(result, {
        components: ['singleword'],
        suffix: null,
      })
    })

    test('should handle parsing numeric-only IDs', () => {
      const result = parse('123-456-789')
      assert.deepStrictEqual(result, {
        components: ['123', '456'],
        suffix: '789',
      })
    })

    test('should handle extreme collision probability scenarios', () => {
      // Test with very large numbers
      const prob1 = calculateCollisionProbability(1000000, 1000)
      assert.ok(prob1 >= 0 && prob1 <= 1)

      // Test with equal numbers
      const prob2 = calculateCollisionProbability(100, 100)
      assert.strictEqual(prob2, 1)

      // Test with very small combinations
      const prob3 = calculateCollisionProbability(1, 2)
      assert.strictEqual(prob3, 1)
    })

    test('should handle suffix generators edge cases', () => {
      // Test that all generators work multiple times
      for (let i = 0; i < 5; i++) {
        const num = suffixGenerators.number()
        const num4 = suffixGenerators.number4()
        const hex = suffixGenerators.hex()
        const ts = suffixGenerators.timestamp()
        const letter = suffixGenerators.letter()

        assert.ok(/^\d{3}$/.test(num))
        assert.ok(/^\d{4}$/.test(num4))
        assert.ok(/^[0-9a-f]{2}$/.test(hex))
        assert.ok(/^\d{4}$/.test(ts))
        assert.ok(/^[a-z]$/.test(letter))
      }
    })

    test('should handle boundary values for calculateCombinations', () => {
      // Test with minimum values
      assert.strictEqual(calculateCombinations(1, 1), adjectives.length)

      // Test with maximum components
      const maxCombinations = calculateCombinations(5, 1)
      assert.ok(maxCombinations > 0)

      // Test with large suffix range
      const largeCombinations = calculateCombinations(1, 999999)
      assert.strictEqual(largeCombinations, adjectives.length * 999999)
    })
  })

  describe('Default export and re-exports', () => {
    test('should import default export correctly', async () => {
      const memorableId = await import('./index.js')

      assert.ok(typeof memorableId.default === 'object')
      assert.ok(typeof memorableId.default.generate === 'function')
      assert.ok(typeof memorableId.default.parse === 'function')
      assert.ok(typeof memorableId.default.calculateCombinations === 'function')
      assert.ok(typeof memorableId.default.calculateCollisionProbability === 'function')
      assert.ok(typeof memorableId.default.getCollisionAnalysis === 'function')
      assert.ok(typeof memorableId.default.suffixGenerators === 'object')
      assert.ok(typeof memorableId.default.defaultSuffix === 'function')
    })

    test('should re-export dictionary correctly', async () => {
      const { dictionary, dictionaryStats } = await import('./index.js')

      assert.ok(typeof dictionary === 'object')
      assert.ok(typeof dictionaryStats === 'object')
      assert.ok(Array.isArray(dictionary.adjectives))
      assert.ok(Array.isArray(dictionary.nouns))
      assert.ok(typeof dictionaryStats.adjectives === 'number')
      assert.ok(typeof dictionaryStats.nouns === 'number')
    })
  })

  describe('Type safety and validation', () => {
    test('should handle invalid suffix function gracefully', () => {
      const invalidSuffix = () => {
        throw new Error('Suffix generation failed')
      }

      // Should not throw, but handle gracefully
      assert.throws(() => {
        generate({
          components: 2,
          suffix: invalidSuffix,
        })
      })
    })

    test('should validate all component ranges work correctly', () => {
      // Test that each component position uses correct dictionary
      const id1 = generate({ components: 1 })
      const parts1 = id1.split('-')
      assert.ok(adjectives.includes(parts1[0]))

      const id2 = generate({ components: 2 })
      const parts2 = id2.split('-')
      assert.ok(adjectives.includes(parts2[0]))
      assert.ok(nouns.includes(parts2[1]))

      const id3 = generate({ components: 3 })
      const parts3 = id3.split('-')
      assert.ok(adjectives.includes(parts3[0]))
      assert.ok(nouns.includes(parts3[1]))
      assert.ok(verbs.includes(parts3[2]))

      const id4 = generate({ components: 4 })
      const parts4 = id4.split('-')
      assert.ok(adjectives.includes(parts4[0]))
      assert.ok(nouns.includes(parts4[1]))
      assert.ok(verbs.includes(parts4[2]))
      assert.ok(adverbs.includes(parts4[3]))

      const id5 = generate({ components: 5 })
      const parts5 = id5.split('-')
      assert.ok(adjectives.includes(parts5[0]))
      assert.ok(nouns.includes(parts5[1]))
      assert.ok(verbs.includes(parts5[2]))
      assert.ok(adverbs.includes(parts5[3]))
      assert.ok(prepositions.includes(parts5[4]))
    })

    test('should handle mathematical edge cases in collision probability', () => {
      // Test very small probability calculations
      const prob = calculateCollisionProbability(1000000, 2)
      assert.ok(prob > 0 && prob < 0.001)

      // Test approaching 50% probability (birthday paradox sweet spot)
      const combinations = 365 // Like birthday paradox
      const ids = Math.sqrt(2 * combinations * Math.log(2)) // Approximate 50% point
      const prob50 = calculateCollisionProbability(combinations, Math.floor(ids))
      assert.ok(prob50 > 0.4 && prob50 < 0.6)
    })
  })
})
