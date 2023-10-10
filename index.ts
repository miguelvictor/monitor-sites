import cron from "node-cron"
import { runCheckSensitiveFilesJob, runHealthCheckJob } from "./jobs"

cron.schedule("*/3 * * * *", runHealthCheckJob)
cron.schedule("* * * * *", runCheckSensitiveFilesJob)
