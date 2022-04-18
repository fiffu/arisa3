package envconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type person struct {
	Name      string  `envvar:"NAME"`
	Age       int     `envvar:"AGE"`
	HeightM   float64 `envvar:"HEIGHT"`
	IsMortal  bool    `envvar:"MORTAL"`
	faveFruit string  `envvar:"FRUIT"` // private info
}

func Test_MergeEnvVars(t *testing.T) {
	for k, v := range map[string]string{
		"TEST_NAME":   "Socrates",
		"TEST_AGE":    "71",
		"TEST_HEIGHT": "1.92",
		"TEST_MORTAL": "1",
		"TEST_FRUIT":  "guava",
	} {
		os.Setenv(k, v)
	}
	s := &person{}
	replaced, err := MergeEnvVars(s, "TEST_")

	assert.NoError(t, err)
	assert.Equal(t, "Socrates", s.Name)
	assert.Equal(t, 71, s.Age)
	assert.Equal(t, 1.92, s.HeightM)
	assert.Equal(t, true, s.IsMortal)
	assert.Equal(t, "", s.faveFruit)

	envKeys := make([]string, 0)
	for key := range replaced {
		envKeys = append(envKeys, key)
	}
	expectKeys := []string{"TEST_NAME", "TEST_AGE", "TEST_HEIGHT", "TEST_MORTAL"}
	assert.ElementsMatch(t, expectKeys, envKeys)
}

func Test_MergeEnvVars_Asserts(t *testing.T) {
	type Stru struct{}

	var (
		a *string // pointer to wrong type
		b string  // wrong non-pointer
		c Stru    // non-pointer to correct type
		d *Stru   // nil pointer to correct type
	)

	for _, each := range []interface{}{a, b, c, d} {
		_, err := MergeEnvVars(each, "")
		assert.Error(t, err)
	}
}
