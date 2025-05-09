package resend

type ResendService interface {
	SendEmailVerificationCode(sendTo, code string)
}
