import nodemailer from "nodemailer";
import dayjs from "dayjs";
import { getConfig } from "./config";

const transporter = nodemailer.createTransport({
  host: "smtp.gmail.com",
  port: 465,
  secure: true,
  auth: {
    user: Bun.env.SMTP_USER,
    pass: Bun.env.SMTP_PASS,
  },
});

export async function notifyEmails(site: string, status: string) {
  const config = await getConfig();
  await transporter.sendMail({
    from: `DF DevOps <${Bun.env.SMTP_USER}>`,
    to: config.emails,
    subject: `${new URL(site).host} is down`,
    html: `
          <div>Site: ${site}</div>
          <div>Status: ${status}</div>
          <div>Timestamp: ${dayjs().format("MMMM D, HH:mm")}</div>
      `,
  });
}

export async function notifySlack() {}
