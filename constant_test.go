package charlatan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestConst(t *testing.T, v interface{}) *Const {
	c, err := NewConst(v)
	assert.Nil(t, err)
	require.NotNil(t, c)
	return c
}

func TestNewConstNull(t *testing.T) {
	assert.True(t, makeTestConst(t, nil).IsNull())
}

func TestNewConstInt(t *testing.T) {
	c := makeTestConst(t, 42)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())

	c = makeTestConst(t, int8(42))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())

	c = makeTestConst(t, int16(42))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())

	c = makeTestConst(t, int32(42))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())

	c = makeTestConst(t, int64(42))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())
}

func TestNewConstIntPtr(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(42)

	c := makeTestConst(t, &i)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, i64, c.AsInt())

	c = makeTestConst(t, &i8)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, i64, c.AsInt())

	c = makeTestConst(t, &i16)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, i64, c.AsInt())

	c = makeTestConst(t, &i32)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, i64, c.AsInt())

	c = makeTestConst(t, &i64)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, i64, c.AsInt())
}

func TestNewConstFloat(t *testing.T) {
	c := makeTestConst(t, float32(42.0))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, float64(42.0), c.AsFloat())

	c = makeTestConst(t, float64(42.0))
	assert.True(t, c.IsNumeric())
	assert.Equal(t, float64(42.0), c.AsFloat())
}

func TestNewConstFloatPtr(t *testing.T) {
	f32 := float32(42.0)
	f64 := float64(42.0)

	c := makeTestConst(t, &f32)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, f64, c.AsFloat())

	c = makeTestConst(t, &f64)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, f64, c.AsFloat())
}

func TestNewConstBool(t *testing.T) {
	c := makeTestConst(t, true)
	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())

	c = makeTestConst(t, false)
	assert.True(t, c.IsBool())
	assert.Equal(t, false, c.AsBool())
}

func TestNewConstBoolPtr(t *testing.T) {
	btrue := true
	bfalse := false

	c := makeTestConst(t, &btrue)
	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())

	c = makeTestConst(t, &bfalse)
	assert.True(t, c.IsBool())
	assert.Equal(t, false, c.AsBool())
}

func TestNewConstString(t *testing.T) {
	c := makeTestConst(t, "yolo")
	assert.True(t, c.IsString())
	assert.Equal(t, "yolo", c.AsString())
}

func TestNewConstStringPtr(t *testing.T) {
	s := "yolo"
	c := makeTestConst(t, &s)
	assert.True(t, c.IsString())
	assert.Equal(t, s, c.AsString())
}

func TestNullConst(t *testing.T) {
	assert.True(t, NullConst().IsNull())
}

func TestIntConst(t *testing.T) {
	c := IntConst(42)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())
}

func TestFloatConst(t *testing.T) {
	c := FloatConst(42)
	assert.True(t, c.IsNumeric())
	assert.Equal(t, float64(42), c.AsFloat())
}

func TestBoolConst(t *testing.T) {
	c := BoolConst(true)
	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())
}

func TestStringConst(t *testing.T) {
	c := StringConst("yo")
	assert.True(t, c.IsString())
	assert.Equal(t, "yo", c.AsString())
}

func TestConstFromStringInt(t *testing.T) {
	c := ConstFromString("42")
	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(42), c.AsInt())
}

func TestConstFromStringFloat(t *testing.T) {
	c := ConstFromString("42.0")
	assert.True(t, c.IsNumeric())
	assert.Equal(t, float64(42.0), c.AsFloat())
}

func TestConstFromStringBool(t *testing.T) {
	c := ConstFromString("true")
	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())
}

func TestConstFromStringString(t *testing.T) {
	c := ConstFromString("[]")
	assert.True(t, c.IsString())
	assert.Equal(t, "[]", c.AsString())
}

func TestConstGetType(t *testing.T) {
	assert.Equal(t, constNull, NullConst().GetType())
	assert.Equal(t, constInt, IntConst(42).GetType())
	assert.Equal(t, constFloat, FloatConst(1.0).GetType())
}

func TestConstValue(t *testing.T) {
	assert.Nil(t, NullConst().Value())

	{
		v, ok := IntConst(41).Value().(int64)
		assert.True(t, ok)
		assert.Equal(t, int64(41), v)
	}

	{
		v, ok := FloatConst(41).Value().(float64)
		assert.True(t, ok)
		assert.Equal(t, float64(41), v)
	}

	{
		v, ok := BoolConst(true).Value().(bool)
		assert.True(t, ok)
		assert.Equal(t, true, v)
	}

	{
		v, ok := StringConst("ya").Value().(string)
		assert.True(t, ok)
		assert.Equal(t, "ya", v)
	}
}

func TestConstString(t *testing.T) {
	assert.Equal(t, "42", IntConst(42).String())
}

func TestConstAsFloat(t *testing.T) {
	assert.Equal(t, float64(3.14), FloatConst(3.14).AsFloat())
	assert.Equal(t, float64(3), IntConst(3).AsFloat())
	assert.Equal(t, float64(1), BoolConst(true).AsFloat())
	assert.Equal(t, float64(0), BoolConst(false).AsFloat())
	assert.Equal(t, float64(0), StringConst("yo").AsFloat())
	assert.Equal(t, float64(0), NullConst().AsFloat())
}

func TestConstAsInt(t *testing.T) {
	assert.Equal(t, int64(3), FloatConst(3.14).AsInt())
	assert.Equal(t, int64(3), IntConst(3).AsInt())
	assert.Equal(t, int64(1), BoolConst(true).AsInt())
	assert.Equal(t, int64(0), BoolConst(false).AsInt())
	assert.Equal(t, int64(0), StringConst("yo").AsInt())
	assert.Equal(t, int64(0), NullConst().AsInt())
}

func TestConstAsBool(t *testing.T) {
	assert.Equal(t, true, FloatConst(3.14).AsBool())
	assert.Equal(t, false, FloatConst(0).AsBool())
	assert.Equal(t, true, IntConst(3).AsBool())
	assert.Equal(t, false, IntConst(0).AsBool())
	assert.Equal(t, true, BoolConst(true).AsBool())
	assert.Equal(t, false, BoolConst(false).AsBool())
	assert.Equal(t, true, StringConst("yo").AsBool())
	assert.Equal(t, true, StringConst("").AsBool())
	assert.Equal(t, false, NullConst().AsBool())
}

func TestConstAsString(t *testing.T) {
	assert.Equal(t, "3.14", FloatConst(3.14).AsString())
	assert.Equal(t, "3", IntConst(3).AsString())
	assert.Equal(t, "0", IntConst(0).AsString())
	assert.Equal(t, "true", BoolConst(true).AsString())
	assert.Equal(t, "false", BoolConst(false).AsString())
	assert.Equal(t, "yo", StringConst("yo").AsString())
	assert.Equal(t, "", StringConst("").AsString())
	assert.Equal(t, "null", NullConst().AsString())
}

func testCmpConsts(t *testing.T, c1, c2 *Const) int {
	i, err := c1.CompareTo(c2)
	assert.Nil(t, err)
	return i
}

func TestConstCompareToSameTypes(t *testing.T) {
	assert.Equal(t, 0, testCmpConsts(t, NullConst(), NullConst()))
	assert.Equal(t, 0, testCmpConsts(t, IntConst(42), IntConst(42)))
	assert.Equal(t, 0, testCmpConsts(t, FloatConst(42), FloatConst(42)))
	assert.Equal(t, 0, testCmpConsts(t, BoolConst(true), BoolConst(true)))
	assert.Equal(t, 0, testCmpConsts(t, BoolConst(false), BoolConst(false)))
	assert.Equal(t, 0, testCmpConsts(t, StringConst("tutu"), StringConst("tutu")))
	assert.Equal(t, 0, testCmpConsts(t, StringConst(""), StringConst("")))

	// is it inversed?
	assert.True(t, 0 > testCmpConsts(t, IntConst(41), IntConst(42)))
	assert.True(t, 0 > testCmpConsts(t, FloatConst(41), FloatConst(42)))
	assert.True(t, 0 > testCmpConsts(t, BoolConst(false), BoolConst(true)))
	assert.True(t, 0 > testCmpConsts(t, StringConst("tatu"), StringConst("tutu")))

	assert.True(t, 0 < testCmpConsts(t, IntConst(43), IntConst(42)))
	assert.True(t, 0 < testCmpConsts(t, FloatConst(100), FloatConst(42)))
	assert.True(t, 0 < testCmpConsts(t, BoolConst(true), BoolConst(false)))
	assert.True(t, 0 < testCmpConsts(t, StringConst("tytu"), StringConst("tutu")))
}

func TestConstCompareToDifferentTypes(t *testing.T) {
	assert.Equal(t, 0, testCmpConsts(t, FloatConst(42.0), IntConst(42)))
	assert.Equal(t, 0, testCmpConsts(t, IntConst(42), FloatConst(42.0)))

	assert.True(t, 0 < testCmpConsts(t, FloatConst(42.0), IntConst(18)))
	assert.True(t, 0 < testCmpConsts(t, IntConst(42), FloatConst(18.0)))

	assert.True(t, 0 > testCmpConsts(t, FloatConst(2.0), IntConst(18)))
	assert.True(t, 0 > testCmpConsts(t, IntConst(2), FloatConst(18.0)))
}

func TestConstTypeString(t *testing.T) {
	for _, ty := range []constType{
		constNull, constInt, constFloat, constBool, constString,
	} {
		assert.NotEqual(t, "", ty.String())
		assert.NotEqual(t, "UNDEFINED", ty.String())
	}
}
