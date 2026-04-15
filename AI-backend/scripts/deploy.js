// 把dist文件夹和.env、system_prompt.txt文件拷贝到服务器上

const ADDRESS = "root@47.109.57.7";
const DIRECTORY = "/opt/oj-agent";

import { execSync } from "child_process";

const files = ["dist", ".env", "system_prompt.txt", "package.json"];

for (const file of files) {
  console.log(`上传 ${file} ...`);
  execSync(`scp -r ${file} ${ADDRESS}:${DIRECTORY}`, { stdio: "inherit" });
}

console.log("安装依赖...");
execSync(`ssh ${ADDRESS} "cd ${DIRECTORY} && npm install --omit=dev --legacy-peer-deps"`, {
  stdio: "inherit",
});

console.log("重启服务...");
execSync(`ssh ${ADDRESS} "cd ${DIRECTORY} && pkill -f 'node dist/index.js' || true && nohup node dist/index.js > /dev/null 2>&1 &"`, {
  stdio: "inherit",
});

console.log("部署完成!");
