/**
 * IndexedDB Storage Service
 * 用于存储题目代码和测试用例
 */

import { DEFAULT_CODES } from "../../constants"

const DB_NAME = 'NexusOJ_DB'
const DB_VERSION = 3
const STORE_CODE = 'problem_code'
const STORE_TEST_CASES = 'test_cases'

/**
 * 语言代码映射
 */
interface LanguageCodeMap {
  [language: string]: string // language -> code
}

/**
 * 代码存储结构 - 每个题目一条记录，包含所有语言的代码
 * key: problemId
 * 例如: "1", "2", "3"
 */
interface CodeRecord {
  id: string // problemId
  languageCodeMap: LanguageCodeMap // 语言代码映射: { "cpp": "...", "python": "..." }
  updated_at: number
  created_at?: number
}

interface TestCaseRecord {
  id: string // problem_id
  input: string
  expected: string
  updated_at: number
}

class IndexedDBService {
  private db: IDBDatabase | null = null

  /**
   * 初始化数据库
   */
  async init(): Promise<void> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION)

      request.onerror = () => {
        reject(new Error('Failed to open IndexedDB'))
      }

      request.onsuccess = () => {
        this.db = request.result
        resolve()
      }

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result
        const transaction = request.transaction

        // 创建代码存储 (新版本 - 一个题目一条记录，包含所有语言)
        if (!db.objectStoreNames.contains(STORE_CODE)) {
          const codeStore = db.createObjectStore(STORE_CODE, { keyPath: 'id' })
          codeStore.createIndex('updated_at', 'updated_at', { unique: false })
        } else if (transaction) {
          // 如果存储已存在，确保索引存在
          const codeStore = transaction.objectStore(STORE_CODE)
          if (!codeStore.indexNames.contains('updated_at')) {
            codeStore.createIndex('updated_at', 'updated_at', { unique: false })
          }
        }

        // 创建测试用例存储
        if (!db.objectStoreNames.contains(STORE_TEST_CASES)) {
          const testCaseStore = db.createObjectStore(STORE_TEST_CASES, { keyPath: 'id' })
          testCaseStore.createIndex('updated_at', 'updated_at', { unique: false })
        } else if (transaction) {
          const testCaseStore = transaction.objectStore(STORE_TEST_CASES)
          if (!testCaseStore.indexNames.contains('updated_at')) {
            testCaseStore.createIndex('updated_at', 'updated_at', { unique: false })
          }
        }
      }
    })
  }

  /**
   * 确保数据库已初始化
   */
  private async ensureDB(): Promise<void> {
    if (!this.db) {
      await this.init()
    }
  }

  /**
   * 保存题目代码
   * @param problemId 题目ID
   * @param language 编程语言
   * @param code 代码内容
   */
  async saveCode(problemId: string, language: string, code: string): Promise<void> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readwrite')
      const store = transaction.objectStore(STORE_CODE)

      // 先获取现有记录
      const getRequest = store.get(problemId)

      getRequest.onsuccess = () => {
        const existingRecord = getRequest.result
        let languageCodeMap: LanguageCodeMap = {}

        // 如果记录存在，保留其他语言的代码
        if (existingRecord) {
          languageCodeMap = existingRecord.languageCodeMap || {}
        }

        // 更新当前语言的代码
        languageCodeMap[language] = code

        const record: CodeRecord = {
          id: problemId,
          languageCodeMap,
          updated_at: Date.now(),
          created_at: existingRecord?.created_at || Date.now()
        }

        const putRequest = store.put(record)

        putRequest.onsuccess = () => resolve()
        putRequest.onerror = () => reject(new Error('Failed to save code'))
      }

      getRequest.onerror = () => reject(new Error('Failed to get existing code'))
    })
  }

  /**
   * 默认代码模板
   */
  private readonly DEFAULT_CODES: Record<string, string> = DEFAULT_CODES

  /**
   * 获取默认代码
   * @param language 编程语言
   * @returns 默认代码
   */
  getDefaultCode(language: string): string {
    return this.DEFAULT_CODES[language] || '// Your code here'
  }

  /**
   * 获取题目指定语言的代码
   * @param problemId 题目ID
   * @param language 编程语言
   * @returns 代码内容，如果不存在返回 null
   */
  async getCode(problemId: string, language: string): Promise<{ code: string } | null> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readonly')
      const store = transaction.objectStore(STORE_CODE)
      const request = store.get(problemId)

      request.onsuccess = () => {
        const record = request.result
        // 明确检查：记录存在、languageCodeMap 存在、键存在于 map 中、且值是非空字符串
        if (record && record.languageCodeMap && Object.prototype.hasOwnProperty.call(record.languageCodeMap, language) && record.languageCodeMap[language] && record.languageCodeMap[language].trim().length > 0) {
          console.log(':', record.languageCodeMap[language])
          resolve({ code: record.languageCodeMap[language] })
        } else {
          // 没有找到代码，返回 null，由调用方决定是否使用默认代码
          resolve(null)
        }
      }
      request.onerror = () => reject(new Error('Failed to get code'))
    })
  }

  /**
   * 获取题目的所有语言代码
   * @param problemId 题目ID
   * @returns 包含所有语言代码的记录
   */
  async getAllCodeForProblem(problemId: string): Promise<CodeRecord | null> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readonly')
      const store = transaction.objectStore(STORE_CODE)
      const request = store.get(problemId)

      request.onsuccess = () => {
        resolve(request.result || null)
      }
      request.onerror = () => reject(new Error('Failed to get all codes for problem'))
    })
  }

  /**
   * 删除题目指定语言的代码
   * @param problemId 题目ID
   * @param language 编程语言
   */
  async deleteCode(problemId: string, language: string): Promise<void> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readwrite')
      const store = transaction.objectStore(STORE_CODE)

      // 先获取现有记录
      const getRequest = store.get(problemId)

      getRequest.onsuccess = () => {
        const record = getRequest.result

        if (!record || !record.languageCodeMap) {
          resolve()
          return
        }

        // 删除指定语言的代码
        delete record.languageCodeMap[language]
        record.updated_at = Date.now()

        // 如果没有语言代码了，删除整个记录
        if (Object.keys(record.languageCodeMap).length === 0) {
          const deleteRequest = store.delete(problemId)
          deleteRequest.onsuccess = () => resolve()
          deleteRequest.onerror = () => reject(new Error('Failed to delete code'))
        } else {
          const putRequest = store.put(record)
          putRequest.onsuccess = () => resolve()
          putRequest.onerror = () => reject(new Error('Failed to update code'))
        }
      }

      getRequest.onerror = () => reject(new Error('Failed to get existing code'))
    })
  }

  /**
   * 删除题目的所有语言代码
   * @param problemId 题目ID
   */
  async deleteAllCodeForProblem(problemId: string): Promise<void> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readwrite')
      const store = transaction.objectStore(STORE_CODE)
      const request = store.delete(problemId)

      request.onsuccess = () => resolve()
      request.onerror = () => reject(new Error('Failed to delete all codes'))
    })
  }

  /**
   * 保存测试用例
   * @param problemId 题目ID
   * @param input 输入
   * @param expected 期望输出
   */
  async saveTestCase(problemId: string, input: string, expected: string): Promise<void> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_TEST_CASES], 'readwrite')
      const store = transaction.objectStore(STORE_TEST_CASES)

      const record: TestCaseRecord = {
        id: problemId,
        input,
        expected,
        updated_at: Date.now()
      }

      const request = store.put(record)

      request.onsuccess = () => resolve()
      request.onerror = () => reject(new Error('Failed to save test case'))
    })
  }

  /**
   * 获取测试用例
   * @param problemId 题目ID
   * @returns 测试用例记录
   */
  async getTestCase(problemId: string): Promise<TestCaseRecord | null> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_TEST_CASES], 'readonly')
      const store = transaction.objectStore(STORE_TEST_CASES)
      const request = store.get(problemId)

      request.onsuccess = () => {
        resolve(request.result || null)
      }
      request.onerror = () => reject(new Error('Failed to get test case'))
    })
  }

  /**
   * 删除测试用例
   * @param problemId 题目ID
   */
  async deleteTestCase(problemId: string): Promise<void> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_TEST_CASES], 'readwrite')
      const store = transaction.objectStore(STORE_TEST_CASES)
      const request = store.delete(problemId)

      request.onsuccess = () => resolve()
      request.onerror = () => reject(new Error('Failed to delete test case'))
    })
  }

  /**
   * 清空所有数据
   */
  async clearAll(): Promise<void> {
    await this.ensureDB()

    const stores = [STORE_CODE, STORE_TEST_CASES]

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(stores, 'readwrite')

      stores.forEach((storeName) => {
        const store = transaction.objectStore(storeName)
        store.clear()
      })

      transaction.oncomplete = () => resolve()
      transaction.onerror = () => reject(new Error('Failed to clear all data'))
    })
  }

  /**
   * 获取所有题目代码
   * @returns 所有代码记录
   */
  async getAllCodes(): Promise<CodeRecord[]> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_CODE], 'readonly')
      const store = transaction.objectStore(STORE_CODE)
      const request = store.getAll()

      request.onsuccess = () => resolve(request.result || [])
      request.onerror = () => reject(new Error('Failed to get all codes'))
    })
  }

  /**
   * 获取所有测试用例
   * @returns 所有测试用例记录
   */
  async getAllTestCases(): Promise<TestCaseRecord[]> {
    await this.ensureDB()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_TEST_CASES], 'readonly')
      const store = transaction.objectStore(STORE_TEST_CASES)
      const request = store.getAll()

      request.onsuccess = () => resolve(request.result || [])
      request.onerror = () => reject(new Error('Failed to get all test cases'))
    })
  }

  /**
   * 导出所有数据
   * @returns 包含所有代码和测试用例的数据
   */
  async exportData(): Promise<{
    codes: CodeRecord[]
    testCases: TestCaseRecord[]
  }> {
    const [codes, testCases] = await Promise.all([
      this.getAllCodes(),
      this.getAllTestCases()
    ])

    return { codes, testCases }
  }

  /**
   * 导入数据
   * @param data 要导入的数据
   */
  async importData(data: {
    codes?: CodeRecord[]
    testCases?: TestCaseRecord[]
  }): Promise<void> {
    await this.ensureDB()

    const { codes = [], testCases = [] } = data

    // 导入代码
    const codeTransaction = this.db!.transaction([STORE_CODE], 'readwrite')
    const codeStore = codeTransaction.objectStore(STORE_CODE)

    for (const code of codes) {
      codeStore.put(code)
    }

    // 导入测试用例
    const testCaseTransaction = this.db!.transaction([STORE_TEST_CASES], 'readwrite')
    const testCaseStore = testCaseTransaction.objectStore(STORE_TEST_CASES)

    for (const testCase of testCases) {
      testCaseStore.put(testCase)
    }

    await new Promise((resolve) => {
      testCaseTransaction.oncomplete = () => resolve(null)
    })
  }
}

// 导出单例
export const indexedDBService = new IndexedDBService()
export default indexedDBService
