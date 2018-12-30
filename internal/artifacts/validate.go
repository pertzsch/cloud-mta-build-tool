package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/validations"
)

// ExecuteValidation - executes validation of MTA
func ExecuteValidation(source, desc, mode string, getWorkingDir func() (string, error)) error {
	logs.Logger.Info("validation started")
	loc, err := dir.Location(source, "", desc, getWorkingDir)
	if err != nil {
		return errors.Wrap(err, "validation failed when initializing location")
	}
	validateSchema, validateProject, err := validate.GetValidationMode(mode)
	if err != nil {
		return errors.Wrap(err, "validation failed when analyzing the validation mode")
	}
	err = validate.MtaYaml(source, loc.GetMtaYamlFilename(), validateSchema, validateProject)
	if err != nil {
		return err
	}
	logs.Logger.Info("validation finished successfully")
	return nil
}