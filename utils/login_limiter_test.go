package utils

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

type MockCaptchaProvider struct{}

func (p *MockCaptchaProvider) Generate() (string, string, string, error) {
	id := uuid.New().String()
	content := uuid.New().String()
	answer := uuid.New().String()
	return id, content, answer, nil
}

func (p *MockCaptchaProvider) Expiration() time.Duration {
	return 2 * time.Second
}
func (p *MockCaptchaProvider) Draw(content string) (string, error) {
	return "MOCK", nil
}

func TestSecurityWorkflow(t *testing.T) {
	policy := SecurityPolicy{
		CaptchaThreshold: 3,
		BanThreshold:     5,
		AttemptsWindow:   5 * time.Minute,
		BanDuration:      5 * time.Minute,
	}
	limiter := NewLoginLimiter(policy)
	ip := "192.168.1.100"

	// Test normal failed-attempt recording
	for i := 0; i < 3; i++ {
		limiter.RecordFailedAttempt(ip)
	}
	isBanned, capRequired := limiter.CheckSecurityStatus(ip)
	fmt.Printf("IP: %s, Banned: %v, Captcha Required: %v\n", ip, isBanned, capRequired)
	if isBanned {
		t.Error("IP should not be banned yet")
	}
	if !capRequired {
		t.Error("Captcha should be required")
	}
	// Test triggering a ban
	for i := 0; i < 3; i++ {
		limiter.RecordFailedAttempt(ip)
		isBanned, capRequired = limiter.CheckSecurityStatus(ip)
		fmt.Printf("IP: %s, Banned: %v, Captcha Required: %v\n", ip, isBanned, capRequired)
	}

	// Test ban status
	if isBanned, _ = limiter.CheckSecurityStatus(ip); !isBanned {
		t.Error("IP should be banned")
	}
}

func TestCaptchaFlow(t *testing.T) {
	policy := SecurityPolicy{CaptchaThreshold: 2}
	limiter := NewLoginLimiter(policy)
	limiter.RegisterProvider(&MockCaptchaProvider{})
	ip := "10.0.0.1"

	// Trigger captcha requirement
	limiter.RecordFailedAttempt(ip)
	limiter.RecordFailedAttempt(ip)

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should be required")
	}

	// Generate captcha
	err, capc := limiter.RequireCaptcha()
	if err != nil {
		t.Fatalf("captcha generation failed: %v", err)
	}
	fmt.Printf("captcha content: %#v\n", capc)

	// Verify success
	if !limiter.VerifyCaptcha(capc.Id, capc.Answer) {
		t.Error("captcha verification should succeed")
	}

	// Captcha should be deleted after use
	if limiter.VerifyCaptcha(capc.Id, capc.Answer) {
		t.Error("captcha should have been deleted")
	}

	limiter.RemoveAttempts(ip)
	// Status after verification
	if banned, need := limiter.CheckSecurityStatus(ip); banned || need {
		t.Error("status should be reset after successful verification")
	}
}

func TestCaptchaMustFlow(t *testing.T) {
	policy := SecurityPolicy{CaptchaThreshold: 0}
	limiter := NewLoginLimiter(policy)
	limiter.RegisterProvider(&MockCaptchaProvider{})
	ip := "10.0.0.1"

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should be required")
	}

	// Generate captcha
	err, capc := limiter.RequireCaptcha()
	if err != nil {
		t.Fatalf("captcha generation failed: %v", err)
	}
	fmt.Printf("captcha content: %#v\n", capc)

	// Verify success
	if !limiter.VerifyCaptcha(capc.Id, capc.Answer) {
		t.Error("captcha verification should succeed")
	}

	// Status after verification
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should still be required")
	}
}
func TestAttemptTimeout(t *testing.T) {
	policy := SecurityPolicy{CaptchaThreshold: 2, AttemptsWindow: 1 * time.Second}
	limiter := NewLoginLimiter(policy)
	limiter.RegisterProvider(&MockCaptchaProvider{})
	ip := "10.0.0.1"

	// Trigger captcha requirement
	limiter.RecordFailedAttempt(ip)
	limiter.RecordFailedAttempt(ip)

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should be required")
	}

	// Generate captcha
	err, _ := limiter.RequireCaptcha()
	if err != nil {
		t.Fatalf("captcha generation failed: %v", err)
	}
	// Wait until the AttemptsWindow expires
	time.Sleep(2 * time.Second)
	// Trigger captcha check again
	limiter.RecordFailedAttempt(ip)

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); need {
		t.Error("captcha should not be required")
	}
}

func TestCaptchaTimeout(t *testing.T) {
	policy := SecurityPolicy{CaptchaThreshold: 2}
	limiter := NewLoginLimiter(policy)
	limiter.RegisterProvider(&MockCaptchaProvider{})
	ip := "10.0.0.1"

	// Trigger captcha requirement
	limiter.RecordFailedAttempt(ip)
	limiter.RecordFailedAttempt(ip)

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should be required")
	}

	// Generate captcha
	err, capc := limiter.RequireCaptcha()
	if err != nil {
		t.Fatalf("captcha generation failed: %v", err)
	}

	// Wait until the CaptchaValidPeriod expires
	time.Sleep(3 * time.Second)

	// Captcha should be expired
	if limiter.VerifyCaptcha(capc.Id, capc.Answer) {
		t.Error("captcha should have expired")
	}

}

func TestBanFlow(t *testing.T) {
	policy := SecurityPolicy{BanThreshold: 5}
	limiter := NewLoginLimiter(policy)
	ip := "10.0.0.1"
	// Trigger ban
	for i := 0; i < 5; i++ {
		limiter.RecordFailedAttempt(ip)
	}

	// Check status
	if banned, _ := limiter.CheckSecurityStatus(ip); !banned {
		t.Error("should be banned")
	}
}
func TestBanDisableFlow(t *testing.T) {
	policy := SecurityPolicy{BanThreshold: 0}
	limiter := NewLoginLimiter(policy)
	ip := "10.0.0.1"
	// Trigger ban
	for i := 0; i < 5; i++ {
		limiter.RecordFailedAttempt(ip)
	}

	// Check status
	if banned, _ := limiter.CheckSecurityStatus(ip); banned {
		t.Error("should not be banned")
	}
}
func TestBanTimeout(t *testing.T) {
	policy := SecurityPolicy{BanThreshold: 5, BanDuration: 1 * time.Second}
	limiter := NewLoginLimiter(policy)
	ip := "10.0.0.1"
	// Trigger ban
	for i := 0; i < 5; i++ {
		limiter.RecordFailedAttempt(ip)
	}

	time.Sleep(2 * time.Second)

	// Check status
	if banned, _ := limiter.CheckSecurityStatus(ip); banned {
		t.Error("should not be banned")
	}
}

func TestLimiterDisabled(t *testing.T) {
	policy := SecurityPolicy{BanThreshold: 0, CaptchaThreshold: -1}
	limiter := NewLoginLimiter(policy)
	ip := "10.0.0.1"
	// Trigger ban
	for i := 0; i < 5; i++ {
		limiter.RecordFailedAttempt(ip)
	}

	// Check status
	if banned, capNeed := limiter.CheckSecurityStatus(ip); banned || capNeed {
		fmt.Printf("IP: %s, Banned: %v, Captcha Required: %v\n", ip, banned, capNeed)
		t.Error("should not be banned or need captcha")
	}
}

func TestB64CaptchaFlow(t *testing.T) {
	limiter := NewLoginLimiter(defaultSecurityPolicy)
	limiter.RegisterProvider(B64StringCaptchaProvider{})
	ip := "10.0.0.1"

	// Trigger captcha requirement
	limiter.RecordFailedAttempt(ip)
	limiter.RecordFailedAttempt(ip)
	limiter.RecordFailedAttempt(ip)

	// Check status
	if _, need := limiter.CheckSecurityStatus(ip); !need {
		t.Error("captcha should be required")
	}

	// Generate captcha
	err, capc := limiter.RequireCaptcha()
	if err != nil {
		t.Fatalf("captcha generation failed: %v", err)
	}
	fmt.Printf("captcha content: %#v\n", capc)

	// Draw captcha
	err, b64 := limiter.DrawCaptcha(capc.Content)
	if err != nil {
		t.Fatalf("captcha rendering failed: %v", err)
	}
	fmt.Printf("captcha image: %#v\n", b64)

	// Verify success
	if !limiter.VerifyCaptcha(capc.Id, capc.Answer) {
		t.Error("captcha verification should succeed")
	}
	limiter.RemoveAttempts(ip)
	// Status after verification
	if banned, need := limiter.CheckSecurityStatus(ip); banned || need {
		t.Error("status should be reset after successful verification")
	}
}
