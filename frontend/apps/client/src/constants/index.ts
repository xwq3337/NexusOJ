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




export const difficultyMap = [
  { text: '简单', color: 'text-green-400', type: 'success' },
  { text: '容易', color: 'text-yellow-400', type: 'warning' },
  { text: '中等', color: 'text-orange-400', type: 'info' },
  { text: '困难', color: 'text-red-400', type: 'error' },
  { text: '极难', color: 'text-purple-400', type: 'error' }
]




// 语言配置映射表
export const LANGUAGE_CONFIG = {
  c: {
    value: 'c', // 
    label: 'Gcc 13.3.0', // 前端显示的标签
    aceMode: 'c_cpp', // ace 编辑器使用的模式，用于主题
    apiValue: 'c', // sandbox/formatter 使用的值
    color: {
      color: '#2080f01a',
      textColor: 'rgb(32, 128, 240)',
      borderColor: 'rgba(32, 128, 240, 0.3)'
    },
    defaultCode: `#include <stdio.h>
int main() {
    printf("Hello, World!\\n");
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
    cout << "Hello, World!" << endl;
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
        System.out.println("Hello, World!");
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
    defaultCode: `
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
    defaultCode: `
fn main(){
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
        Console.WriteLine("Hello, World!");
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
    defaultCode: `const std = @import("std");
pub fn main() {
    std.debug.print("Hello, World!\\n", .{});
}`
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
