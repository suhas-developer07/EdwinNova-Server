package mail

const RegistrationTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Registration Successful - EdwinNova Hackathon</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family:Arial, sans-serif;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:20px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff; border-radius:8px; overflow:hidden; box-shadow:0 4px 12px rgba(0,0,0,0.1);">
                    
                    <!-- Header -->
                    <tr>
                        <td style="background:linear-gradient(90deg,#0f2027,#203a43,#2c5364); padding:25px; text-align:center; color:#ffffff;">
                            <h1 style="margin:0; font-size:26px;">🚀 EdwinNova Hackathon</h1>
                            <p style="margin:5px 0 0; font-size:14px;">Registration Successful</p>
                        </td>
                    </tr>

                    <!-- Body -->
                    <tr>
                        <td style="padding:30px;">
                            <h2 style="color:#2c5364; margin-top:0;">Hello {{.TeamName}},</h2>
                            
                            <p style="font-size:15px; color:#333; line-height:1.6;">
                                Congratulations! 🎉 Your team has successfully registered for the 
                                <strong>EdwinNova Hackathon</strong>.
                            </p>

                            <table width="100%" cellpadding="10" cellspacing="0" style="background:#f8f9fa; border-radius:6px; margin:20px 0;">
                                <tr>
                                    <td style="font-size:14px; color:#333;">
                                        <strong>Team Name:</strong> {{.TeamName}}<br>
                                        <strong>Team Leader:</strong> {{.PMName}}<br>
                                        <strong>Email:</strong> {{.PMEmail}}<br>
                                        <strong>Contact:</strong> {{.PMContact}}<br>
                                        <strong>Registration ID:</strong> {{.ApplicationID}}<br>
                                        <strong>Registered On (IST):</strong> {{.IndianTime}}
                                    </td>
                                </tr>
                            </table>

                            <p style="font-size:15px; color:#333; line-height:1.6;">
                                Please keep this email for your records. Further updates regarding event schedule, 
                                venue details, and guidelines will be shared soon.
                            </p>

                            <div style="text-align:center; margin:30px 0;">
                                <a href="edwinslab.com" style="background:#2c5364; color:#ffffff; text-decoration:none; padding:12px 25px; border-radius:5px; font-size:14px;">
                                    View Event Details
                                </a>
                            </div>

                            <p style="font-size:14px; color:#555;">
                                If you have any questions, feel free to contact the organizing team.
                            </p>

                            <p style="margin-top:30px; font-size:14px; color:#333;">
                                Best Regards,<br>
                                <strong>Team EdwinNova</strong>
                            </p>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="background:#f1f1f1; text-align:center; padding:15px; font-size:12px; color:#777;">
                            © 2026 EdwinNova Hackathon. All rights reserved.<br>
                            This is an automated email. Please do not reply.
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
