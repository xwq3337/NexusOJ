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

export const SUPPORTED_LANGUAGES = ['C++', 'Java', 'Python', 'JavaScript', 'Go', 'Rust', 'C']
