// Package loader provides image I/O operations with validation and format support
// for the Omarchy Theme Generator. It handles loading, validating, and extracting
// metadata from image files in a memory-efficient manner.
//
// Supported Formats:
//   - JPEG (.jpg, .jpeg)
//   - PNG (.png)
//
// Key Features:
//   - Memory-efficient image processing
//   - Format validation and error handling
//   - Image metadata extraction (dimensions, pixel count)
//   - Settings-as-methods pattern enforcement
//   - Comprehensive error handling with context
//
// Usage:
//
//	settings := settings.DefaultSettings()
//	loader := loader.NewFileLoader(settings)
//	
//	// Load image from file
//	img, err := loader.LoadImage(ctx, "wallpaper.jpg")
//	if err != nil {
//	    return err
//	}
//	
//	// Get image metadata
//	info, err := loader.GetImageInfo(ctx, "wallpaper.jpg")
//	if err != nil {
//	    return err
//	}
//	
//	fmt.Printf("Image: %dx%d (%d pixels)\n", 
//	    info.Width, info.Height, info.PixelCount())
//
// The loader validates file formats against configured allowed formats
// and provides detailed error information for debugging and user feedback.
package loader