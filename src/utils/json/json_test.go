package json

import (
	"log"
	"testing"
)

func TestJsonGroup(t *testing.T) {
	b := jres{
		Obj{"a": 1, "b": 2, "v": 1},
		Obj{"a": 1, "b": 2, "v": 2},
		Obj{"a": 1, "b": 3, "v": 2},
		Obj{"a": 2, "b": 4, "v": 3}}
	r := JsonGroup(b, []string{"a", "b"})
	if ToJson(r) != `{
    "1": {
        "2": [
            {
                "v": 1
            },
            {
                "v": 2
            }
        ],
        "3": [
            {
                "v": 2
            }
        ]
    },
    "2": {
        "4": [
            {
                "v": 3
            }
        ]
    }
}` {
		t.Errorf("strip option")
	}

	r = JsonGroup(b, []string{"a", "b"}, "nostrip")
	log.Printf("%+v", ToJson(r))

}

func TestAppend(t *testing.T) {
	a := 1
	c := Append(a, "fuck")
	log.Printf("%+v", c)
}

func TestOmit(t *testing.T) {
	a := Obj{"a": 1, "b": 2, "v": 1}
	c := Omit(a, []string{"a", "v"})
	log.Printf("%+v", c)
}

func TestGroup(t *testing.T) {
	b := jres{
		Obj{"a": 1, "b": 2, "v": 1},
		Obj{"a": 1, "b": 2, "v": 2},
		Obj{"a": 1, "b": 3, "v": 2},
		Obj{"a": 2, "b": 4, "v": 3}}
	r := Group(b, []string{"a", "b"})
	if ToJson(r) != `{
    "1": {
        "2": [
            {
                "v": 1
            },
            {
                "v": 2
            }
        ],
        "3": [
            {
                "v": 2
            }
        ]
    },
    "2": {
        "4": [
            {
                "v": 3
            }
        ]
    }
}` {
		t.Errorf("strip option")
	}

	r = Group(b, []string{"a", "b"}, "nostrip")

	if ToJson(r) != `{
    "1": {
        "2": [
            {
                "a": 1,
                "b": 2,
                "v": 1
            },
            {
                "a": 1,
                "b": 2,
                "v": 2
            }
        ],
        "3": [
            {
                "a": 1,
                "b": 3,
                "v": 2
            }
        ]
    },
    "2": {
        "4": [
            {
                "a": 2,
                "b": 4,
                "v": 3
            }
        ]
    }
}` {
		t.Errorf("no strip")
	}

	log.Printf("%+v", ToJson(r))
}

func Test_str2var(t *testing.T) {
	str := `{
    "1": {
        "2": [
            {
                "v": 1
            },
            {
                "v": 2
            }
        ],
        "3": [
            {
                "v": 2
            }
        ]
    },
    "2": {
        "4": [
            {
                "v": 3
            }
        ]
    }
}`
	log.Printf("%+v", Str2Var(str))

}
