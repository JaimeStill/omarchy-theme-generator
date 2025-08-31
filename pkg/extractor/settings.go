package extractor

// Settings contains all configurable thresholds and parameters for the extraction system.
// These values were empirically derived from analysis of 15 diverse wallpaper images
// using tests/analyze-thresholds/main.go.
//
// To regenerate optimal values:
//   go run tests/analyze-thresholds/main.go
type Settings struct {
	// Strategy selection thresholds
	Strategy StrategySettings
	
	// Image analysis parameters
	Analysis AnalysisSettings
	
	// Saliency calculation weights
	Saliency SaliencySettings
	
	// Frequency strategy weights
	Frequency FrequencySettings
	
	// Extraction options
	Extraction ExtractionSettings
}

// StrategySettings controls when different extraction strategies are selected
type StrategySettings struct {
	// Edge density threshold for triggering saliency strategy
	// Empirically derived: 0.036 achieves 73% accuracy on test set
	SaliencyEdgeThreshold float64
	
	// Color complexity threshold for saliency consideration
	// Images with >10K colors often have visually important regions
	SaliencyColorComplexity int
	
	// Minimum saturation for color complexity to trigger saliency
	// Prevents grayscale images from triggering complexity check
	SaliencySaturationThreshold float64
}

// AnalysisSettings controls image characteristic analysis
type AnalysisSettings struct {
	// Edge detection thresholds
	EdgeDetectionMinStrength float64 // Minimum gradient to count as edge
	EdgeDetectionSampleRate  int     // Sample every Nth pixel for performance
	
	// Image type classification thresholds
	HighDetailEdgeThreshold   float64 // EdgeDensity > this = HighDetail
	SmoothEdgeThreshold       float64 // EdgeDensity < this = potential Smooth
	SmoothColorThreshold      int     // ColorComplexity > this for Smooth
	LowDetailColorThreshold   int     // ColorComplexity < this = LowDetail
	ComplexColorThreshold     int     // ColorComplexity > this = potential Complex
	ComplexEdgeThreshold      float64 // EdgeDensity > this for Complex
	
	// Region detection
	RegionMinEdgeDensity float64 // Minimum edge density for distinct regions
	RegionMaxEdgeDensity float64 // Maximum edge density for distinct regions
}

// SaliencySettings controls saliency map generation and scoring
type SaliencySettings struct {
	// Saliency calculation weights (must sum to 1.0)
	LocalContrastWeight   float64
	EdgeStrengthWeight    float64
	ColorUniquenessWeight float64
	
	// Sampling parameters for performance
	SaliencyMapSampleRate    int // Sample every Nth pixel
	SaliencyMapSpreadRadius  int // Radius for spreading saliency values
	
	// Local analysis windows
	ContrastWindowRadius    int // Radius for local contrast calculation
	UniquenessWindowRadius  int // Radius for color uniqueness calculation
	
	// Color similarity threshold for uniqueness calculation
	ColorSimilarityThreshold float64 // Euclidean distance in RGB space
	
	// Final color weighting
	FrequencyWeight float64 // Weight for original frequency (vs saliency)
	SaliencyWeight  float64 // Weight for saliency score
}

// FrequencySettings controls frequency-based extraction
type FrequencySettings struct {
	// Visual importance scoring weights by image type
	// Each set of weights: [frequency, saturation, lightness, contrast]
	HighDetailWeights [4]float64
	LowDetailWeights  [4]float64
	SmoothWeights     [4]float64
	ComplexWeights    [4]float64
	DefaultWeights    [4]float64
	
	// Lightness penalties for extreme values
	DarkLightThreshold     float64 // L < this gets penalty
	BrightLightThreshold   float64 // L > this gets penalty
	ExtremeLightnessPenalty float64 // Score multiplier for extreme lightness
	OptimalLightnessBonus   float64 // Score multiplier for optimal range
	
	// Saturation bonus threshold
	MinSaturationForBonus float64 // Minimum saturation to avoid penalty
	
	// Contrast calculation
	MaxContrastSamples int // Maximum colors to check for contrast
}

// ExtractionSettings controls general extraction behavior
type ExtractionSettings struct {
	// Memory optimization
	InitialMapCapacity int // Initial capacity for frequency map
	MaxMapCapacity     int // Maximum capacity to prevent excessive memory
	
	// Performance optimization
	EdgeDensitySampleRate int // Sample rate for edge density calculation
	ColorDistSampleRate   int // Sample rate for color distribution analysis
}

// DefaultSettings returns the empirically-derived optimal settings
func DefaultSettings() *Settings {
	return &Settings{
		Strategy: StrategySettings{
			SaliencyEdgeThreshold:        0.036,
			SaliencyColorComplexity:      10000,
			SaliencySaturationThreshold:  0.4,
		},
		Analysis: AnalysisSettings{
			EdgeDetectionMinStrength:   30.0,
			EdgeDetectionSampleRate:    4,
			HighDetailEdgeThreshold:    0.15,
			SmoothEdgeThreshold:        0.05,
			SmoothColorThreshold:       100,
			LowDetailColorThreshold:    50,
			ComplexColorThreshold:      200,
			ComplexEdgeThreshold:       0.08,
			RegionMinEdgeDensity:       0.05,
			RegionMaxEdgeDensity:       0.25,
		},
		Saliency: SaliencySettings{
			LocalContrastWeight:      0.5,
			EdgeStrengthWeight:       0.3,
			ColorUniquenessWeight:    0.2,
			SaliencyMapSampleRate:    4,
			SaliencyMapSpreadRadius:  2,
			ContrastWindowRadius:     2,
			UniquenessWindowRadius:   3,
			ColorSimilarityThreshold: 30.0,
			FrequencyWeight:          0.3,
			SaliencyWeight:           0.7,
		},
		Frequency: FrequencySettings{
			HighDetailWeights:       [4]float64{0.2, 0.3, 0.2, 0.3},
			LowDetailWeights:        [4]float64{0.5, 0.2, 0.2, 0.1},
			SmoothWeights:           [4]float64{0.2, 0.4, 0.3, 0.1},
			ComplexWeights:          [4]float64{0.25, 0.25, 0.25, 0.25},
			DefaultWeights:          [4]float64{0.3, 0.3, 0.2, 0.2},
			DarkLightThreshold:      0.1,
			BrightLightThreshold:    0.9,
			ExtremeLightnessPenalty: 0.3,
			OptimalLightnessBonus:   1.2,
			MinSaturationForBonus:   0.05,
			MaxContrastSamples:      5,
		},
		Extraction: ExtractionSettings{
			InitialMapCapacity:    1024,
			MaxMapCapacity:        65536,
			EdgeDensitySampleRate: 4,
			ColorDistSampleRate:   2,
		},
	}
}

// Global settings instance - can be replaced for testing or tuning
var CurrentSettings = DefaultSettings()

// WithSettings temporarily replaces settings for testing
func WithSettings(settings *Settings, fn func()) {
	old := CurrentSettings
	CurrentSettings = settings
	defer func() { CurrentSettings = old }()
	fn()
}