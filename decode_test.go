package nvelope_test

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"testing"

	"github.com/muir/nvelope"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type Complex128 complex128

func (c Complex128) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprint(c)), nil
}

type Complex64 complex64

func (c Complex64) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprint(c)), nil
}

func TestDecodeQuerySimpleParameters(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		Int        int         `json:",omitempty" nvelope:"query,name=int"`
		Int8       int8        `json:",omitempty" nvelope:"query,name=int8"`
		Int16      int16       `json:",omitempty" nvelope:"query,name=int16"`
		Int32      int32       `json:",omitempty" nvelope:"query,name=int32"`
		Int64      int64       `json:",omitempty" nvelope:"query,name=int64"`
		Uint       uint        `json:",omitempty" nvelope:"query,name=uint"`
		Uint8      uint8       `json:",omitempty" nvelope:"query,name=uint8"`
		Uint16     uint16      `json:",omitempty" nvelope:"query,name=uint16"`
		Uint32     uint32      `json:",omitempty" nvelope:"query,name=uint32"`
		Uint64     uint64      `json:",omitempty" nvelope:"query,name=uint64"`
		Float32    float32     `json:",omitempty" nvelope:"query,name=float32"`
		Float64    float64     `json:",omitempty" nvelope:"query,name=float64"`
		String     string      `json:",omitempty" nvelope:"query,name=string"`
		IntP       *int        `json:",omitempty" nvelope:"query,name=intp"`
		Int8P      *int8       `json:",omitempty" nvelope:"query,name=int8p"`
		Int16P     *int16      `json:",omitempty" nvelope:"query,name=int16p"`
		Int32P     *int32      `json:",omitempty" nvelope:"query,name=int32p"`
		Int64P     *int64      `json:",omitempty" nvelope:"query,name=int64p"`
		UintP      *uint       `json:",omitempty" nvelope:"query,name=uintp"`
		Uint8P     *uint8      `json:",omitempty" nvelope:"query,name=uint8p"`
		Uint16P    *uint16     `json:",omitempty" nvelope:"query,name=uint16p"`
		Uint32P    *uint32     `json:",omitempty" nvelope:"query,name=uint32p"`
		Uint64P    *uint64     `json:",omitempty" nvelope:"query,name=uint64p"`
		Float32P   *float32    `json:",omitempty" nvelope:"query,name=float32p"`
		Float64P   *float64    `json:",omitempty" nvelope:"query,name=float64p"`
		StringP    *string     `json:",omitempty" nvelope:"query,name=stringp"`
		Complex64  *Complex64  `json:",omitempty" nvelope:"query,name=complex64"`
		Complex128 *Complex128 `json:",omitempty" nvelope:"query,name=complex128"`
		BoolP      *bool       `json:",omitempty" nvelope:"query,name=boolp"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"Int":135}`, do("/x?int=135"))
	assert.Equal(t, `200->{"Int8":-5}`, do("/x?int8=-5"))
	assert.Equal(t, `200->{"Int16":127}`, do("/x?int16=127"))
	assert.Equal(t, `200->{"Int32":11}`, do("/x?int32=11"))
	assert.Equal(t, `200->{"Int64":-38}`, do("/x?int64=-38"))
	assert.Equal(t, `200->{"Uint":135}`, do("/x?uint=135"))
	assert.Equal(t, `200->{"Uint8":5}`, do("/x?uint8=5"))
	assert.Equal(t, `200->{"Uint16":127}`, do("/x?uint16=127"))
	assert.Equal(t, `200->{"Uint32":11}`, do("/x?uint32=11"))
	assert.Equal(t, `200->{"Uint64":38}`, do("/x?uint64=38"))
	assert.Equal(t, `200->{"Float64":38.7}`, do("/x?float64=38.7"))
	assert.Equal(t, `200->{"Float32":11.1}`, do("/x?float32=11.1"))
	assert.Equal(t, `200->{"String":"fred"}`, do("/x?string=fred"))
	assert.Equal(t, `200->{"IntP":135}`, do("/x?intp=135"))
	assert.Equal(t, `200->{"Int8P":-5}`, do("/x?int8p=-5"))
	assert.Equal(t, `200->{"Int16P":127}`, do("/x?int16p=127"))
	assert.Equal(t, `200->{"Int32P":11}`, do("/x?int32p=11"))
	assert.Equal(t, `200->{"Int64P":-38}`, do("/x?int64p=-38"))
	assert.Equal(t, `200->{"UintP":135}`, do("/x?uintp=135"))
	assert.Equal(t, `200->{"Uint8P":5}`, do("/x?uint8p=5"))
	assert.Equal(t, `200->{"Uint16P":127}`, do("/x?uint16p=127"))
	assert.Equal(t, `200->{"Uint32P":11}`, do("/x?uint32p=11"))
	assert.Equal(t, `200->{"Uint64P":38}`, do("/x?uint64p=38"))
	assert.Equal(t, `200->{"Float64P":38.7}`, do("/x?float64p=38.7"))
	assert.Equal(t, `200->{"Float32P":11.1}`, do("/x?float32p=11.1"))
	assert.Equal(t, `200->{"StringP":"fred"}`, do("/x?stringp=fred"))
	assert.Equal(t, `200->{"Complex64":"(38.7-9.3i)"}`, do("/x?complex64="+url.QueryEscape("38.7-9.3i")))
	assert.Equal(t, `200->{"Complex128":"(11.1+22.1i)"}`, do("/x?complex128="+url.QueryEscape("11.1+22.1i")))
	assert.Equal(t, `200->{"BoolP":false}`, do("/x?boolp=false"))
}

func TestDecodeQueryComplexParameters(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		IntSlice     []int          `json:",omitempty" nvelope:"query,name=intslice,explode=false"`
		Int8Slice    []*int8        `json:",omitempty" nvelope:"query,name=int8slice,explode=true"`
		Int16Slice   []*int8        `json:",omitempty" nvelope:"query,name=int16slice,explode=false,delimiter=space"`
		Int32Slice   *[]*int8       `json:",omitempty" nvelope:"query,name=int32slice,explode=false,delimiter=pipe"`
		MapIntBool   map[int]bool   `json:",omitempty" nvelope:"query,name=mapintbool,explode=false"`
		MapIntString map[int]string `json:",omitempty" nvelope:"query,name=mapintstring,deepObject=true"`
		IntArrayP    *[3]int        `json:",omitempty" nvelope:"query,name=intarrayp,explode=false"`
		Emb1         *struct {
			Int    int    `json:",omitempty" nvelope:"eint"`
			Int8   int8   `json:",omitempty" nvelope:"eint8"`
			Int16  int16  `json:",omitempty" nvelope:"eint16"`
			String string `json:",omitempty"`
		} `json:",omitempty" nvelope:"query,name=emb1,explode=false"`
		Emb2 *struct {
			Int    int    `json:",omitempty" nvelope:"eint"`
			Int8   int8   `json:",omitempty" nvelope:"eint8"`
			Int16  int16  `json:",omitempty" nvelope:"eint16"`
			String string `json:",omitempty"`
		} `json:",omitempty" nvelope:"query,name=emb2,deepObject=true"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"IntSlice":[1,7]}`, do("/x?intslice=1,7"))
	assert.Equal(t, `200->{"Int8Slice":[10,11,12]}`, do("/x?int8slice=10&int8slice=11&int8slice=12"))
	assert.Equal(t, `200->{"Int16Slice":[8,22,-3]}`, do("/x?int16slice=8%2022%20-3"))
	assert.Equal(t, `200->{"Int32Slice":[7,11,13]}`, do("/x?int32slice=7|11|13"))
	assert.Equal(t, `200->{"MapIntBool":{"-9":false,"7":true}}`, do("/x?mapintbool=7,true,-9,false"))
	assert.Equal(t, `200->{"MapIntString":{"-9":"hi","7":"bye"}}`, do("/x?mapintstring[7]=bye&mapintstring[-9]=hi"))
	assert.Equal(t, `200->{"Emb1":{"Int":192,"Int8":-3,"String":"foo"}}`, do("/x?emb1=eint,192,eint8,-3,String,foo"))
	assert.Equal(t, `200->{"Emb2":{"Int":193,"Int8":-4,"String":"bar"}}`, do("/x?emb2[eint]=193&emb2[eint8]=-4&emb2[String]=bar"))
	assert.Equal(t, `200->{"IntArrayP":[7,22,0]}`, do("/x?intarrayp=7,22"))
}

type Foo string

func (fp *Foo) UnmarshalText(b []byte) error {
	*fp = Foo("~" + string(b) + "~")
	return nil
}

func TestDecodeQueryJSONParameters(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		Foo  Foo      `json:",omitempty" nvelope:"query,name=foo,explode=false"`
		FooP *Foo     `json:",omitempty" nvelope:"query,name=foop,explode=false"`
		FooA []Foo    `json:",omitempty" nvelope:"query,name=fooa,explode=true"`
		FooB *[]*Foo  `json:",omitempty" nvelope:"query,name=foob,explode=false"`
		S1   string   `json:",omitempty" nvelope:"query,name=s1,content=application/json"`
		S2   *string  `json:",omitempty" nvelope:"query,name=s2,content=application/json"`
		S3   **string `json:",omitempty" nvelope:"query,name=s3,content=application/json"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"Foo":"~bar~"}`, do("/x?foo=bar"))
	assert.Equal(t, `200->{"FooP":"~baz~"}`, do("/x?foop=baz"))
	assert.Equal(t, `200->{"FooA":["~bar~","~baz~"]}`, do("/x?fooa=bar&fooa=baz"))
	assert.Equal(t, `200->{"FooB":["~bing~","~baz~"]}`, do("/x?foob=bing,baz"))
	assert.Equal(t, `200->{"S1":"doof"}`, do(`/x?s1="doof"`))
	assert.Equal(t, `200->{"S2":"boor"}`, do(`/x?s2="boor"`))
	assert.Equal(t, `200->{"S3":"ppp"}`, do(`/x?s3="ppp"`))
}

func TestDecodeQueryHeaderParameters(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		S  string   `json:",omitempty" nvelope:"header,name=S"`
		A1 []string `json:",omitempty" nvelope:"header,name=A1"`
		A2 []string `json:",omitempty" nvelope:"header,name=A2"`
		A3 []string `json:",omitempty" nvelope:"header,explode=false,name=A3"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"S":"yip"}`, do("/x", header("S", "yip")))
	assert.Equal(t, `200->{"A1":["eee"]}`, do("/x", header("A1", "eee")))
	assert.Equal(t, `200->{"A2":["yia","yo"]}`, do("/x", header("A2", "yia"), header("A2", "yo")))
	assert.Equal(t, `200->{"A3":["cow","boy"]}`, do("/x", header("A3", "cow,boy")))
}

func TestDecodeQueryCookieParameters(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		S  string   `json:",omitempty" nvelope:"cookie,name=S"`
		A1 []string `json:",omitempty" nvelope:"cookie,name=A1"`
		A3 []string `json:",omitempty" nvelope:"cookie,explode=false,name=A3"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"S":"yip"}`, do("/x", cookie("S", "yip")))
	assert.Equal(t, `200->{"A1":["eee"]}`, do("/x", cookie("A1", "eee")))
	assert.Equal(t, `200->{"A3":["cow","boy"]}`, do("/x", cookie("A3", "cow,boy")))
}

func TestDecodeQueryPathParameters(t *testing.T) {
	do := captureOutput("/x/{a}/{b}/{c}", func(s struct {
		A string `json:",omitempty" nvelope:"path,name=a"`
		B *int   `json:",omitempty" nvelope:"path,name=b"`
		C Foo    `json:",omitempty" nvelope:"path,name=c"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"A":"foobar","B":38,"C":"~john~"}`, do("/x/foobar/38/john"))
}

func TestDecodeQueryExplode(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		M map[string]int `json:",omitempty" nvelope:"query,name=m,explode=true"`
		S []string       `json:",omitempty" nvelope:"query,name=s,explode=true"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})
	assert.Equal(t, `200->{"M":{"a":7,"b":8}}`, do("/x?m=a%3D7&m=b%3D8"))
	assert.Equal(t, `200->{"S":["x","y","z"]}`, do("/x?s=x&s=y&s=z"))
}

type thing struct {
	I int     `json:"I,omitempty"`
	F float64 `json:"F,omitempty"`
}

func e(s string) string { return url.QueryEscape(s) }

func TestDecodeQueryContentExplode(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		MEA map[int]thing `json:",omitempty" nvelope:"query,name=mea,explode=true,content=application/json"`
		SEA []thing       `json:",omitempty" nvelope:"query,name=sea,explode=true,content=application/json"`
		ME  map[int]int   `json:",omitempty" nvelope:"query,name=me,explode=true"`
		SE  []int         `json:",omitempty" nvelope:"query,name=se,explode=true"`
		MA  map[int]thing `json:",omitempty" nvelope:"query,name=ma,explode=false,content=application/json"`
		SA  []thing       `json:",omitempty" nvelope:"query,name=sa,explode=false,content=application/json"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})

	assert.Equal(t, `200->{"MEA":{"3":{"I":8},"4":{"F":3.9}}}`, do("/x?mea="+e(`3={"I":8}`)+"&mea="+e(`4={"F":3.9}`)))
	assert.Equal(t, `200->{"SEA":[{"I":8},{"F":3.9}]}`, do("/x?sea="+e(`{"I":8}`)+"&sea="+e(`{"F":3.9}`)))
	assert.Equal(t, `200->{"ME":{"3":4,"9":0}}`, do("/x?me=3%3D4&me=9%3D0"))
	assert.Equal(t, `200->{"SE":[3,9,2]}`, do("/x?se=3&se=9&se=2"))
	assert.Equal(t, `200->{"SA":[{"I":8},{"F":3.9}]}`, do("/x?sa="+e(`[{"I":8},{"F":3.9}]`)))
	assert.Equal(t, `200->{"MA":{"3":{"I":8},"4":{"F":3.9}}}`, do("/x?ma="+e(`{"3":{"I":8},"4":{"F":3.9}}`)))
}

func TestDecodeQueryOtherEncoders(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		XML  *thing `json:",omitempty" nvelope:"query,name=xml,explode=false,content=application/xml"`
		YAML *thing `json:",omitempty" nvelope:"query,name=yaml,explode=false,content=text/yaml"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})

	xmle := func(i interface{}) string {
		enc, err := xml.Marshal(i)
		require.NoError(t, err, "marshal xml")
		return e(string(enc))
	}
	yamle := func(i interface{}) string {
		enc, err := yaml.Marshal(i)
		require.NoError(t, err, "marshal yaml")
		return e(string(enc))
	}

	assert.Equal(t, `200->{"XML":{"I":3,"F":6.2}}`, do("/x?xml="+xmle(thing{I: 3, F: 6.2})))
	assert.Equal(t, `200->{"YAML":{"I":8,"F":2.2}}`, do("/x?yaml="+yamle(thing{I: 8, F: 2.2})))
}

func TestDecodeFormValues(t *testing.T) {
	do := captureOutput("/x", func(s struct {
		A int `json:",omitempty" nvelope:"query,name=a"`
		B int `json:",omitempty" nvelope:"query,form,name=b"`
		C int `json:",omitempty" nvelope:"query,formOnly,name=c"`
		D int `json:",omitempty" nvelope:"query,formOnly,name=d"`
	},
	) (nvelope.Response, error) {
		return s, nil
	})

	assert.Equal(t, `200->{"A":7,"B":8,"C":9}`, do("/x?a=7&b=8", header("Content-type", "application/x-www-form-urlencoded"), body(`c=9`)))
	assert.Equal(t, `200->{"A":7,"B":8}`, do("/x?a=7&b=8", header("Content-type", "application/json"), body(`{}`)))
	assert.Equal(t, `200->{"A":7,"B":8,"C":9,"D":2}`, do("/x?a=7", header("Content-type", "application/x-www-form-urlencoded"), body(`c=9&b=8&d=2`)))
}
