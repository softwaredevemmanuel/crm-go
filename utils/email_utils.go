// utils/email_utils.go
package utils

import (
	"fmt"
	"time"
)

// BuildVerificationEmail builds the HTML email body for email verification
func BuildVerificationEmail(firstName, verificationURL, token, email string) string {
	now := time.Now()
	year := now.Year()

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            line-height: 1.6;
            color: #1a202c;
            background: #f7fafc;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.08);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #134574 0%%, #1a5a9e 100%%);
            padding: 40px 30px 30px;
            text-align: center;
        }
        .header h1 { color: #ffffff; font-size: 28px; font-weight: 700; margin-bottom: 5px; }
        .header p { color: rgba(255,255,255,0.8); font-size: 16px; }
        .content { padding: 40px 35px 35px; }
        .greeting { font-size: 18px; font-weight: 600; color: #134574; margin-bottom: 15px; }
        .greeting span { color: #F13178; }
        .message { color: #4a5568; font-size: 16px; margin: 10px 0 20px; line-height: 1.8; }
        .verification-box {
            background: linear-gradient(135deg, #f7fafc 0%%, #edf2f7 100%%);
            border-radius: 12px;
            padding: 30px;
            margin: 25px 0;
            text-align: center;
            border: 2px dashed #e2e8f0;
        }
        .verification-box .icon { font-size: 48px; margin-bottom: 15px; display: block; }
        .verification-box h3 { color: #134574; font-size: 18px; font-weight: 600; margin-bottom: 10px; }
        .btn {
            display: inline-block;
            padding: 14px 40px;
            background: linear-gradient(135deg, #F13178 0%%, #d42a6a 100%%);
            color: #ffffff;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
            box-shadow: 0 4px 15px rgba(241,49,120,0.3);
            transition: all 0.3s ease;
        }
        .btn:hover { transform: translateY(-2px); box-shadow: 0 8px 25px rgba(241,49,120,0.4); }
        .divider { border: none; height: 1px; background: linear-gradient(90deg, transparent, #e2e8f0, transparent); margin: 25px 0; }
        .token-box {
            background: #f7fafc;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            padding: 12px 16px;
            font-family: 'Courier New', monospace;
            font-size: 13px;
            color: #2d3748;
            word-break: break-all;
            margin-top: 8px;
        }
        .warning-box {
            background: #fefcbf;
            border: 1px solid #f6e05e;
            border-radius: 8px;
            padding: 15px 20px;
            margin: 20px 0;
            display: flex;
            align-items: flex-start;
            gap: 12px;
        }
        .warning-box .text { font-size: 14px; color: #744210; line-height: 1.6; }
        .footer {
            background: #f7fafc;
            padding: 30px 35px;
            text-align: center;
            border-top: 1px solid #e2e8f0;
        }
        .footer p { color: #718096; font-size: 13px; line-height: 1.8; }
        .footer .copyright { margin-top: 15px; font-size: 12px; color: #a0aec0; }
        @media (max-width: 480px) {
            .header { padding: 30px 20px 25px; }
            .header h1 { font-size: 22px; }
            .content { padding: 25px 20px; }
            .verification-box { padding: 20px; }
            .btn { display: block; text-align: center; padding: 14px 20px; }
            .footer { padding: 20px 20px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email</h1>
            <p>Confirm your email address to get started</p>
        </div>
        <div class="content">
            <div class="greeting">Hello <span>%s</span>! 👋</div>
            <p class="message">
                Thanks for creating an account with <strong>EduTech</strong>. 
                Please verify your email address to unlock all features and start learning.
            </p>
            <div class="verification-box">
                <span class="icon">📧</span>
                <h3>Verify Your Email</h3>
                <p>Click the button below to verify your email address</p>
                <a href="%s" class="btn">Verify Email →</a>
            </div>
            <hr class="divider">
            <div style="margin: 20px 0 10px;">
                <div style="font-size:13px;color:#718096;font-weight:500;text-transform:uppercase;letter-spacing:0.5px;">Or use this link</div>
                <div class="token-box">%s</div>
            </div>
            <div class="warning-box">
                <div class="text">
                    <strong>⏰ This verification link will expire in 24 hours.</strong>
                    If you didn't create an account with EduTech, you can safely ignore this email.
                </div>
            </div>
        </div>
        <div class="footer">
            <p><strong>EduTech</strong> — Empowering Education Through Technology</p>
            <div class="copyright">© %d EduTech. All rights reserved.<br>This is an automated message, please do not reply.</div>
        </div>
    </div>
</body>
</html>
`, firstName, verificationURL, verificationURL, year)
}

// BuildPasswordResetEmail builds the HTML email body for password reset
func BuildPasswordResetEmail(firstName, resetURL, token, email string) string {
	now := time.Now()
	year := now.Year()

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            line-height: 1.6;
            color: #1a202c;
            background: #f7fafc;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.08);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #1a202c 0%%, #2d3748 100%%);
            padding: 40px 30px 30px;
            text-align: center;
        }
        .header h1 { color: #ffffff; font-size: 28px; font-weight: 700; margin-bottom: 5px; }
        .header p { color: rgba(255,255,255,0.8); font-size: 16px; }
        .content { padding: 40px 35px 35px; }
        .greeting { font-size: 18px; font-weight: 600; color: #1a202c; margin-bottom: 15px; }
        .greeting span { color: #F13178; }
        .message { color: #4a5568; font-size: 16px; margin: 10px 0 20px; line-height: 1.8; }
        .reset-box {
            background: linear-gradient(135deg, #fefcbf 0%%, #f6e05e 20%%, #fefcbf 100%%);
            border-radius: 12px;
            padding: 30px;
            margin: 25px 0;
            text-align: center;
            border: 2px solid #f6e05e;
        }
        .reset-box .icon { font-size: 48px; margin-bottom: 15px; display: block; }
        .reset-box h3 { color: #744210; font-size: 18px; font-weight: 600; margin-bottom: 10px; }
        .btn {
            display: inline-block;
            padding: 14px 40px;
            background: linear-gradient(135deg, #2d3748 0%%, #1a202c 100%%);
            color: #ffffff;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
            box-shadow: 0 4px 15px rgba(26,32,44,0.3);
            transition: all 0.3s ease;
        }
        .btn:hover { transform: translateY(-2px); box-shadow: 0 8px 25px rgba(26,32,44,0.4); }
        .divider { border: none; height: 1px; background: linear-gradient(90deg, transparent, #e2e8f0, transparent); margin: 25px 0; }
        .token-box {
            background: #f7fafc;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            padding: 12px 16px;
            font-family: 'Courier New', monospace;
            font-size: 13px;
            color: #2d3748;
            word-break: break-all;
            margin-top: 8px;
        }
        .security-notice {
            background: #fed7d7;
            border: 1px solid #feb2b2;
            border-radius: 8px;
            padding: 15px 20px;
            margin: 20px 0;
            display: flex;
            align-items: flex-start;
            gap: 12px;
        }
        .security-notice .text { font-size: 14px; color: #9b2c2c; line-height: 1.6; }
        .tips-box {
            background: #ebf8ff;
            border: 1px solid #bee3f8;
            border-radius: 8px;
            padding: 15px 20px;
            margin: 20px 0;
        }
        .tips-box h4 { color: #2a69ac; font-size: 14px; font-weight: 600; margin-bottom: 8px; }
        .tips-box ul { list-style: none; padding: 0; margin: 0; }
        .tips-box ul li { color: #2b6cb0; font-size: 14px; padding: 4px 0; padding-left: 24px; position: relative; }
        .tips-box ul li::before { content: '✓'; position: absolute; left: 0; color: #2b6cb0; font-weight: bold; }
        .footer {
            background: #f7fafc;
            padding: 30px 35px;
            text-align: center;
            border-top: 1px solid #e2e8f0;
        }
        .footer p { color: #718096; font-size: 13px; line-height: 1.8; }
        .footer .copyright { margin-top: 15px; font-size: 12px; color: #a0aec0; }
        @media (max-width: 480px) {
            .header { padding: 30px 20px 25px; }
            .header h1 { font-size: 22px; }
            .content { padding: 25px 20px; }
            .reset-box { padding: 20px; }
            .btn { display: block; text-align: center; padding: 14px 20px; }
            .footer { padding: 20px 20px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Reset Your Password</h1>
            <p>Secure your account with a new password</p>
        </div>
        <div class="content">
            <div class="greeting">Hello <span>%s</span>!</div>
            <p class="message">
                We received a request to reset the password for your <strong>EduTech</strong> account. 
                Click the button below to create a new password.
            </p>
            <div class="reset-box">
                <span class="icon">🔑</span>
                <h3>Reset Your Password</h3>
                <p>Create a new password to secure your account</p>
                <a href="%s" class="btn">Reset Password →</a>
            </div>
            <hr class="divider">
            <div style="margin: 20px 0 10px;">
                <div style="font-size:13px;color:#718096;font-weight:500;text-transform:uppercase;letter-spacing:0.5px;">Or use this link</div>
                <div class="token-box">%s</div>
            </div>
            <div class="security-notice">
                <div class="text">
                    <strong>⚠️ Security Notice:</strong> 
                    This link will expire in 1 hour for your security. 
                    If you didn't request a password reset, please ignore this email 
                    and ensure your account is secure.
                </div>
            </div>
            <div class="tips-box">
                <h4>💡 Password Tips</h4>
                <ul>
                    <li>Use at least 8 characters</li>
                    <li>Include uppercase and lowercase letters</li>
                    <li>Add numbers and special characters</li>
                    <li>Avoid using common words or personal information</li>
                </ul>
            </div>
        </div>
        <div class="footer">
            <p><strong>EduTech</strong> — Empowering Education Through Technology</p>
            <div class="copyright">© %d EduTech. All rights reserved.<br>This is an automated message, please do not reply.</div>
        </div>
    </div>
</body>
</html>
`, firstName, resetURL, resetURL, year)
}