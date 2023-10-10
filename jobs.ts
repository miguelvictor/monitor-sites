import { getConfig } from "./config"
import { notifyEmails, notifySlack } from "./notifiers"

/**
 * Checks if the site is up and running.
 * The site is up and running if the response status is 2xx.
 */
export async function runHealthCheckJob() {
  const config = await getConfig()
  const checkPath = async (url: string) => {
    try {
      const response = await fetch(url)
      if (!response.ok) {
        console.error(
          `[runHealthCheckJob] ${url} is down: ${response.status} - ${response.statusText}`
        )
        return notifyEmails(url, `${response.status} - ${response.statusText}`)
      }
    } catch (e) {
      const reason = (e as any)?.message || "Unknown Error"
      console.error(`[runHealthCheckJob] ${url} is unreachable: ${reason}`)
      return notifyEmails(url, `Site Unreachable (${reason})`)
    }
  }

  // perform health check for each site concurrently
  await Promise.allSettled(config.sites.map(({ site }) => checkPath(site)))
}

/**
 * Checks if there are any leaked potentially sensitive files in the site.
 */
export async function runCheckSensitiveFilesJob() {
  const config = await getConfig()
  const paths = [".env", ".env.production", "php.ini"]
  const checkPath = async (url: string) => {
    const response = await fetch(url, { redirect: "manual" })
    if (response.status >= 300) return

    const body = await response.text()
    const title = `Exposed file: ${url}`
    const status = `${response.status} ${response.statusText}`
    console.error(`[runCheckSensitiveFilesJob] exposed file: ${url} - ${status}`)

    return notifySlack(
      title,
      `*Exposed file:* ${url} (${status})\n\`\`\`${body.slice(0, 500)}\`\`\``
    )
  }

  // perform health check for each site concurrently
  await Promise.allSettled(
    config.sites
      .flatMap(({ site }) => paths.map(path => new URL(path, site).toString()))
      .map(url => checkPath(url).catch(console.error))
  )
}
