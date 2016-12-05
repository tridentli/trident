package trident

import (
	pf "trident.li/pitchfork/lib"
)

const CRLF = pf.CRLF

func Mail_PassResetUser(ctx pf.PfCtx, email pf.PfUserEmail, is_reset bool, nom_email pf.PfUserEmail, user_portion string) (err error) {
	sys := pf.System_Get()
	subject := ""
	body := "Dear " + email.FullName + "," + CRLF +
		CRLF

	if is_reset {
		subject = "Password Reset (User portion)"
		body += "A password reset request was made." + CRLF
	} else {
		subject = "New Account (User portion)"
		body += "Your account has been created. A new password thus has to be set." + CRLF
	}

	body += CRLF +
		"We are therefor sending you two token portions." + CRLF +
		"The user portion is in this email, the other portion " + CRLF +
		"has been sent to your nominator who will forward it in " + CRLF +
		"a secure method towards you." + CRLF +
		CRLF +
		"Your nominator is:" + CRLF +
		" " + nom_email.FullName + " <" + nom_email.Email + ">" + CRLF +
		CRLF +
		"When both parts have been received by you, please proceed to:" + CRLF +
		"  " + sys.PublicURL + "/recover/" +
		CRLF +
		"and enter the following password in the User Portion:" + CRLF +
		"  " + user_portion + CRLF +
		CRLF +
		"If you do not perform this reset the request will be canceled." + CRLF

	err = pf.Mail(ctx,
		"", "",
		email.FullName, email.Email,
		true,
		subject,
		body,
		true,
		"",
		true)

	return
}

func Mail_PassResetNominator(ctx pf.PfCtx, email pf.PfUserEmail, is_reset bool, user_email pf.PfUserEmail, nom_portion string) (err error) {
	subject := ""
	body := "Dear " + email.FullName + "," + CRLF +
		CRLF

	if is_reset {
		subject = "Password Reset (Nominator portion)"
		body += "A password reset request was made for:" + CRLF
	} else {
		subject = "New account (Nominator portion)"
		body += "A new account has been created for: " + CRLF
	}

	body += " " + user_email.FullName + " <" + user_email.Email + ">" + CRLF +
		CRLF +
		"As you are a nominator of this person, you are receiving " + CRLF +
		"the second portion of this email. " + CRLF +
		CRLF +
		"Please securely inform " + user_email.FullName + CRLF +
		"of the following Nominator Portion of the password reset: " + CRLF +
		"  " + nom_portion + CRLF

	err = pf.Mail(ctx,
		"", "",
		email.FullName, email.Email,
		true,
		subject,
		body,
		true,
		"",
		true)

	return
}
