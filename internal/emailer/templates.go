package emailer

var emailTemplate = `
<!doctype html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Новое обращение в поддержку</title>
</head>
<body style="margin:0;padding:0;background-color:#f5f7fb;font-family:Arial,sans-serif;color:#1f2937;">
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color:#f5f7fb;padding:24px 0;">
    <tr>
      <td align="center">
        <table role="presentation" width="600" cellspacing="0" cellpadding="0" style="background:#ffffff;border-radius:12px;padding:32px;box-shadow:0 4px 20px rgba(0,0,0,0.08);">
          <tr>
            <td>
              <h2 style="margin:0 0 24px;font-size:24px;color:#111827;">Новое обращение в поддержку</h2>

              <p style="margin:0 0 12px;font-size:14px;color:#6b7280;">Отправитель</p>
              <p style="margin:0 0 20px;font-size:16px;color:#111827;"><strong>%s</strong></p>

              <p style="margin:0 0 12px;font-size:14px;color:#6b7280;">Тема</p>
              <p style="margin:0 0 20px;font-size:16px;color:#111827;"><strong>%s</strong></p>

              <p style="margin:0 0 12px;font-size:14px;color:#6b7280;">Описание</p>
              <div style="margin:0 0 24px;font-size:15px;line-height:1.6;color:#111827;background:#f9fafb;border-radius:8px;padding:16px;white-space:pre-wrap;">%s</div>

              <hr style="border:none;border-top:1px solid #e5e7eb;margin:24px 0;">

              <p style="margin:0;font-size:12px;color:#9ca3af;">
                Это письмо было сформировано автоматически сервисом Dormitory Life.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`

var supportHeader string = "Новое обращение в поддержку"