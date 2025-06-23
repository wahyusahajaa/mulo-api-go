package resend

import (
	"context"

	resendlib "github.com/resend/resend-go/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type resendService struct {
	secret string
	log    *logrus.Logger
}

func NewResendService(cfg *config.Config, log *logrus.Logger) ResendService {
	return &resendService{
		secret: cfg.ResendKey,
		log:    log,
	}
}

func (r *resendService) SendEmailVerificationCode(sendTo, code string) {
	// Send email code verification
	client := resendlib.NewClient(r.secret)
	params := &resendlib.SendEmailRequest{
		From:    "noreply@craftedfolio.my.id",
		To:      []string{sendTo},
		Subject: "Mulo Email Verification",
		Html:    "<p>Your verification code is <strong>" + code + "</strong></p>",
	}

	if _, err := client.Emails.Send(params); err != nil {
		utils.LogError(r.log, context.Background(), "resend_service", "SendEmailVerificationCode", err)
		return
	}
}
