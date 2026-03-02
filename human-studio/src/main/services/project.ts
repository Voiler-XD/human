import { readdir, readFile, writeFile, mkdir, rm, rename as fsRename, stat } from 'fs/promises'
import { join, relative } from 'path'
import { watch, FSWatcher } from 'fs'
import { spawn } from 'child_process'
import { app } from 'electron'
import { existsSync } from 'fs'

export interface FileEntry {
  name: string
  path: string
  isDirectory: boolean
  children?: FileEntry[]
}

export class ProjectService {
  private watcher: FSWatcher | null = null

  async openProject(dirPath: string): Promise<FileEntry[]> {
    return this.listFiles(dirPath)
  }

  async createProject(name: string, parentDir: string): Promise<string> {
    const projectDir = join(parentDir, name)
    await mkdir(projectDir, { recursive: true })

    // Run `human init` in the parent directory
    const binary = this.getBinaryPath()
    return new Promise((resolve, reject) => {
      const proc = spawn(binary, ['init', name], {
        cwd: parentDir,
        stdio: 'pipe',
      })

      proc.on('error', (err) => reject(err))
      proc.on('close', (code) => {
        if (code === 0) {
          resolve(projectDir)
        } else {
          reject(new Error(`human init exited with code ${code}`))
        }
      })

      // Send default answers to interactive prompts (press Enter for each)
      proc.stdin?.write('\n\n\n\n\n')
      proc.stdin?.end()
    })
  }

  async readFile(filePath: string): Promise<string> {
    return readFile(filePath, 'utf-8')
  }

  async writeFile(filePath: string, content: string): Promise<void> {
    await writeFile(filePath, content, 'utf-8')
  }

  async listFiles(dirPath: string): Promise<FileEntry[]> {
    const entries = await readdir(dirPath, { withFileTypes: true })
    const result: FileEntry[] = []

    for (const entry of entries) {
      // Skip hidden files and common non-essential dirs
      if (entry.name.startsWith('.') && entry.name !== '.human') continue
      if (entry.name === 'node_modules') continue

      const fullPath = join(dirPath, entry.name)
      const fileEntry: FileEntry = {
        name: entry.name,
        path: fullPath,
        isDirectory: entry.isDirectory(),
      }

      if (entry.isDirectory()) {
        fileEntry.children = await this.listFiles(fullPath)
      }

      result.push(fileEntry)
    }

    // Sort: directories first, then alphabetically
    result.sort((a, b) => {
      if (a.isDirectory !== b.isDirectory) return a.isDirectory ? -1 : 1
      return a.name.localeCompare(b.name)
    })

    return result
  }

  async createFile(filePath: string, content = ''): Promise<void> {
    await writeFile(filePath, content, 'utf-8')
  }

  async createDir(dirPath: string): Promise<void> {
    await mkdir(dirPath, { recursive: true })
  }

  async deletePath(targetPath: string): Promise<void> {
    await rm(targetPath, { recursive: true, force: true })
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    await fsRename(oldPath, newPath)
  }

  watch(
    dirPath: string,
    onChange: (event: string, filePath: string) => void
  ): void {
    this.unwatch()
    this.watcher = watch(dirPath, { recursive: true }, (event, filename) => {
      if (filename) {
        onChange(event, join(dirPath, filename))
      }
    })
  }

  unwatch(): void {
    if (this.watcher) {
      this.watcher.close()
      this.watcher = null
    }
  }

  private getBinaryPath(): string {
    if (!app.isPackaged) {
      const devBinary = join(__dirname, '../../../../../human')
      if (existsSync(devBinary)) return devBinary
      return 'human'
    }
    const platform = process.platform === 'win32' ? 'human.exe' : 'human'
    const bundledPath = join(process.resourcesPath, 'bin', platform)
    if (existsSync(bundledPath)) return bundledPath
    return 'human'
  }
}
