package tui

import "testing"

type sensitiveSample struct{}

func (sensitiveSample) Name() string     { return "secret" }
func (sensitiveSample) Label() string    { return "Secret" }
func (sensitiveSample) Value() string    { return "s3cr3t" }
func (sensitiveSample) Change(_ string)  {}
func (sensitiveSample) Sensitive() bool { return true }

func TestSensitiveIsOptionalCapability(t *testing.T) {
	_, ok := any(editSample{}).(Sensitive)
	if ok {
		t.Fatalf("editSample must NOT satisfy Sensitive — optional capability leaked to existing handler")
	}
}

func TestSensitiveDetected(t *testing.T) {
	var h any = sensitiveSample{}
	s, ok := h.(Sensitive)
	if !ok {
		t.Fatalf("sensitiveSample must satisfy Sensitive, got ok=false")
	}
	if !s.Sensitive() {
		t.Fatalf("Sensitive() = false, want true")
	}
	if _, ok := h.(HandlerEdit); !ok {
		t.Fatalf("sensitiveSample must also satisfy HandlerEdit")
	}
}
