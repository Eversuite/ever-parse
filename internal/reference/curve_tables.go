package reference

import (
	"encoding/json"
	"ever-parse/internal/util"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

// CurveTableReference type alias to work with reference entries.
// The CurveTableReference contains an arbitrary amount of properties with their corresponding CurveTable references.
type CurveTableReference map[string]CurveTableReferenceEntry

// CurveTableReferenceEntry contains a reference to the CurveTable and its corresponding row associated with the property
type CurveTableReferenceEntry struct {
	CurveTable struct {
		ObjectPath string
	}
	RowName string
}

type CurvePropertiesMapping interface {
	GetCurveProperty() CurveTableReference
}

func GetCurveProperties(m CurvePropertiesMapping) string {
	if m.GetCurveProperty() == nil {
		return ""
	}
	jsonBytes, err := json.Marshal(m.GetCurveProperty().GetValues())
	util.Check(err, m)
	return string(jsonBytes)
}

// GetValues retrieves all CurvePoint s for each entry in the reference.
// A CurveTableReference may contain multiple entries.
func (c CurveTableReference) GetValues() map[string][]CurvePoint {
	result := make(map[string][]CurvePoint, len(c))
	for key, entry := range c {
		points := entry.getValue()
		result[key] = points
	}
	return result
}

type CurvePoint struct {
	Time  float64 `json:"Time"`
	Value float64 `json:"Value"`
}

type CurveDefinition struct {
	Mode   string       `json:"InterpMode"`
	Points []CurvePoint `json:"Keys"`
}

// getValue retrieves all values from a CurveTableReferenceEntry.
// A single entry consists of all points on that curve.
// A CurvePoint is defined by their x and y value where x is usually the 'time' and y the value f(x)
func (ce CurveTableReferenceEntry) getValue() (curvePoints []CurvePoint) {
	if ce.RowName == noneName || len(ce.CurveTable.ObjectPath) == 0 {
		return
	}
	correctRoot := FixRoot(ce.CurveTable.ObjectPath)
	cleanedPath := jsonRegex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, ce, correctRoot, cleanedPath)
	tableValue := gjson.Get(string(content), "#.Rows."+ce.RowName+"|0").String()
	tableValue = whitespaceRegex.ReplaceAllString(tableValue, "")
	if len(tableValue) == 0 {
		return
	}
	definition := CurveDefinition{}
	err = json.Unmarshal([]byte(tableValue), &definition)
	util.Check(err, ce, tableValue)
	curvePoints = interpolateMissingValues(definition)
	return
}

func interpolateMissingValues(cd CurveDefinition) (res []CurvePoint) {
	if strings.Compare(cd.Mode, "ERichCurveInterpMode::RCIM_Linear") != 0 {
		return cd.Points
	}
	start := cd.Points[0]
	end := cd.Points[len(cd.Points)-1]
	for i := start.Time; i <= end.Time; i++ {
		point := interpolateLinearPoint(start, end, i)
		res = append(res, CurvePoint{
			Time:  i,
			Value: point,
		})
	}
	return
}

func interpolateLinearPoint(p1, p2 CurvePoint, x float64) float64 {
	if x < p1.Time || x > p2.Time {
		return 0
	}
	// calculate the slope of the line
	slope := (p2.Value - p1.Value) / (p2.Time - p1.Time)
	// use the slope and point1 to calculate the y-value for x
	return slope*(x-p1.Time) + p1.Value
}

func FixCurveTableValues(ctDescription gjson.Result) (error, CurveTableReference) {
	if ctDescription.Type == gjson.Null {
		return nil, CurveTableReference{}
	}

	// Now we fix our curve table values.
	if !ctDescription.IsArray() {
		// If it's an object we can use the usual unmarshal function of the json package... easy
		var ctRef CurveTableReference
		err := json.Unmarshal([]byte(ctDescription.String()), &ctRef)
		return err, ctRef
	} else {
		// In case of an array we need to get a bit creative
		return fixArrayCurveTableValues(ctDescription)
	}
}

func fixArrayCurveTableValues(res gjson.Result) (error, CurveTableReference) {
	results := res.Array()
	// Under the hood a reference.CurveTableReference is nothing but a map of [string]CTREntry
	entries := map[string]CurveTableReferenceEntry{}
	var err error = nil
	// Each entry of the gjson.Result array is its own entry to our CTR-Map
	for _, result := range results {
		// And here we can unmarshal it again using the json package
		var entry map[string]CurveTableReferenceEntry
		err = json.Unmarshal([]byte(result.String()), &entry)
		// And add everything to the complete map
		// This could technically cause an issue if a key ways present at least twice
		// But I doubt that ever happens. It would cause way more issues for the ECH Devs than for us.
		for key, value := range entry {
			entries[key] = value
		}
	}
	return err, entries
}
