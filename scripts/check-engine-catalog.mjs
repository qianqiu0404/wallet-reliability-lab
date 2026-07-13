import { createHash } from 'node:crypto'
import { readFileSync } from 'node:fs'

const lock = JSON.parse(readFileSync(new URL('../catalog/engine-catalog.lock.json', import.meta.url), 'utf8'))
const response = await fetch(lock.sourceUrl)
if (!response.ok) throw new Error(`Unable to load pinned engine catalog: HTTP ${response.status}`)

const bytes = Buffer.from(await response.arrayBuffer())
const checksum = createHash('sha256').update(bytes).digest('hex')
if (checksum !== lock.catalogSha256) {
  throw new Error(`Engine catalog checksum mismatch: received ${checksum}`)
}

const catalog = JSON.parse(bytes.toString('utf8'))
if (catalog.version !== lock.catalogVersion) {
  throw new Error(`Engine catalog version ${catalog.version} does not match lock ${lock.catalogVersion}`)
}

const engineIds = new Set(catalog.scenarios.map(scenario => scenario.id))
const labSource = readFileSync(new URL('../src/data/scenarios.ts', import.meta.url), 'utf8')
for (const mapping of lock.scenarioMappings) {
  if (!engineIds.has(mapping.engineScenarioId)) throw new Error(`Missing engine scenario ${mapping.engineScenarioId}`)
  if (!labSource.includes(`id: '${mapping.labScenarioId}'`)) throw new Error(`Missing lab scenario ${mapping.labScenarioId}`)
}

console.log(`Scenario Catalog v${catalog.version} verified at ${lock.sourceRelease} (${checksum.slice(0, 12)}…).`)
