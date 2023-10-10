import { getConfig } from "./config";
import { notifyEmails } from "./notifiers";

export async function runHealthCheckJob() {
  const config = await getConfig();
  await Promise.all(
    config.sites.map(({ site }) =>
      fetch(site)
        .then((res) => {
          if (!res.ok) notifyEmails(site, `${res.status} - ${res.statusText}`);
        })
        .catch((e) => {
          notifyEmails(
            site,
            `Site Unreachable (${e.message || "Unknown Error"})`
          );
        })
    )
  );
}

export async function runCheckSensitiveFilesJob() {}
