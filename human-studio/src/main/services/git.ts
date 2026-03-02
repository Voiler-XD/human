import simpleGit, { SimpleGit, StatusResult } from 'simple-git'

export class GitService {
  private getGit(projectDir: string): SimpleGit {
    return simpleGit(projectDir)
  }

  async isRepo(projectDir: string): Promise<boolean> {
    try {
      const git = this.getGit(projectDir)
      return await git.checkIsRepo()
    } catch {
      return false
    }
  }

  async currentBranch(projectDir: string): Promise<string | null> {
    try {
      const git = this.getGit(projectDir)
      const branch = await git.branchLocal()
      return branch.current || null
    } catch {
      return null
    }
  }

  async status(projectDir: string): Promise<StatusResult | null> {
    try {
      const git = this.getGit(projectDir)
      return await git.status()
    } catch {
      return null
    }
  }

  async push(projectDir: string): Promise<void> {
    const git = this.getGit(projectDir)
    await git.push()
  }

  async pull(projectDir: string): Promise<void> {
    const git = this.getGit(projectDir)
    await git.pull()
  }

  async createBranch(projectDir: string, name: string): Promise<void> {
    const git = this.getGit(projectDir)
    await git.checkoutLocalBranch(name)
  }
}
