package receipt

import (
	"math"
	"sort"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/invopop/jsonschema"
)

// WordInfo holds the description and average x/y of a word box.
type WordInfo struct {
	Text string
	AvgY float64
	AvgX float64
}

// GroupTextByLines groups annotations into lines based on Y-coordinate proximity.
func GroupTextByLines(annotations []*visionpb.EntityAnnotation, yThreshold int) []string {
	if len(annotations) < 2 {
		return []string{}
	}

	wordAnnotations := annotations[1:] // skip the full text

	var words []WordInfo
	for _, ann := range wordAnnotations {
		vertices := ann.BoundingPoly.Vertices
		if len(vertices) == 0 {
			continue
		}

		var sumX, sumY float64
		for _, v := range vertices {
			sumX += float64(v.X)
			sumY += float64(v.Y)
		}
		avgX := sumX / float64(len(vertices))
		avgY := sumY / float64(len(vertices))

		words = append(words, WordInfo{
			Text: ann.Description,
			AvgX: avgX,
			AvgY: avgY,
		})
	}

	// Group by Y axis
	var lines [][]WordInfo
	for _, word := range words {
		placed := false
		for i, line := range lines {
			if len(line) > 0 && math.Abs(line[0].AvgY-word.AvgY) <= float64(yThreshold) {
				lines[i] = append(lines[i], word)
				placed = true
				break
			}
		}
		if !placed {
			lines = append(lines, []WordInfo{word})
		}
	}

	// Sort each line by X and join into final strings
	var result []string
	for _, line := range lines {
		sort.SliceStable(line, func(i, j int) bool {
			return line[i].AvgX < line[j].AvgX
		})
		text := ""
		for i, word := range line {
			if i > 0 {
				text += " "
			}
			text += word.Text
		}
		result = append(result, text)
	}

	return result
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
