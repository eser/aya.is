package configfx_test

import (
	"reflect"
	"testing"

	"github.com/eser/aya.is/services/pkg/ajan/configfx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Host string `conf:"host" default:"localhost"`
}

type TestConfigNestedKV struct {
	Name string `conf:"name"`
}

type TestConfigNested struct {
	TestConfig

	Port     int    `conf:"port"      default:"8080"`
	MaxRetry uint16 `conf:"max_retry" default:"10"`

	Dictionary map[string]string    `conf:"dict"`
	Array      []TestConfigNestedKV `conf:"arr"`
}

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("should load config", func(t *testing.T) {
		t.Parallel()

		config := TestConfigNested{} //nolint:exhaustruct

		cl := configfx.NewConfigManager()
		err := cl.Load(&config)

		require.NoError(t, err)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 8080, config.Port)
		assert.Equal(t, uint16(10), config.MaxRetry)
	})

	t.Run("should load config from string", func(t *testing.T) {
		t.Parallel()

		config := TestConfigNested{} //nolint:exhaustruct

		cl := configfx.NewConfigManager()
		err := cl.Load(
			&config,
			cl.FromJSONFile("testdata/config.json"),
			cl.FromEnvFile("testdata/.env", true),
		)

		require.NoError(t, err)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 8081, config.Port)
		assert.Equal(t, uint16(20), config.MaxRetry)
		assert.Equal(
			t,
			map[string]string{"key": "value", "key2": "value2", "key3": "value3"},
			config.Dictionary,
		)
		// assert.Equal(t, []TestConfigNestedKV{{Name: "eser"}}, config.Array)
	})
}

func TestLoadMeta(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("should get config meta", func(t *testing.T) {
		t.Parallel()

		config := TestConfig{} //nolint:exhaustruct

		cl := configfx.NewConfigManager()
		meta, err := cl.LoadMeta(&config)

		expected := []configfx.ConfigItemMeta{
			{
				Name:            "host",
				Field:           meta.Children[0].Field,
				Type:            reflect.TypeFor[string](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "localhost",

				Children: nil,
			},
		}

		require.NoError(t, err)
		assert.Equal(t, "root", meta.Name)
		assert.Nil(t, meta.Type)

		assert.ElementsMatch(t, expected, meta.Children)
	})

	t.Run("should get config meta from nested definition", func(t *testing.T) {
		t.Parallel()

		config := TestConfigNested{} //nolint:exhaustruct

		cl := configfx.NewConfigManager()
		meta, err := cl.LoadMeta(&config)

		expected := []configfx.ConfigItemMeta{
			{
				Name:            "host",
				Field:           meta.Children[0].Field,
				Type:            reflect.TypeFor[string](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "localhost",

				Children: nil,
			},
			{
				Name:            "port",
				Field:           meta.Children[1].Field,
				Type:            reflect.TypeFor[int](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "8080",

				Children: nil,
			},
			{
				Name:            "max_retry",
				Field:           meta.Children[2].Field,
				Type:            reflect.TypeFor[uint16](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "10",

				Children: nil,
			},
			{
				Name:            "dict",
				Field:           meta.Children[3].Field,
				Type:            reflect.TypeFor[map[string]string](),
				IsRequired:      false,
				HasDefaultValue: false,
				DefaultValue:    "",

				Children: nil,
			},
			{
				Name:            "arr",
				Field:           meta.Children[4].Field,
				Type:            reflect.TypeFor[[]TestConfigNestedKV](),
				IsRequired:      false,
				HasDefaultValue: false,
				DefaultValue:    "",

				Children: nil,
			},
		}

		require.NoError(t, err)
		assert.Equal(t, "root", meta.Name)
		assert.Nil(t, meta.Type)

		assert.ElementsMatch(t, expected, meta.Children)
	})
}
