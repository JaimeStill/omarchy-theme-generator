# Test Suite Documentation

Integration tests for the omarchy-theme-generator color extraction system using real wallpaper images.

## Test Structure

### Strategy Selection Tests (`strategies_test.go`)

#### TestStrategySelection
Validates that images trigger the expected extraction strategies based on their visual characteristics.

| Image | Strategy | Colors | Dominant Color | Percentage |
|-------|----------|---------|----------------|------------|
| ![nebula](images/nebula.jpeg) | saliency | 69,489 | #200b1e | 0.4% |
| ![night-city](images/night-city.jpeg) | saliency | 121,176 | #000101 | 1.6% |
| ![grayscale](images/grayscale.jpeg) | frequency | 256 | #0e0e0e | 2.0% |
| ![mountains](images/mountains.jpeg) | saliency | 199,953 | #000400 | 0.7% |
| ![abstract](images/abstract.jpeg) | saliency | 127,756 | #f39c71 | 0.3% |

```
=== RUN   TestStrategySelection
=== RUN   TestStrategySelection/nebula.jpeg
    strategies_test.go:58: nebula.jpeg: strategy=saliency, colors=69489, dominant=#200b1e (0.4%)
=== RUN   TestStrategySelection/night-city.jpeg
    strategies_test.go:58: night-city.jpeg: strategy=saliency, colors=121176, dominant=#000101 (1.6%)
=== RUN   TestStrategySelection/grayscale.jpeg
    strategies_test.go:58: grayscale.jpeg: strategy=frequency, colors=256, dominant=#0e0e0e (2.0%)
=== RUN   TestStrategySelection/mountains.jpeg
    strategies_test.go:58: mountains.jpeg: strategy=saliency, colors=199953, dominant=#000400 (0.7%)
=== RUN   TestStrategySelection/abstract.jpeg
    strategies_test.go:58: abstract.jpeg: strategy=saliency, colors=127756, dominant=#f39c71 (0.3%)
--- PASS: TestStrategySelection (5.19s)
```

#### TestThemeGenerationAnalysis
Validates theme generation analysis for different image types.

| Image | Strategy | Grayscale | Monochromatic | Saturation |
|-------|----------|-----------|---------------|------------|
| ![nebula](images/nebula.jpeg) | extract | false | true | 0.474 |
| ![grayscale](images/grayscale.jpeg) | extract | true | false | 0.000 |
| ![mountains](images/mountains.jpeg) | extract | false | false | 0.600 |

```
=== RUN   TestThemeGenerationAnalysis
=== RUN   TestThemeGenerationAnalysis/nebula.jpeg
    strategies_test.go:101: nebula.jpeg: strategy=extract, grayscale=false, monochromatic=true, sat=0.474
=== RUN   TestThemeGenerationAnalysis/grayscale.jpeg
    strategies_test.go:101: grayscale.jpeg: strategy=extract, grayscale=true, monochromatic=false, sat=0.000
=== RUN   TestThemeGenerationAnalysis/mountains.jpeg
    strategies_test.go:101: mountains.jpeg: strategy=extract, grayscale=false, monochromatic=false, sat=0.600
--- PASS: TestThemeGenerationAnalysis (2.86s)
```

#### TestSaliencyVsFrequency
Validates that saliency strategy extracts saturated colors from complex images.

![nebula](images/nebula.jpeg)

```
=== RUN   TestSaliencyVsFrequency
--- PASS: TestSaliencyVsFrequency (2.09s)
```

### Benchmark Tests

Performance benchmarks for extraction strategies.

| Strategy | Target Image | Performance | Operations |
|----------|-------------|-------------|------------|
| Saliency | ![nebula](images/nebula.jpeg) | 2.07 seconds | 1 |
| Frequency | ![grayscale](images/grayscale.jpeg) | 206ms | 5 |

```
goos: linux
goarch: amd64
pkg: github.com/JaimeStill/omarchy-theme-generator/tests
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkSaliencyStrategy
BenchmarkSaliencyStrategy-8    	       1	2070687429 ns/op
BenchmarkFrequencyStrategy
BenchmarkFrequencyStrategy-8   	       5	 206148251 ns/op
```

## Test Images

The `images/` directory contains wallpaper samples for validation. See `images/README.md` for detailed characteristics analysis.

## Utilities

- `analyze-images/` - Generates comprehensive image analysis documentation

## Running Tests

```bash
# Run all tests
go test ./tests -v

# Run specific test suites
go test ./tests -run TestStrategySelection -v
go test ./tests -run TestThemeGeneration -v

# Run benchmarks
go test ./tests -bench=. -v
```