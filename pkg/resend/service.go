package resend

import (
	"log"

	resendlib "github.com/resend/resend-go/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
)

type resendService struct {
	secret string
}

func NewResendService(cfg *config.Config) ResendService {
	return &resendService{
		secret: cfg.ResendKey,
	}
}

func (r *resendService) SendEmailVerificationCode(sendTo, code string) {
	// Send email code verification
	client := resendlib.NewClient(r.secret)
	params := &resendlib.SendEmailRequest{
		From:    "mulo@resend.dev",
		To:      []string{sendTo},
		Subject: "Mulo Email Verification",
		Html:    "<p>Your verification code is <strong>" + code + "</strong></p>",
	}

	if _, err := client.Emails.Send(params); err != nil {
		log.Printf("Resend code err: %v", err)
		return
	}
}
