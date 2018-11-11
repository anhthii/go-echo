package lrc_parser_test

import (
	"io/ioutil"
	"testing"

	"github.com/anhthii/go-echo/lrc_parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	// t.Fatal("not implemented")
	bytes, err := ioutil.ReadFile("song.lrc")
	if err != nil {
		t.Error(err)
	}

	data := string(bytes)

	result := lrc_parser.Parse(data)
	lyrics := result["scripts"].([]lrc_parser.Lyric)
	assert.Equal(t, result["length"].(string), "03:55", "length should be equal")
	assert.Equal(t, result["ar"].(string), "Maroon 5, Cardi B", "artist should be the same")
	assert.Equal(t, result["al"].(string), "Girls Like You (Single)", "album should be the same")
	assert.Equal(t, result["ti"].(string), "Girls Like You", "title should be the same")
	assert.Equal(t, lyrics[1].Script, "Ca sĩ: Maroon 5, Cardi B", "script should match")
	assert.Equal(t, lyrics[1].Start, float64(2), "time sould match")
	assert.Equal(t, lyrics[1].End, float64(4), "time sould match")
	assert.Equal(t, lyrics[2].Script, "Spent 24 hours", "script should match")
	assert.Equal(t, lyrics[2].Start, float64(6.79), "time sould match")
	assert.Equal(t, lyrics[2].End, float64(9.10), "time sould match")
	assert.Equal(t, lyrics[len(lyrics)-1].Start, float64(222.21), "time sould match")
	assert.Equal(t, lyrics[len(lyrics)-1].End, float64(226.25), "time sould match")
}
