import os from 'os'
import { exec } from 'child_process'
const cmd = 'rsync -avz  dist/ root@47.109.57.7:/www/wwwroot/oj'

if (os.platform() === 'win32') {
  exec(`sshpass -p '123456' ${cmd}`, (err, stdout, stderr) => {
    if (err) {
      console.error(`Error: ${stderr}`)
      return
    }
    console.log(`Output: ${stdout}`)
  })
} else {
  exec(cmd, { maxBuffer: 1024 * 1024 * 1024 }, (err, stdout, stderr) => {
    if (err) {
      console.error(`Error: ${stderr}`)
      return
    }
    console.log(`Output: ${stdout}`)
  })
}
