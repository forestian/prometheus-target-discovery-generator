package validate

import (
	"fmt"
	"strings"

	"github.com/sdfnj/ptdgen/internal/model"
)

// Validate checks a parsed TargetFile for correctness.
// It returns a slice of human-readable error strings; an empty slice means valid.
func Validate(tf *model.TargetFile) []string {
	var errs []string

	seenNames := map[string]bool{}
	// key: "job/address:port"
	seenEndpoints := map[string]bool{}

	for i, t := range tf.Targets {
		prefix := fmt.Sprintf("target[%d]", i)
		if t.Name == "" {
			errs = append(errs, fmt.Sprintf("%s: name must not be empty", prefix))
		} else {
			key := strings.ToLower(t.Name)
			if seenNames[key] {
				errs = append(errs, fmt.Sprintf("%s: duplicate target name %q", prefix, t.Name))
			}
			seenNames[key] = true
			prefix = fmt.Sprintf("target %q", t.Name)
		}

		if t.Job == "" {
			errs = append(errs, fmt.Sprintf("%s: job must not be empty", prefix))
		}

		if t.Address == "" {
			errs = append(errs, fmt.Sprintf("%s: address must not be empty", prefix))
		}

		if t.Port < 1 || t.Port > 65535 {
			errs = append(errs, fmt.Sprintf("%s: port %d is invalid (must be 1–65535)", prefix, t.Port))
		}

		if t.Job != "" && t.Address != "" && t.Port >= 1 && t.Port <= 65535 {
			epKey := fmt.Sprintf("%s/%s:%d", t.Job, t.Address, t.Port)
			if seenEndpoints[epKey] {
				errs = append(errs, fmt.Sprintf("%s: duplicate address:port %s:%d under job %q", prefix, t.Address, t.Port, t.Job))
			}
			seenEndpoints[epKey] = true
		}
	}

	return errs
}
