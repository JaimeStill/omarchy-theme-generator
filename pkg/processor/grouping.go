package processor

import (
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func (p *Processor) groupByLightness(colors []WeightedColor) LightnessGroups {
	groups := NewLightnessGroups()

	darkMax := p.settings.LightnessDarkMax
	lightMin := p.settings.LightnessLightMin

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		switch {
		case hsla.L <= darkMax:
			groups.Dark = append(groups.Dark, wc)
		case hsla.L >= lightMin:
			groups.Light = append(groups.Light, wc)
		default:
			groups.Mid = append(groups.Mid, wc)
		}
	}

	sortByWeight(groups.Dark)
	sortByWeight(groups.Mid)
	sortByWeight(groups.Light)

	return groups
}

func (p *Processor) groupBySaturation(colors []WeightedColor) SaturationGroups {
	groups := NewSaturationGroups()

	grayMax := p.settings.SaturationGrayMax
	mutedMax := p.settings.SaturationMutedMax
	normalMax := p.settings.SaturationNormalMax

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		switch {
		case hsla.S <= grayMax:
			groups.Gray = append(groups.Gray, wc)
		case hsla.S <= mutedMax:
			groups.Muted = append(groups.Muted, wc)
		case hsla.S <= normalMax:
			groups.Normal = append(groups.Normal, wc)
		default:
			groups.Vibrant = append(groups.Vibrant, wc)
		}
	}

	sortByWeight(groups.Gray)
	sortByWeight(groups.Muted)
	sortByWeight(groups.Normal)
	sortByWeight(groups.Vibrant)

	return groups
}

func (p *Processor) groupByHue(colors []WeightedColor) HueFamilies {
	families := make(HueFamilies)
	sectorSize := p.settings.HueSectorSize
	sectorCount := p.settings.HueSectorCount
	grayscaleThreshold := p.settings.GrayscaleThreshold

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		if hsla.S < grayscaleThreshold {
			continue
		}

		sector := int(hsla.H / sectorSize)

		if sector >= sectorCount {
			sector = 0
		}

		families[sector] = append(families[sector], wc)
	}

	for sector := range families {
		sortByWeight(families[sector])
	}

	return families
}

func sortByWeight(colors []WeightedColor) {
	sort.Slice(colors, func(i, j int) bool {
		return colors[i].Weight > colors[j].Weight
	})
}
