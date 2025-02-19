package main

import (
	"bytes"
	"testing"
)

func split_(s string) (gcodes [][]byte) {
	b := bytes.NewBufferString(s)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		gcodes = append(gcodes, line[0:len(line)-1])
	}
	return gcodes
}

func TestConvertThumbnail(t *testing.T) {
	gcodes := split_(`; top_infill_extrusion_width = 0.4
; top_solid_infill_speed = 60%
; top_solid_layers = 6
; generated by PrusaSlicer 2.4.2+arm64 on 2022-05-12 at 06:33:34 UTC
G0
G1

;
; thumbnail begin 16x16 536
; aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
; bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
; thumbnail end
;
G0
G1
;
; thumbnail begin 220x124 6528
; xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
; yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy
; zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
; thumbnail end
;
;

; external perimeters extrusion width = 0.45mm
; perimeters extrusion width = 0.45mm
; infill extrusion width = 0.45mm

`)

	comp := []byte("data:image/png;base64,xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	r := convertThumbnail(gcodes)
	if 0 != bytes.Compare(r, comp) {
		t.Error(r, comp)
	}
}

func TestFindEstimatedTime(t *testing.T) {
	gcodes := split_(`; estimated printing time = 2s
	`)
	if r := findEstimatedTime(gcodes); 2 != r {
		t.Error("2s /", r)
	}
	gcodes = split_(`; estimated printing time = 1m 2s
	`)
	if r := findEstimatedTime(gcodes); 60+2 != r {
		t.Error("1m 2s /", r)
	}
	gcodes = split_(`; estimated printing time = 1h  2s
	`)
	if r := findEstimatedTime(gcodes); 3600+2 != r {
		t.Error("1h 2s /", r)
	}
	gcodes = split_(`; estimated printing time = 2d 1m  2s
	`)
	if r := findEstimatedTime(gcodes); 2*86400+1*60+2 != r {
		t.Error("2d 1m 2s /", r)
	}
	gcodes = split_(`; estimated printing time(first) = 2d
; estimated printing time = 2d 1m  3s
	`)
	if r := findEstimatedTime(gcodes); 2*86400 != r {
		t.Error("2d /", r)
	}
}
