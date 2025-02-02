package users

import (
	"testing"
)

func TestVerifyEmail(t *testing.T) {
	t.Run("Empty email", func(t *testing.T) {
		if ok := VerifyEmail(""); ok {
			t.Error("Failed to recognize empty email")
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		if ok := VerifyEmail("skh@.dj"); ok {
			t.Error("Failed to recognize invalid email")
		}
	})

	t.Run("Valid email", func(t *testing.T) {
		if ok := VerifyEmail("sjkh@kshj.com"); !ok {
			t.Error("Failed to verify email")
		}
	})
}