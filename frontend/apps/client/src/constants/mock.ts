export const MOCK_PROBLEMS = [
  {
    id: 1,
    title: '两数之和',
    difficulty: 0,
    tags: ['数组', '哈希表'],
    user_id: 3953240110606292,
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
    user_id: 3953240110606292,
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
    user_id: 3953240110606292,
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
    user_id: 3953240110606292,
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
    user_id: 3953240110606292,
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
    user_id: 3953240110606292,
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

export const MOCK_CONTESTS = [
  {
    id: '792516153393221',
    title: '周赛 1',
    begin_at: '2025/12/15 14:44',
    duration: 90,
    participants: 12450,
    status: 'Live'
  },
  {
    id: '792516153393221',
    title: '双周赛 1',
    begin_at: '2025/12/16 14:44',
    duration: 90,
    participants: 4300,
    status: 'Upcoming'
  },
  {
    id: '792516153393221',
    title: 'Nexus新手杯 #1',
    begin_at: '2025/12/17 14:44',
    duration: 120,
    participants: 1200,
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


export const ACTIVITY_DATA = [
  { name: 'Mon', submissions: 400 },
  { name: 'Tue', submissions: 300 },
  { name: 'Wed', submissions: 550 },
  { name: 'Thu', submissions: 500 },
  { name: 'Fri', submissions: 700 },
  { name: 'Sat', submissions: 800 },
  { name: 'Sun', submissions: 750 }
]

export interface UserRanking {
  rank: number
  username: string
  nickname: string
  avatar: string
  rating: number
  school?: string
}
