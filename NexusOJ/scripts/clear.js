import os from 'os'
import { exec } from 'child_process'
// darwin 我的mac

if (os.platform() === 'win32') {
  exec(`rmdir /s /q dist && cls`, (err, stdout, stderr) => {
    if (err) {
      console.error(`Error: ${stderr}`)
      return
    }
    console.log(`Output: ${stdout}`)
  })
} else {
  exec('rm -rf dist && clear', (err, stdout, stderr) => {
    if (err) {
      console.error(`Error: ${stderr}`)
      return
    }
    console.log(`Output: ${stdout}`)
  })
}
