package pkg

import (
	"os"
	"testing"

	"github.com/miniscruff/envexample/then"
)

func Test_Config_BadArgs(t *testing.T) {
	args := []string{"-invalid-arg", "value"}
	_, err := NewConfig(args)
	then.Err(t, errUnableToParseFlags, err)
}

func Test_Config_Defaults(t *testing.T) {
	cfg, err := NewConfig([]string{})
	then.Nil(t, err)
	then.Equals(t, ".env.example", cfg.ExportFile)
	then.Equals(t, "", cfg.ConfigStruct)
	then.Equals(t, "", cfg.Package)
	then.False(t, cfg.DryRun)
	then.False(t, cfg.ShowVersion)
	then.False(t, cfg.ShowHelp)
}

func Test_Config_OverrideDefaults(t *testing.T) {
	args := []string{
		"-export", "test.env.example",
		"-type", "myStruct",
		"-pkg", "github.com/your/project/config",
		"-prefix", "C_",
		"-tag", "cenv",
		"-required-if-no-def",
		"-use-field-name",
		"-dry", "-v", "-h",
	}

	cfg, err := NewConfig(args)
	then.Nil(t, err)
	then.Equals(t, "test.env.example", cfg.ExportFile)
	then.Equals(t, "myStruct", cfg.ConfigStruct)
	then.Equals(t, "github.com/your/project/config", cfg.Package)
	then.Equals(t, "cenv", cfg.TagName)
	then.Equals(t, "C_", cfg.Prefix)
	then.True(t, cfg.DryRun)
	then.True(t, cfg.ShowVersion)
	then.True(t, cfg.ShowHelp)
	then.True(t, cfg.RequiredIfNoDef)
	then.True(t, cfg.UseFieldNameByDefault)
}

func Test_Config_ValidArgs(t *testing.T) {
	args := []string{
		"-type", "MyConfig",
		"-pkg", "github.com/your/project/config",
	}

	cfg, err := NewConfig(args)
	then.Nil(t, err)
	then.Nil(t, cfg.Validate())
}

func Test_Config_InvalidArgsNoStruct(t *testing.T) {
	cfg, err := NewConfig([]string{})
	then.Nil(t, err)
	then.Err(t, errInvalidConfigNoStruct, cfg.Validate())
}

func Test_Config_InvalidArgsNoExport(t *testing.T) {
	args := []string{
		"-type", "MyConfigStruct",
		"-export", "",
	}

	cfg, err := NewConfig(args)
	then.Nil(t, err)
	then.Err(t, errInvalidConfigNoExport, cfg.Validate())
}

func Test_Config_StdoutOnDryRun(t *testing.T) {
	cfg := &Config{
		DryRun: true,
	}
	writer, err := cfg.Writer()
	then.Nil(t, err)
	then.Equals(t, os.Stdout, writer.(*os.File))
}

func Test_Config_ExportFileWriter(t *testing.T) {
	dir := t.TempDir()

	cfg := &Config{
		ExportFile: dir + "/.env.example",
	}
	writer, err := cfg.Writer()
	then.Nil(t, err)

	defer writer.Close()

	exportFile, ok := writer.(*os.File)
	then.True(t, ok)
	then.Equals(t, cfg.ExportFile, exportFile.Name())
}

func Test_Config_ExportFileError(t *testing.T) {
	dir := t.TempDir()

	cfg := &Config{
		ExportFile: dir + "/.",
	}
	_, err := cfg.Writer()
	then.Err(t, errUnableToCreateWriter, err)
}
