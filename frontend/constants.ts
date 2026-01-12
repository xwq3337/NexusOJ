/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-12-16 13:14:02
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-12 20:01:28
 * @FilePath: /frontend/constants.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Problem, Contest } from './src/types'

export const MOCK_PROBLEMS: Problem[] = [
  {
    id: '1',
    title: '两数之和',
    difficulty: 'Easy',
    acceptance: '48.5%',
    tags: ['数组', '哈希表'],
    status: 'todo'
  },
  {
    id: '2',
    title: '两数相加',
    difficulty: 'Medium',
    acceptance: '39.1%',
    tags: ['链表', '数学'],
    status: 'todo'
  },
  {
    id: '3',
    title: '无重复字符的最长子串',
    difficulty: 'Medium',
    acceptance: '32.1%',
    tags: ['字符串', '滑动窗口']
  },
  {
    id: '4',
    title: '寻找两个正序数组的中位数',
    difficulty: 'Hard',
    acceptance: '35.2%',
    tags: ['二分查找', '数组']
  },
  {
    id: '5',
    title: '最长回文子串',
    difficulty: 'Medium',
    acceptance: '32.1%',
    tags: ['字符串', '动态规划']
  },
  { id: '6', title: 'Z 字形变换', difficulty: 'Medium', acceptance: '42.4%', tags: ['字符串'] }
]

export const MOCK_CONTESTS: Contest[] = [
  {
    id: '358',
    title: '周赛 358',
    startTime: '2025/12/15 14:44',
    duration: '1h 30m',
    registered: 12450,
    type: 'Weekly',
    status: 'Live'
  },
  {
    id: '110',
    title: '双周赛 110',
    startTime: '2025/12/16 14:44',
    duration: '1h 30m',
    registered: 4300,
    type: 'Biweekly',
    status: 'Upcoming'
  },
  {
    id: '4',
    title: 'Nexus新手杯 #4',
    startTime: '2025/12/17 14:44',
    duration: '2h 00m',
    registered: 1200,
    type: 'Cup',
    status: 'Upcoming'
  }
]

export const ACTIVITY_DATA = [
  { name: 'Mon', submissions: 400 },
  { name: 'Tue', submissions: 300 },
  { name: 'Wed', submissions: 550 },
  { name: 'Thu', submissions: 500 },
  { name: 'Fri', submissions: 700 },
  { name: 'Sat', submissions: 800 },
  { name: 'Sun', submissions: 750 }
]

// 语言配置映射表
export const LANGUAGE_CONFIG = {
  c: {
    value: 'c',
    label: 'C',
    aceMode: 'c_cpp',
    apiValue: 'c',
    defaultCode: `#include <stdio.h>
int main() {
    // Your code here
    return 0;
}`
  },
  cpp: {
    value: 'cpp',
    label: 'C++',
    aceMode: 'c_cpp',
    apiValue: 'cpp',
    defaultCode: `#include <iostream>
using namespace std;

int main() {
    // Your code here
    return 0;
}`
  },
  python: {
    value: 'python',
    label: 'Python',
    aceMode: 'python',
    apiValue: 'python',
    defaultCode: `def main():
    # Your code here
    pass

if __name__ == "__main__":
    main()`
  },
  java: {
    value: 'java',
    label: 'Java',
    aceMode: 'java',
    apiValue: 'java',
    defaultCode: `public class Main {
    public static void main(String[] args) {
        // Your code here
    }
}`
  },
  javascript: {
    value: 'javascript',
    label: 'JavaScript',
    aceMode: 'javascript',
    apiValue: 'javascript',
    defaultCode: `// Your code here
console.log("Hello, World!");`
  },
  go: {
    value: 'go',
    label: 'Go',
    aceMode: 'go',
    apiValue: 'javascript',
    defaultCode: `package main
    
import "fmt"

func main(){
    fmt.Println("Hello world")
}`
  },
  rust: {
    value: 'rust',
    label: 'Rust',
    aceMode: 'rust',
    apiValue: 'rust',
    defaultCode: `fn main(){
  println!("Hello world")
}`
  },
  csharp: {
    value: 'csharp',
    label: 'C#',
    aceMode: 'csharp',
    apiValue: 'csharp',
    defaultCode: `using System;

class Program {
    static void Main() {
        // Your code here
    }
}`
  }
} as const

function transformCodeObject(inputObj) {
    const transformed = {};
    // 遍历输入对象的每个键（如 'c'）
    for (const key in inputObj) {
        // 检查属性是否为对象自身所有（非继承）
        if (inputObj.hasOwnProperty(key)) {
            const value = inputObj[key];
            
            // 判断该属性的值是否为对象，并且包含 defaultCode 属性
            if (value && typeof value === 'object' && 'defaultCode' in value) {
                // 如果满足条件，则提取 defaultCode 的值
                transformed[key] = value.defaultCode;
            } else {
                // 如果不满足条件，则保留原始值
                transformed[key] = value;
            }
        }
    }
    
    return transformed;
}

export const DEFAULT_CODES = transformCodeObject(LANGUAGE_CONFIG) as Record<string,string>
export const LANGUAGE_OPTIONS = Object.keys(LANGUAGE_CONFIG)
export const ACE_MODE_OPTIONS = LANGUAGE_OPTIONS.map(
  (lang) => LANGUAGE_CONFIG[lang as keyof typeof LANGUAGE_CONFIG].aceMode
)
export type ACE_MODE = (typeof ACE_MODE_OPTIONS)[number]
export type LanguageValue = keyof typeof LANGUAGE_CONFIG

export const EDITOR_THEME_OPTIONS = [
  'github',
  'chrome',
  'monokai',
  'xcode',
  'dracula',
  'clouds',
  'terminal'
]
export type EDITOR_THEHE = (typeof EDITOR_THEME_OPTIONS)[number]
