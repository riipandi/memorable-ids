# Memorable IDs

A flexible TypeScript library for generating human-readable, memorable identifiers.
Uses combinations of adjectives, nouns, verbs, adverbs, and prepositions with optional numeric/custom suffixes.

**Compatibility**: ESM-only.

## Features

- 🎯 **Human-readable** - Generate IDs like `cute-rabbit`, `quick-owl-dance-quietly`, etc
- 🔧 **Flexible** - 1-5 word components with customizable separators
- 📊 **Predictable** - Built-in collision analysis and capacity planning
- 🎲 **Extensible** - Custom suffix generators and vocabulary
- 📝 **TypeScript** - Full type safety and IntelliSense support
- ⚡ **Fast** - ~1M IDs per second generation speed
- 🪶 **Lightweight** - ~10KB vocabulary, zero dependencies

## Installation

```bash
npm install memorable-ids
```

## Quick Start

```typescript
import { generate, suffixGenerators } from 'memorable-ids';

// Basic usage - 2 components
generate(); // "cute-rabbit"

// More components for uniqueness
generate({ components: 3 }); // "large-fox-swim"

// Add numeric suffix for extra capacity
generate({ 
  components: 2, 
  suffix: suffixGenerators.number 
}); // "quick-mouse-042"

// Custom separator
generate({ 
  components: 2, 
  separator: "_" 
}); // "warm_duck"
```

## API Reference

### `generate(options?)`

Generate a memorable ID with customizable options.

**Parameters:**
- `options.components` (number, 1-5): Number of word components (default: 2)
- `options.suffix` (function): Suffix generator function (default: null)
- `options.separator` (string): Separator between parts (default: "-")

**Returns:** `string` - Generated memorable ID

**Examples:**

```typescript
// Different component counts
generate({ components: 1 }); // "bright"
generate({ components: 2 }); // "cute-rabbit" 
generate({ components: 3 }); // "large-fox-swim"
generate({ components: 4 }); // "happy-owl-dance-quietly"
generate({ components: 5 }); // "clever-fox-run-quickly-through"

// With suffixes
generate({ 
  components: 2, 
  suffix: suffixGenerators.number 
}); // "safe-rabbit-042"

generate({ 
  components: 2, 
  suffix: suffixGenerators.hex 
}); // "bright-owl-a7"

// Custom separators
generate({ separator: "_" }); // "warm_duck"
generate({ separator: "." }); // "cute.rabbit"
```

### `parse(id, separator?)`

Parse a memorable ID back to its components.

**Parameters:**
- `id` (string): The memorable ID to parse
- `separator` (string): Separator used (default: "-")

**Returns:** `ParsedId` - Object with `components` array and `suffix` string

**Examples:**

```typescript
parse("cute-rabbit-042");
// { components: ["cute", "rabbit"], suffix: "042" }

parse("large-fox-swim");
// { components: ["large", "fox", "swim"], suffix: null }

parse("warm_duck_123", "_");
// { components: ["warm", "duck"], suffix: "123" }
```

### Suffix Generators

Pre-built suffix generators for common use cases:

```typescript
import { suffixGenerators } from 'memorable-ids';

// 3-digit number (000-999) - adds 1,000x multiplier
suffixGenerators.number(); // "042"

// 4-digit number (0000-9999) - adds 10,000x multiplier  
suffixGenerators.number4(); // "1337"

// 2-digit hex (00-ff) - adds 256x multiplier
suffixGenerators.hex(); // "a7"

// Timestamp (last 4 digits) - time-based
suffixGenerators.timestamp(); // "8429"

// Single letter (a-z) - adds 26x multiplier
suffixGenerators.letter(); // "k"
```

### Analysis Functions

Plan capacity and understand collision probabilities:

```typescript
import { 
  calculateCombinations, 
  calculateCollisionProbability,
  getCollisionAnalysis 
} from 'memorable-ids';

// Calculate total possible combinations
calculateCombinations(2); // 5,304 (2 components)
calculateCombinations(2, 1000); // 5,304,000 (2 components + 3-digit suffix)
calculateCombinations(3); // 212,160 (3 components)

// Calculate collision probability (Birthday Paradox)
calculateCollisionProbability(5304, 100); // 0.0093 (0.93% chance)

// Get comprehensive analysis
getCollisionAnalysis(2);
// {
//   totalCombinations: 5304,
//   scenarios: [
//     { ids: 50, probability: 0.0023, percentage: "0.23%" },
//     { ids: 100, probability: 0.0093, percentage: "0.93%" },
//     { ids: 200, probability: 0.037, percentage: "3.7%" },
//     { ids: 500, probability: 0.218, percentage: "21.8%" }
//   ]
// }
```

## Capacity & Collision Analysis

### Total Combinations by Component Count

| Components | Total IDs   | Example                          |
|------------|-------------|----------------------------------|
| 1          | 78          | `bright`                         |
| 2          | 5,304       | `cute-rabbit`                    |
| 3          | 212,160     | `large-fox-swim`                 |
| 4          | 5,728,320   | `happy-owl-dance-quietly`        |
| 5          | 148,936,320 | `clever-fox-run-quickly-through` |

### Suffix Multipliers

| Suffix Type    | Multiplier | Example            |
|----------------|------------|--------------------|
| 3-digit number | ×1,000     | `cute-rabbit-042`  |
| 4-digit number | ×10,000    | `cute-rabbit-1337` |
| 2-digit hex    | ×256       | `cute-rabbit-a7`   |
| Single letter  | ×26        | `cute-rabbit-k`    |

### Collision Probability Examples

**For 2 components (5,304 total combinations):**
- 50 IDs: 0.23% collision chance
- 100 IDs: 0.93% collision chance  
- 200 IDs: 3.7% collision chance
- 500 IDs: 21.8% collision chance

**For 3 components (212,160 total combinations):**
- 1,000 IDs: 0.002% collision chance
- 5,000 IDs: 0.059% collision chance
- 10,000 IDs: 0.235% collision chance
- 20,000 IDs: 0.94% collision chance

**For 2 components + 3-digit suffix (5,304,000 total combinations):**
- 10,000 IDs: 0.0009% collision chance
- 50,000 IDs: 0.023% collision chance
- 100,000 IDs: 0.094% collision chance
- 500,000 IDs: 2.35% collision chance

## Configuration Recommendations

Choose the right configuration based on your expected ID volume:

| Use Case                  | Recommendation          | Capacity | Example               |
|---------------------------|-------------------------|----------|-----------------------|
| Small apps (<1K IDs)      | 2 components            | 5,304    | `cute-rabbit`         |
| Medium apps (1K-50K IDs)  | 3 components            | 212,160  | `large-fox-swim`      |
| Large apps (50K-500K IDs) | 2-3 components + suffix | 5M+      | `cute-rabbit-042`     |
| Enterprise (500K+ IDs)    | 4+ components + suffix  | 50M+     | `happy-owl-dance-042` |

## Advanced Usage

### Custom Suffix Generators

Create your own suffix logic:

```typescript
// Custom timestamp suffix
const timestampSuffix = () => {
  return Date.now().toString().slice(-6); // Last 6 digits
};

// Custom random string
const randomString = () => {
  return Math.random().toString(36).substring(2, 5); // 3 random chars
};

// Use custom suffix
generate({ 
  components: 2, 
  suffix: timestampSuffix 
}); // "cute-rabbit-123456"
```

### Dictionary Access

Access the underlying word collections:

```typescript
import { dictionary } from 'memorable-ids';

console.log(dictionary.adjectives.length); // 78
console.log(dictionary.nouns.length); // 68
console.log(dictionary.verbs.length); // 40
console.log(dictionary.adverbs.length); // 27
console.log(dictionary.prepositions.length); // 26

// Access individual words
console.log(dictionary.adjectives[0]); // "cute"
console.log(dictionary.nouns[0]); // "rabbit"
```

### Error Handling

```typescript
try {
  generate({ components: 6 }); // Invalid: max is 5
} catch (error) {
  console.error(error.message); // "Components must be between 1 and 5"
}
```

## Performance Considerations

### Generation Speed
- **~1M IDs per second** on modern hardware
- No significant performance difference between component counts
- Suffix generation adds minimal overhead

### Randomness Quality
- Uses `Math.random()` - suitable for non-cryptographic purposes
- For cryptographic security, replace with `crypto.getRandomValues()`
- Distribution is uniform across all vocabulary combinations

## Security Notes

⚠️ **Important Security Information:**

- IDs are **NOT cryptographically secure**
- Predictable if `Math.random()` seed is known
- **Suitable for**: user-friendly identifiers, temporary IDs, non-sensitive references
- **NOT suitable for**: session tokens, passwords, security-critical identifiers

For cryptographic security, implement custom random function:

```typescript
import { randomBytes } from 'crypto';

const cryptoRandom = (max: number) => {
  const bytes = randomBytes(4);
  return bytes.readUInt32BE(0) % max;
};

// Replace Math.random() in your custom implementation
```

## TypeScript Support

Full TypeScript support with comprehensive type definitions:

```typescript
import type { 
  GenerateOptions, 
  ParsedId, 
  SuffixGenerator,
  CollisionAnalysis 
} from 'memorable-ids';

const options: GenerateOptions = {
  components: 3,
  suffix: suffixGenerators.number,
  separator: "-"
};

const result: ParsedId = parse("cute-rabbit-042");
const analysis: CollisionAnalysis = getCollisionAnalysis(2);
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
# Clone repository
git clone https://github.com/riipandi/memorable-ids.git
cd memorable-ids

# Install dependencies
pnpm install

# Build
pnpm build

# Run tests
pnpm test
```

### Adding Custom Vocabulary

1. Extend existing arrays in `src/dictionary.ts`
2. Ensure words are URL-safe and human-readable
3. Avoid duplicates to maintain combination count accuracy
4. Update tests and documentation

## License

This project is open-sourced software licensed under the [MIT license](https://choosealicense.com/licenses/mit/).

Copyrights in this project are retained by their contributors.
See the [license file](./LICENSE) for more information.

---

<sub>🤫 Psst! If you like my work you can support me via [GitHub sponsors](https://github.com/sponsors/riipandi).</sub>
