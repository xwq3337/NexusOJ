#!/usr/bin/env node

/**
 * 验证 Turborepo 设置
 */

import { existsSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const root = join(__dirname, '..')

const checks = []

function checkFile(path, description) {
  const fullPath = join(root, path)
  const exists = existsSync(fullPath)
  checks.push({ path, description, exists })
  console.log(`  ${exists ? '✅' : '❌'} ${path} (${description})`)
  return exists
}

console.log('🔍 验证 Turborepo 设置...\n')

// 检查核心配置文件
console.log('📄 核心配置文件:')
checkFile('package.json', '根 package.json')
checkFile('turbo.json', 'Turborepo 配置')
checkFile('pnpm-workspace.yaml', 'pnpm workspace 配置')
checkFile('tsconfig.json', '根 TypeScript 配置')
checkFile('.gitignore', 'Git 忽略文件')
console.log('')

// 检查应用
console.log('📱 应用:')
checkFile('apps/client/package.json', '客户端应用')
checkFile('apps/admin/package.json', '管理端应用')
console.log('')

// 检查共享包
console.log('📦 共享包:')
checkFile('packages/ui/package.json', 'UI 组件包')
checkFile('packages/config/package.json', '配置包')
checkFile('packages/eslint-config/package.json', 'ESLint 配置包')
checkFile('packages/typescript/package.json', 'TypeScript 配置包')
checkFile('packages/utils/package.json', '工具函数包')
checkFile('packages/server/package.json', '后端服务包')
console.log('')

// 汇总结果
const failed = checks.filter(c => !c.exists)
if (failed.length > 0) {
  console.log('❌ 以下文件缺失:')
  failed.forEach(({ path, description }) => {
    console.log(`   - ${path} (${description})`)
  })
  process.exit(1)
} else {
  console.log('✅ 所有检查通过!')
  console.log('\n📝 下一步:')
  console.log('   1. 运行 pnpm install 安装依赖')
  console.log('   2. 运行 pnpm dev 启动开发服务器')
  console.log('   3. 运行 pnpm build 构建所有应用')
}
