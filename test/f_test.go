/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/andypangaribuan/project9"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

func TestRetry(t *testing.T) {
	var i int = 0

	f := func() (*string, error) {
		i++
		if i == 10 {
			v := "value"
			return &v, nil
		}

		log.Printf("invalid i: %v\n", i)
		return nil, fmt.Errorf("invalid i : %v", i)
	}

	r, err := f9.Retry(f, 4, time.Second)
	if err != nil {
		log.Printf("error: %v\n", err)
		t.Error(err)
		return
	}

	log.Printf("done: %v\n", *r)
	log.Printf("%v", *r)
}

func TestXID(t *testing.T) {
	project9.Initialize()
	// tm := f9.TimeNow()

	// t.Log(tm.Unix())
	// t.Log(tm.UnixMilli())
	// t.Log(tm.UnixMicro())
	// t.Log(tm.UnixNano())

	hex := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	nine := p9.Util.GetRandomAlphabetNumber(9)
	xid := hex + nine
	log.Println(xid)
}

func TestCIDLv1(t *testing.T) {
	// total: 52 item
	var (
		rl = rune('a')
		ru = rune('A')
		ch = make([]string, 0)
	)

	for i := 0; i < 26; i++ {
		l1 := string(rl + int32(i))
		u1 := string(ru + int32(i))

		ch = append(ch,
			l1,
			u1,
		)
	}

	checkUnique(ch)

	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}
}

func TestCIDLv2(t *testing.T) {
	// total: 2.704 item
	var (
		rl = rune('a')
		ru = rune('A')
		ch = make([]string, 0)
	)

	for i := 0; i < 26; i++ {
		l1 := string(rl + int32(i))
		u1 := string(ru + int32(i))

		for j := 0; j < 26; j++ {
			l2 := string(rl + int32(j))
			u2 := string(ru + int32(j))

			ch = append(ch,
				l1+l2,
				l1+u2,
				u1+l2,
				u1+u2,
			)
		}
	}

	checkUnique(ch)

	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}
}

func TestCIDLv3(t *testing.T) {
	// total: 140.608 item
	// sub ms: 10
	var (
		rl = rune('a')
		ru = rune('A')
		ch = make([]string, 0)
	)

	tmStart := f9.TimeNow()
	for i := 0; i < 26; i++ {
		l1 := string(rl + int32(i))
		u1 := string(ru + int32(i))

		for j := 0; j < 26; j++ {
			l2 := string(rl + int32(j))
			u2 := string(ru + int32(j))

			for k := 0; k < 26; k++ {
				l3 := string(rl + int32(k))
				u3 := string(ru + int32(k))

				ch = append(ch,
					l1+l2+l3,
					l1+l2+u3,
					l1+u2+l3,
					l1+u2+u3,
					u1+l2+l3,
					u1+l2+u3,
					u1+u2+l3,
					u1+u2+u3,
				)
			}
		}
	}
	tmFinish := f9.TimeNow()

	checkUnique(ch)
	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}

	ms := tmFinish.Sub(tmStart).Milliseconds()
	fmt.Printf("sub ms: %v\n", ms)
}

func TestCIDLv4(t *testing.T) {
	// total: 7.311.616 item
	var (
		rl = rune('a')
		ru = rune('A')
		ch = make([]string, 0)
	)

	for i := 0; i < 26; i++ {
		l1 := string(rl + int32(i))
		u1 := string(ru + int32(i))

		for j := 0; j < 26; j++ {
			l2 := string(rl + int32(j))
			u2 := string(ru + int32(j))

			for k := 0; k < 26; k++ {
				l3 := string(rl + int32(k))
				u3 := string(ru + int32(k))

				for l := 0; l < 26; l++ {
					l4 := string(rl + int32(l))
					u4 := string(ru + int32(l))

					ch = append(ch,
						l1+l2+l3+l4, // llll
						l1+l2+l3+u4, // lllu
						l1+l2+u3+l4, // llul
						l1+l2+u3+u4, // lluu
						l1+u2+l3+l4, // lull
						l1+u2+l3+u4, // lulu
						l1+u2+u3+l4, // luul
						l1+u2+u3+u4, // luuu
						u1+l2+l3+l4, // ulll
						u1+l2+l3+u4, // ullu
						u1+l2+u3+l4, // ulul
						u1+l2+u3+u4, // uluu
						u1+u2+l3+l4, // uull
						u1+u2+l3+u4, // uulu
						u1+u2+u3+l4, // uuul
						u1+u2+u3+u4, // uuuu
					)
				}
			}
		}
	}

	checkUnique(ch)

	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}
}

func TestCIDN1(t *testing.T) {
	// total: 62 item
	ch := getCN()

	checkUnique(ch)
	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}
}

func TestCIDN2(t *testing.T) {
	// total: 3.844 item
	var (
		cn = getCN()
		ch = make([]string, 0)
	)

	for _, v1 := range cn {
		for _, v2 := range cn {
			ch = append(ch, v1+v2)
		}
	}

	checkUnique(ch)
	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}
}

func TestCIDN3(t *testing.T) {
	// total: 238.328 item
	// sub ms: 22
	var (
		cn = getCN()
		ch = make([]string, 0)
	)

	tmStart := f9.TimeNow()
	for _, v1 := range cn {
		for _, v2 := range cn {
			for _, v3 := range cn {
				ch = append(ch, v1+v2+v3)
			}
		}
	}
	tmFinish := f9.TimeNow()

	checkUnique(ch)
	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}

	ms := tmFinish.Sub(tmStart).Milliseconds()
	fmt.Printf("sub ms: %v\n", ms)
}

func TestCIDN4(t *testing.T) {
	// total: 14.776.336 item
	// sub ms: 2006
	var (
		cn = getCN()
		ch = make([]string, 0)
	)

	tmStart := f9.TimeNow()
	for _, v1 := range cn {
		for _, v2 := range cn {
			for _, v3 := range cn {
				for _, v4 := range cn {
					ch = append(ch, v1+v2+v3+v4)
				}
			}
		}
	}
	tmFinish := f9.TimeNow()

	checkUnique(ch)
	for i, v := range ch {
		fmt.Printf("%v: %v\n", i+1, v)
	}

	ms := tmFinish.Sub(tmStart).Milliseconds()
	fmt.Printf("sub ms: %v\n", ms)
}

func getCN() []string {
	var (
		rn = rune('0')
		rl = rune('a')
		ru = rune('A')
		cn = make([]string, 0)
		// cl = make([]string, 0)
		// cu = make([]string, 0)
	)

	for i := 0; i < 10; i++ {
		n1 := string(rn + int32(i))
		cn = append(cn, n1)
	}

	for i := 0; i < 26; i++ {
		l := string(rl + int32(i))
		u := string(ru + int32(i))

		// cl = append(cl, l)
		// cu = append(cu, u)

		cn = append(cn, l, u)
	}

	// cn = append(cn, cl...)
	// cn = append(cn, cu...)

	return cn
}

func TestUnixMicro(t *testing.T) {
	project9.Initialize()

	tm := f9.TimeNow()
	fmt.Printf("unix milli: %v\n", tm.UnixMilli())
	fmt.Printf("unix micro: %v\n", tm.UnixMicro())
	fmt.Printf("unix nano : %v\n", tm.UnixNano())
}

func checkUnique(ch []string) {
	l := len(ch)
	mp := make(map[string]interface{}, 0)

	for _, v := range ch {
		if _, ok := mp[v]; !ok {
			mp[v] = nil
		} else {
			fmt.Printf("duplicate: %v\n", v)
		}
	}

	mpl := len(mp)

	if l != mpl {
		fmt.Printf("HAVE DUPLICATE\n")
	}
}
