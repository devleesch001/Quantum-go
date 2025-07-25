package maps

import (
	"errors"
)

type Maps struct {
	maps  []byte
	doors []*Door
	w     uint
	h     uint
}

func New() (*Maps, error) {
	var m Maps

	return &m, m.init()
}

func (m *Maps) init() error {
	m.maps = DefaultMaps

	for i := 0; i < len(m.maps); i++ {
		if m.maps[i] == 'v' || m.maps[i] == 'h' {
			var d = &Door{
				enabled: 1,
				y:       i / MapW,
				x:       i % MapW,
				state:   0,
			}

			if m.maps[i] == 'v' {
				d.vertical = 1
			} else {
				d.vertical = 0
			}

			m.doors = append(m.doors, d)
		}
	}

	if len(m.doors) > MaxDoors {
		return errors.New("too many doors in the map")
	}

	m.h = MapH
	m.w = MapW
	return nil
}

type Door struct {
	enabled  int
	x        int
	y        int
	vertical int
	state    int
}
