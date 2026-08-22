import { createReadStream } from 'node:fs'
import { stat } from 'node:fs/promises'
import { createServer } from 'node:http'
import { extname, join, normalize } from 'node:path'
import { Readable } from 'node:stream'
import serverEntry from './dist/server/server.js'

const port = Number(process.env.PORT || 3000)
const publicDir = join(process.cwd(), 'dist', 'client')

const contentTypes = {
  '.css': 'text/css; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.webp': 'image/webp',
  '.ico': 'image/x-icon',
}

function safeStaticPath(pathname) {
  const cleanPath = normalize(decodeURIComponent(pathname)).replace(/^([/\\])+/, '')
  const filePath = join(publicDir, cleanPath)
  return filePath.startsWith(publicDir) ? filePath : null
}

async function sendStatic(req, res, pathname) {
  if (!pathname.startsWith('/assets/')) return false

  const filePath = safeStaticPath(pathname)
  if (!filePath) return false

  try {
    const fileStat = await stat(filePath)
    if (!fileStat.isFile()) return false

    res.writeHead(200, {
      'content-length': fileStat.size,
      'content-type': contentTypes[extname(filePath)] || 'application/octet-stream',
      'cache-control': 'public, max-age=31536000, immutable',
    })
    if (req.method === 'HEAD') {
      res.end()
      return true
    }
    createReadStream(filePath).pipe(res)
    return true
  } catch {
    return false
  }
}

createServer(async (req, res) => {
  try {
    const url = new URL(req.url || '/', `http://${req.headers.host || `localhost:${port}`}`)
    if (await sendStatic(req, res, url.pathname)) return

    const headers = new Headers()
    for (const [key, value] of Object.entries(req.headers)) {
      if (Array.isArray(value)) {
        for (const entry of value) headers.append(key, entry)
      } else if (value !== undefined) {
        headers.set(key, value)
      }
    }

    const request = new Request(url, {
      method: req.method,
      headers,
      body: req.method === 'GET' || req.method === 'HEAD' ? undefined : Readable.toWeb(req),
      duplex: req.method === 'GET' || req.method === 'HEAD' ? undefined : 'half',
    })

    const response = await serverEntry.fetch(request)
    res.writeHead(response.status, Object.fromEntries(response.headers.entries()))
    if (!response.body || req.method === 'HEAD') {
      res.end()
      return
    }
    Readable.fromWeb(response.body).pipe(res)
  } catch (error) {
    console.error(error)
    res.writeHead(500, { 'content-type': 'text/plain; charset=utf-8' })
    res.end('Internal Server Error')
  }
}).listen(port, '0.0.0.0', () => {
  console.log(`Mini ERP web listening on ${port}`)
})
