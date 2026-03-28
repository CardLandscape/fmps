'use strict'

const { app, BrowserWindow, ipcMain, dialog } = require('electron')
const path = require('path')
const { spawn } = require('child_process')
const fs = require('fs')
const http = require('http')

let goServerProcess = null
let mainWindow = null

function startGoServer() {
  const isDev = process.env.NODE_ENV === 'development'
  let serverProcess

  const env = {
    ...process.env,
    FMPS_PORT: '8080',
    FMPS_DB_PATH: path.join(app.getPath('userData'), 'fmps.db')
  }

  if (isDev) {
    const serverDir = path.join(__dirname, '../server')
    console.log('[server] Starting in dev mode via go run:', serverDir)
    serverProcess = spawn('go', ['run', '.'], {
      cwd: serverDir,
      env,
      stdio: ['ignore', 'pipe', 'pipe']
    })
  } else {
    const binaryName = process.platform === 'win32' ? 'fmps-server.exe' : 'fmps-server'
    const binaryPath = path.join(process.resourcesPath, 'server', binaryName)
    console.log('[server] Starting production binary:', binaryPath)
    serverProcess = spawn(binaryPath, [], {
      env,
      stdio: ['ignore', 'pipe', 'pipe']
    })
  }

  serverProcess.stdout.on('data', (data) => {
    console.log('[server stdout]', data.toString().trim())
  })

  serverProcess.stderr.on('data', (data) => {
    console.error('[server stderr]', data.toString().trim())
  })

  serverProcess.on('close', (code) => {
    console.log(`[server] process exited with code ${code}`)
  })

  serverProcess.on('error', (err) => {
    console.error('[server] failed to start:', err.message)
  })

  return serverProcess
}

function waitForServer(port, maxAttempts = 30) {
  return new Promise((resolve, reject) => {
    let attempts = 0

    function poll() {
      attempts++
      // Use dedicated health endpoint
      const req = http.get(`http://localhost:${port}/api/health`, (res) => {
        console.log(`[server] ready after ${attempts} attempt(s)`)
        resolve()
      })

      req.on('error', () => {
        if (attempts >= maxAttempts) {
          reject(new Error(`Server on port ${port} did not start after ${maxAttempts} attempts`))
          return
        }
        setTimeout(poll, 1000)
      })

      req.setTimeout(1000, () => {
        req.destroy()
        if (attempts >= maxAttempts) {
          reject(new Error(`Server on port ${port} timed out after ${maxAttempts} attempts`))
          return
        }
        setTimeout(poll, 1000)
      })
    }

    poll()
  })
}

async function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false
    }
  })

  try {
    await mainWindow.loadURL('http://localhost:8080')
  } catch (err) {
    console.warn('[window] Failed to load from server, falling back to dist:', err.message)
    const fallback = path.join(__dirname, '../web/dist/index.html')
    if (fs.existsSync(fallback)) {
      await mainWindow.loadFile(fallback)
    } else {
      dialog.showErrorBox('FMPS', 'Unable to load application. Server may not have started correctly.')
    }
  }

  mainWindow.on('closed', () => {
    if (goServerProcess) {
      goServerProcess.kill()
      goServerProcess = null
    }
    mainWindow = null
  })
}

app.whenReady().then(async () => {
  goServerProcess = startGoServer()

  try {
    await waitForServer(8080, 30)
  } catch (err) {
    console.error('[app] Server did not become ready:', err.message)
  }

  await createWindow()

  app.on('activate', async () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      await createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('before-quit', () => {
  if (goServerProcess) {
    goServerProcess.kill()
    goServerProcess = null
  }
})
