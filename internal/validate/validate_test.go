package validate_test

import (
	"testing"

	"github.com/sdfnj/ptdgen/internal/model"
	"github.com/sdfnj/ptdgen/internal/validate"
)

func makeTarget(name, job, addr string, port int) model.Target {
	t := model.Target{
		Name:    name,
		Job:     job,
		Address: addr,
		Port:    port,
	}
	t.ApplyDefaults()
	return t
}

func TestValidateOK(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{
			makeTarget("t1", "job1", "10.0.0.1", 9100),
			makeTarget("t2", "job1", "10.0.0.2", 9100),
		},
	}
	errs := validate.Validate(tf)
	if len(errs) != 0 {
		t.Errorf("want no errors, got %v", errs)
	}
}

func TestValidateEmptyName(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{makeTarget("", "job1", "10.0.0.1", 9100)},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for empty name")
	}
}

func TestValidateEmptyJob(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{makeTarget("t1", "", "10.0.0.1", 9100)},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for empty job")
	}
}

func TestValidateEmptyAddress(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{makeTarget("t1", "job1", "", 9100)},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for empty address")
	}
}

func TestValidateInvalidPortZero(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{makeTarget("t1", "job1", "10.0.0.1", 0)},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for port=0")
	}
}

func TestValidateInvalidPortTooHigh(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{makeTarget("t1", "job1", "10.0.0.1", 70000)},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for port=70000")
	}
}

func TestValidatePortBoundary(t *testing.T) {
	for _, port := range []int{1, 65535} {
		tf := &model.TargetFile{
			Targets: []model.Target{makeTarget("t1", "job1", "10.0.0.1", port)},
		}
		errs := validate.Validate(tf)
		if len(errs) != 0 {
			t.Errorf("port=%d should be valid, got errors %v", port, errs)
		}
	}
}

func TestValidateDuplicateName(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{
			makeTarget("t1", "job1", "10.0.0.1", 9100),
			makeTarget("t1", "job2", "10.0.0.2", 9100),
		},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for duplicate name")
	}
}

func TestValidateDuplicateEndpointSameJob(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{
			makeTarget("t1", "job1", "10.0.0.1", 9100),
			makeTarget("t2", "job1", "10.0.0.1", 9100),
		},
	}
	errs := validate.Validate(tf)
	if len(errs) == 0 {
		t.Error("want error for duplicate address:port under same job")
	}
}

func TestValidateDuplicateEndpointDifferentJob(t *testing.T) {
	// Same address:port under different jobs is allowed.
	tf := &model.TargetFile{
		Targets: []model.Target{
			makeTarget("t1", "job1", "10.0.0.1", 9100),
			makeTarget("t2", "job2", "10.0.0.1", 9100),
		},
	}
	errs := validate.Validate(tf)
	if len(errs) != 0 {
		t.Errorf("same address:port under different jobs should be OK, got %v", errs)
	}
}

func TestValidateMultipleErrors(t *testing.T) {
	tf := &model.TargetFile{
		Targets: []model.Target{
			makeTarget("", "", "", 0),
		},
	}
	errs := validate.Validate(tf)
	if len(errs) < 3 {
		t.Errorf("want at least 3 errors (name, job, address), got %v", errs)
	}
}
