/**
 * 编程语言配置
 */
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

/**
 * 默认代码模板
 */
export const DEFAULT_CODES = {
  c: LANGUAGE_CONFIG.c.defaultCode,
  cpp: LANGUAGE_CONFIG.cpp.defaultCode,
  python: LANGUAGE_CONFIG.python.defaultCode,
  java: LANGUAGE_CONFIG.java.defaultCode,
  javascript: LANGUAGE_CONFIG.javascript.defaultCode,
  csharp: LANGUAGE_CONFIG.csharp.defaultCode
} as const

/**
 * 语言类型
 */
export type LanguageValue = keyof typeof LANGUAGE_CONFIG

/**
 * 编辑器主题类型
 */
export type EDITOR_THEHE =
  | 'chrome'
  | 'github'
  | 'monokai'
  | 'xcode'
  | 'dracula'
  | 'clouds'
  | 'terminal'
  | 'tomorrow'
  | 'twilight'
  | 'ambiance'
  | 'textmate'

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
