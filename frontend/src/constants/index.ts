/**
 * Markdown 代码主题
 */
export type MarkdownCodeTheme =
  | 'atom'
  | 'dark'
  | 'light'
  | 'github'
  | 'monokai'
  | 'vs'
  | 'tomorrow'
  | 'twilight'
  | 'xcode'

/**
 * 支持的 Markdown 代码主题列表
 */
export const MARKDOWN_CODE_THEMES: { value: MarkdownCodeTheme; label: string }[] = [
  { value: 'github', label: 'GitHub' },
  { value: 'atom', label: 'Atom' },
  { value: 'dark', label: 'Dark' },
  { value: 'light', label: 'Light' },
  { value: 'monokai', label: 'Monokai' },
  { value: 'vs', label: 'VS Code' },
  { value: 'tomorrow', label: 'Tomorrow' },
  { value: 'twilight', label: 'Twilight' },
  { value: 'xcode', label: 'Xcode' }
]
import { Problem } from '@/types/problem'
import { Contest } from '@/types/contest'

export const MOCK_PROBLEMS: Problem[] = [
  {
    id: 1,
    title: '两数之和',
    difficulty: 0,
    tags: ['数组', '哈希表'],
    user_id: '3953240110606292',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  },
  {
    id: 2,
    title: '两数相加',
    difficulty: 0,
    tags: ['链表', '数学'],
    user_id: '1',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  },
  {
    id: 3,
    title: '无重复字符的最长子串',
    difficulty: 0,
    tags: ['字符串', '滑动窗口'],
    user_id: '1',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  },
  {
    id: 4,
    title: '寻找两个正序数组的中位数',
    difficulty: 0,
    tags: ['二分查找', '数组'],
    user_id: '1',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  },
  {
    id: 5,
    title: '最长回文子串',
    difficulty: 0,
    tags: ['字符串', '动态规划'],
    user_id: '1',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  },
  {
    id: 6,
    title: 'Z 字形变换',
    difficulty: 0,
    tags: ['字符串'],
    user_id: '1',
    accept: 21431,
    submission: 63216,
    context: 'null',
    input_description: 'input',
    output_description: 'ouput',
    judge_case: [
      {
        input: '1 2',
        expected: '3'
      }
    ],
    judge_sample: [
      {
        input: '3 4',
        expected: '7'
      }
    ],
    tips: 'null'
  }
]

export const MOCK_CONTESTS: Contest[] = [
  {
    id: '1',
    title: '周赛 1',
    startTime: '2025/12/15 14:44',
    duration: '1h 30m',
    registered: 12450,
    type: 'Weekly',
    status: 'Live'
  },
  {
    id: '2',
    title: '双周赛 1',
    startTime: '2025/12/16 14:44',
    duration: '1h 30m',
    registered: 4300,
    type: 'Biweekly',
    status: 'Upcoming'
  },
  {
    id: '3',
    title: 'Nexus新手杯 #1',
    startTime: '2025/12/17 14:44',
    duration: '2h 00m',
    registered: 1200,
    type: 'Cup',
    status: 'Upcoming'
  }
]

export interface ContestRanking {
  rank: number
  username: string
  avatar: string
  score: number
  solved: number
  time: string
  penalty: number
}

export interface ContestProblem {
  id: string
  title: string
  difficulty: 'Easy' | 'Medium' | 'Hard'
  acceptRate: string
  solved: number
  status?: 'accepted' | 'wrong' | 'pending'
}

export const difficultyMap = [
  { text: '简单', color: 'text-green-400', type: 'success' },
  { text: '容易', color: 'text-yellow-400', type: 'warning' },
  { text: '中等', color: 'text-orange-400', type: 'info' },
  { text: '困难', color: 'text-red-400', type: 'error' },
  { text: '极难', color: 'text-purple-400', type: 'error' }
]
export const MOCK_CONTEST_RANKING: ContestRanking[] = [
  {
    rank: 1,
    username: 'algorithm_master',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=algorithm',
    score: 3200,
    solved: 4,
    time: '1:23:45',
    penalty: 12
  },
  {
    rank: 2,
    username: 'code_ninja',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=ninja',
    score: 3050,
    solved: 4,
    time: '1:28:32',
    penalty: 18
  },
  {
    rank: 3,
    username: 'byte_warrior',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=warrior',
    score: 2900,
    solved: 3,
    time: '1:35:12',
    penalty: 8
  },
  {
    rank: 4,
    username: 'pixel_coder',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=pixel',
    score: 2750,
    solved: 3,
    time: '1:42:05',
    penalty: 22
  },
  {
    rank: 5,
    username: 'data_drifter',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=drifter',
    score: 2600,
    solved: 3,
    time: '1:48:33',
    penalty: 15
  },
  {
    rank: 6,
    username: 'logic_lord',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=logic',
    score: 2450,
    solved: 2,
    time: '1:55:18',
    penalty: 10
  },
  {
    rank: 7,
    username: 'syntax_sage',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=syntax',
    score: 2300,
    solved: 2,
    time: '2:02:45',
    penalty: 25
  },
  {
    rank: 8,
    username: 'quantum_dev',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=quantum',
    score: 2150,
    solved: 2,
    time: '2:10:12',
    penalty: 14
  }
]

export const MOCK_CONTEST_PROBLEMS: ContestProblem[] = [
  {
    id: 'A',
    title: '数字游戏',
    difficulty: 'Easy',
    acceptRate: '78.5%',
    solved: 12450,
    status: 'accepted'
  },
  {
    id: 'B',
    title: '字符串变换',
    difficulty: 'Easy',
    acceptRate: '65.2%',
    solved: 9850,
    status: 'accepted'
  },
  {
    id: 'C',
    title: '矩阵路径',
    difficulty: 'Medium',
    acceptRate: '42.8%',
    solved: 5320,
    status: 'pending'
  },
  {
    id: 'D',
    title: '树形 DP',
    difficulty: 'Medium',
    acceptRate: '28.3%',
    solved: 3520,
    status: 'wrong'
  },
  {
    id: 'E',
    title: '最短路优化',
    difficulty: 'Hard',
    acceptRate: '12.5%',
    solved: 1560,
    status: 'pending'
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
    label: 'Gcc 13.3.0',
    aceMode: 'c_cpp',
    apiValue: 'c',
    color: {
      color: '#2080f01a',
      textColor: 'rgb(32, 128, 240)',
      borderColor: 'rgba(32, 128, 240, 0.3)'
    },
    defaultCode: `#include <stdio.h>
int main() {
    // Your code here
    return 0;
}`
  },
  cpp: {
    value: 'cpp',
    label: 'G++ 13.3.0',
    aceMode: 'c_cpp',
    apiValue: 'cpp',
    color: {
      color: '#2080f01a',
      textColor: 'rgb(32, 128, 240)',
      borderColor: 'rgba(32, 128, 240, 0.3)'
    },
    defaultCode: `#include <iostream>
using namespace std;

int main() {
    // Your code here
    return 0;
}`
  },
  python: {
    value: 'python',
    label: 'Python 3.12.3',
    aceMode: 'python',
    apiValue: 'python',
    color: {
      color: '#18a0581a',
      textColor: '#18a058',
      borderColor: 'rgba(24, 160, 88, 0.3)'
    },
    defaultCode: `def main():
    # Your code here
    pass

if __name__ == "__main__":
    main()`
  },
  java: {
    value: 'java',
    label: 'Java 17',
    aceMode: 'java',
    apiValue: 'java',
    color: {
      color: '#f0a0201a',
      textColor: 'rgb(240, 160, 32)',
      borderColor: 'rgba(240, 160, 32, 0.3)'
    },
    defaultCode: `public class Main {
    public static void main(String[] args) {
        // Your code here
    }
}`
  },
  javascript: {
    value: 'javascript',
    label: 'Node 25.3.0',
    aceMode: 'javascript',
    apiValue: 'javascript',
    color: {
      color: '#f0a0201a',
      textColor: 'rgb(240, 160, 32)',
      borderColor: 'rgba(240, 160, 32, 0.3)'
    },
    defaultCode: `// Your code here
console.log("Hello, World!");`
  },
  go: {
    value: 'go',
    label: 'Go 1.25.5',
    aceMode: 'go',
    apiValue: 'go',
    color: {
      color: '#18a0581a',
      textColor: '#18a058',
      borderColor: 'rgba(24, 160, 88, 0.3)'
    },
    defaultCode: `package main
    
import "fmt"

func main(){
    fmt.Println("Hello world")
}`
  },
  rust: {
    value: 'rust',
    label: 'Rust 1.92.0',
    aceMode: 'rust',
    apiValue: 'rust',
    color: {
      color: '#f0a0201a',
      textColor: '#f0a020',
      borderColor: 'rgba(240, 160, 32, 0.3)'
    },
    defaultCode: `fn main(){
  println!("Hello world")
}`
  },
  csharp: {
    value: 'csharp',
    label: 'C# 6.8.0.105',
    aceMode: 'csharp',
    apiValue: 'csharp',
    color: {
      color: '#2080f01a',
      textColor: 'rgb(32, 128, 240)',
      borderColor: 'rgba(32, 128, 240, 0.3)'
    },
    defaultCode: `using System;

class Program {
    static void Main() {
        // Your code here
    }
}`
  },
  zig: {
    value: 'zig',
    label: 'Zig 0.16.0',
    aceMode: 'zig',
    apiValue: 'zig',
    color: {
      color: '#f0a0201a',
      textColor: 'rgb(240, 160, 32)',
      borderColor: 'rgba(240, 160, 32, 0.3)'
    },
    defaultCode: ''
  }
} as const

function transformCodeObject(inputObj) {
  const transformed = {}
  // 遍历输入对象的每个键（如 'c'）
  for (const key in inputObj) {
    // 检查属性是否为对象自身所有（非继承）
    if (inputObj.hasOwnProperty(key)) {
      const value = inputObj[key]

      // 判断该属性的值是否为对象，并且包含 defaultCode 属性
      if (value && typeof value === 'object' && 'defaultCode' in value) {
        // 如果满足条件，则提取 defaultCode 的值
        transformed[key] = value.defaultCode
      } else {
        // 如果不满足条件，则保留原始值
        transformed[key] = value
      }
    }
  }

  return transformed
}

export const DEFAULT_CODES = transformCodeObject(LANGUAGE_CONFIG) as Record<string, string>
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

export const STATUS_OPTIONS = [
  { label: '通过', value: 'Accepted' },
  { label: '答案错误', value: 'WrongAnswer' },
  { label: '超时', value: 'TimeLimitExceeded' },
  { label: '内存超限', value: 'MemoryLimitExceeded' },
  { label: '运行错误', value: 'RuntimeError' },
  { label: '编译错误', value: 'CompilationError' }
]

export type VerdictType = "Accepted" | "WrongAnswer" | "TimeLimitExceeded" | "MemoryLimitExceeded" | "RuntimeError" | "CompilationError"

export const STATUS_COLORS = {
  // success
  Accepted: {
    color: '#18a0581a',
    textColor: '#18a058',
    borderColor: 'rgba(24, 160, 88, 0.3)'
  },
  //error
  WrongAnswer: {
    color: '#d030581a',
    textColor: '#d03050',
    borderColor: 'rgba(208, 48, 80, 0.3)'
  },
  // error
  TimeLimitExceeded: {
    color: '#d030581a',
    textColor: '#d03050',
    borderColor: 'rgba(208, 48, 80, 0.3)'
  },
  // warning
  MemoryLimitExceeded: {
    color: '#f0a0201a',
    textColor: '#f0a020',
    borderColor: 'rgba(240, 160, 32, 0.3)'
  },
  // error
  RuntimeError: {
    color: '#d030581a',
    textColor: '#d03050',
    borderColor: 'rgba(208, 48, 80, 0.3)'
  },
  // primary
  CompilationError: {
    color: '#2080f01a',
    textColor: 'rgb(32, 128, 240)',
    borderColor: 'rgba(32, 128, 240, 0.3)'
  }
}
