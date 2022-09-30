// Copyright 2022 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	TwoFactorSessionUserId        = "TwoFactorSessionUserId"
	TwoFactorSessionEnableSession = "TwoFactorSessionEnableSession"
	TwoFactorSessionAutoSignIn    = "TwoFactorSessionAutoSignIn"
	NextTwoFactor                 = "nextTwoFactor"
)

type TwoFactorSessionData struct {
	UserId        string
	EnableSession bool
	AutoSignIn    bool
}

func NewTotpKey(issuer string, accountName string) (*otp.Key, error) {
	period := 30
	secretSize := 20
	digits := otp.DigitsSix
	return totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Period:      uint(period),
		SecretSize:  uint(secretSize),
		Digits:      digits,
	})
}

func ValidateTotpPassCode(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}

func RecoverTotp(user *User, recoveryCode string) (bool, error) {
	if user.TwoFactorRecoveryCode != recoveryCode {
		return false, nil
	}
	user.TwoFactorRecoveryCode = ""
	user.TotpSecret = ""
	affected, err := adapter.Engine.
		Cols("totp_secret", "two_factor_recovery_code").Update(user)
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}
