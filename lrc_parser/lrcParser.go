package lrc_parser

import (
	"regexp"
	"strconv"
	"strings"
)

type Lyric struct {
	Start  float64 `json:"start"`
	Script string  `json:"text"`
	End    float64 `json:"end"`
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// ExtractMetadata from data under lrc format
func ExtractMetadata(lines []string, result map[string]interface{}) {
	for _, line := range lines {
		last := len(line) - 1
		infos := strings.Split(line[1:last], ": ")
		tag := infos[0]
		value := infos[1]
		result[tag] = value
	}
}

// convert string to seconds
func convertTime(s string) float64 {
	parts := strings.Split(s, ":")
	minutes, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic(err)
	}
	seconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		panic(err)
	}

	if minutes > 0 {
		seconds = float64(minutes)*60 + seconds

	}
	return float64(int(seconds*100)) / 100
}

// Parse lrc data into map format
func Parse(data string) map[string]interface{} {
	result := make(map[string]interface{})
	var lines []string
	lines = strings.Split(data, "\n")
	_startingTimeRx := `\[(\d*\:\d*\.?\d*)\]` // i.g [00:10.55]
	startingTimeRx := regexp.MustCompile(_startingTimeRx)
	// similar to starting time regex
	endingTimeRx := _startingTimeRx
	scriptRx := `(.+)` // ig Havana ooh na-na (ayy)
	startingTimeAndScriptRx := _startingTimeRx + scriptRx
	lyricLineRegex := regexp.MustCompile(startingTimeAndScriptRx)

	var metadataLines []string
	for i := 0; startingTimeRx.Match([]byte(lines[i])) == false && lines[i] != ""; i++ {
		metadataLines = append(metadataLines, strings.TrimSpace(lines[i]))
	}

	ExtractMetadata(metadataLines, result)
	lyricLines := lines[len(metadataLines):]
	// remove all empty lines
	notEmptyLineRegex := regexp.MustCompile(startingTimeAndScriptRx + "|" + endingTimeRx)
	removeEmptyLine := func(s string) bool { return notEmptyLineRegex.Match([]byte(s)) }
	lyricLines = filter(lyricLines, removeEmptyLine)
	endingTimeRxP := regexp.MustCompile(endingTimeRx)

	result["scripts"] = make([]Lyric, 0)

	for i, length := 0, len(lyricLines); i < length-1; i++ {
		matches := lyricLineRegex.FindStringSubmatch(lyricLines[i])
		endingTimeMatches := endingTimeRxP.FindStringSubmatch(lyricLines[i+1])
		if len(matches) > 0 && len(endingTimeMatches) > 0 {
			startTime := convertTime(matches[1])
			script := matches[2]
			endTime := convertTime(endingTimeMatches[1])
			result["scripts"] = append(result["scripts"].([]Lyric), Lyric{
				Start: startTime, Script: script, End: endTime,
			})
		}
	}

	return result
}
