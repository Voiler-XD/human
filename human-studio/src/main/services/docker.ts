import { execFile } from 'child_process'

export class DockerService {
  async isAvailable(): Promise<boolean> {
    return new Promise((resolve) => {
      execFile('docker', ['version', '--format', '{{.Server.Version}}'], (err) => {
        resolve(!err)
      })
    })
  }
}
