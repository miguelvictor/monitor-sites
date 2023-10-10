import nodemailer from "nodemailer"
import dayjs from "dayjs"
import { getConfig } from "./config"

const transporter = nodemailer.createTransport({
  host: "smtp.gmail.com",
  port: 465,
  secure: true,
  auth: {
    user: Bun.env.SMTP_USER,
    pass: Bun.env.SMTP_PASS,
  },
})

export async function notifyEmails(site: string, status: string) {
  const config = await getConfig()
  await transporter.sendMail({
    from: `DF DevOps <${Bun.env.SMTP_USER}>`,
    to: config.emails,
    subject: `${new URL(site).host} is down`,
    html: `
          <div>Site: ${site}</div>
          <div>Status: ${status}</div>
          <div>Timestamp: ${dayjs().format("MMMM D, HH:mm")}</div>
      `,
  })
}

export async function notifySlack(title: string, contents: string) {
  try {
    const res = await fetch(Bun.env.WEBHOOK_URL!, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        text: title,
        blocks: [
          {
            type: "section",
            text: { type: "mrkdwn", text: contents },
          },
        ],
      }),
    })

    if (!res.ok) {
      console.error(`[notifySlack] failed to notify slack: ${res.status} ${res.statusText}`)
      console.error(`[notifySlack] ${await res.text()}`)
    }
  } catch (e) {
    console.error(`[notifySlack] failed to notify slack: ${e}`)
  }
}
