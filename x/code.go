package x

// GenerateVerificationCode generates a verification code
func GenerateVerificationCode() string {
	return verificationToken()
}

// GeneratePasswordResetCode generates a password reset code
func GeneratePasswordResetCode() string {
	return verificationToken()
}
