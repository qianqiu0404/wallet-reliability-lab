import { execFileSync } from 'node:child_process'
import { readFileSync } from 'node:fs'

const files = execFileSync('git', ['ls-files', '--cached', '--others', '--exclude-standard'], { encoding: 'utf8' })
  .trim()
  .split('\n')
  .filter(Boolean)
  .filter(file => !file.endsWith('package-lock.json') && !file.endsWith('lab-overview.png'))

const forbidden = [
  { label: 'absolute local path', pattern: /\/Users\/[^/\s]+\// },
  { label: 'private key block', pattern: /-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/ },
  { label: 'credential assignment', pattern: /(?:api[_-]?key|private[_-]?key|secret|token)\s*[=:]\s*(?!\[?REDACTED\]?|$)["']?[A-Za-z0-9_\-]{12,}/i },
  { label: 'non-loopback service target', pattern: /LAB_SERVICE_TARGETS=.*(?:10\.|192\.168\.|172\.(?:1[6-9]|2\d|3[01])\.)/ },
]

const failures = []
for (const file of files) {
  const value = readFileSync(file, 'utf8')
  for (const rule of forbidden) {
    if (rule.pattern.test(value)) failures.push(`${file}: ${rule.label}`)
  }
}

if (failures.length) {
  console.error(`Public safety check failed:\n${failures.join('\n')}`)
  process.exit(1)
}

console.log(`Public safety check passed for ${files.length} files.`)
