package useragent

import (
	"container/ring"
	"embed"
	"encoding/json"
	"math/rand"
	"sync"
)

var (
	rings map[UserAgent]*ring.Ring
	mux   sync.Mutex
)

type UserAgent uint16

func (u UserAgent) String() string {
	return stringsKebab[u]
}

const (
	Any            UserAgent = 0x0000
	Desktop        UserAgent = 0x1000
	DesktopWindows UserAgent = 0x1001
	DesktopMacOS   UserAgent = 0x1002
	Mobile         UserAgent = 0x2000
	MobileAndroid  UserAgent = 0x2001
	MobileIOS      UserAgent = 0x2002
	Tablet         UserAgent = 0x3000
	TabletAndroid  UserAgent = 0x3001
	TabletIOS      UserAgent = 0x3002
)

var values = []UserAgent{
	Any,
	Desktop,
	DesktopWindows,
	DesktopMacOS,
	Mobile,
	MobileAndroid,
	MobileIOS,
	Tablet,
	TabletAndroid,
	TabletIOS,
}

var stringsKebab = map[UserAgent]string{
	Any:            "any",
	Desktop:        "desktop",
	DesktopWindows: "desktop-windows",
	DesktopMacOS:   "desktop-macos",
	Mobile:         "mobile",
	MobileAndroid:  "mobile-android",
	MobileIOS:      "mobile-ios",
	Tablet:         "tablet",
	TabletAndroid:  "tablet-android",
	TabletIOS:      "tablet-ios",
}

var stringsCamel = map[UserAgent]string{
	Any:            "Any",
	Desktop:        "Desktop",
	DesktopWindows: "DesktopWindows",
	DesktopMacOS:   "DesktopMacos",
	Mobile:         "Mobile",
	MobileAndroid:  "MobileAndroid",
	MobileIOS:      "MobileIOS",
	Tablet:         "Tablet",
	TabletAndroid:  "TabletAndroid",
	TabletIOS:      "TabletIOS",
}

func init() {
	rings = make(map[UserAgent]*ring.Ring)

	rings[DesktopMacOS] = ringOrPanic(ringFromFilename("desktop-macos.json"))
	rings[DesktopWindows] = ringOrPanic(ringFromFilename("desktop-windows.json"))
	rings[Desktop] = ringOrPanic(ringFromFilenames(
		"desktop-macos.json",
		"desktop-windows.json",
	))

	rings[MobileAndroid] = ringOrPanic(ringFromFilename("mobile-android.json"))
	rings[MobileIOS] = ringOrPanic(ringFromFilename("mobile-ios.json"))
	rings[Mobile] = ringOrPanic(ringFromFilenames(
		"mobile-android.json",
		"mobile-ios.json",
	))

	rings[TabletAndroid] = ringOrPanic(ringFromFilename("tablet-android.json"))
	rings[TabletIOS] = ringOrPanic(ringFromFilename("tablet-ios.json"))
	rings[Tablet] = ringOrPanic(ringFromFilenames(
		"tablet-android.json",
		"tablet-ios.json",
	))

	rings[Any] = ringOrPanic(ringFromFilenames(
		"desktop-macos.json",
		"desktop-windows.json",
		"mobile-android.json",
		"mobile-ios.json",
	))
}

//go:embed *.json
var data embed.FS

func ringFromFilenameX(fn string) *ring.Ring {
	rng, err := ringFromFilename(fn)
	if err != nil {
		panic(err)
	}
	return rng
}

func ringFromFilename(fn string) (*ring.Ring, error) {
	r, err := data.Open(fn)
	if err != nil {
		return nil, err
	}
	var values []string
	err = json.NewDecoder(r).Decode(&values)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
	var rng = ring.New(len(values))
	for i := range values {
		rng.Value = values[i]
		rng = rng.Next()
	}
	return rng, nil
}

func ringFromFilenames(fn ...string) (rng *ring.Ring, err error) {
	for i := range fn {
		sub, err := ringFromFilename(fn[i])
		if err != nil {
			return nil, err
		}
		if rng == nil {
			rng = sub
		} else {
			rng = rng.Link(sub)
		}
	}
	return rng, nil
}

func ringOrPanic(rng *ring.Ring, err error) *ring.Ring {
	if err != nil {
		panic(err)
	}
	return rng
}

func String(ua UserAgent) string {
	mux.Lock()
	r := rings[ua]
	s := r.Value.(string)
	rings[ua] = r.Next()
	mux.Unlock()
	return s
}

var (
	valuesCamel = make(map[string]UserAgent)
	valuesKebab = make(map[string]UserAgent)
)

func init() {
	for i, s := range stringsCamel {
		valuesCamel[s] = i
	}
	for i, s := range stringsKebab {
		valuesKebab[s] = i
	}
}

func FromString(s string) UserAgent {
	if i, ok := valuesCamel[s]; ok {
		return i
	}
	if i, ok := valuesKebab[s]; ok {
		return i
	}
	return Any
}
