import type { Problem, Contest } from '@nexusoj/type'
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

export interface UserRanking {
  rank: number
  username: string
  nickname: string
  avatar: string
  rating: number
  school?: string
}

export const MOCK_USER_RANKING: UserRanking[] = [
  {
    rank: 1,
    username: 'algorithm_master',
    nickname: '算法大师',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=algorithm',
    rating: 2847,
    school: 'MIT'
  },
  {
    rank: 2,
    username: 'code_ninja',
    nickname: '代码忍者',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=ninja',
    rating: 2756,
    school: 'Stanford'
  },
  {
    rank: 3,
    username: 'byte_warrior',
    nickname: '字节战士',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=warrior',
    rating: 2689,
    school: 'Tsinghua'
  },
  {
    rank: 4,
    username: 'pixel_coder',
    nickname: '像素程序员',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=pixel',
    rating: 2623,
    school: 'CMU'
  },
  {
    rank: 5,
    username: 'data_drifter',
    nickname: '数据漂移者',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=drifter',
    rating: 2587,
    school: 'Berkeley'
  },
  {
    rank: 6,
    username: 'logic_lord',
    nickname: '逻辑领主',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=logic',
    rating: 2542,
    school: 'Oxford'
  },
  {
    rank: 7,
    username: 'syntax_sage',
    nickname: '语法圣贤',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=syntax',
    rating: 2498,
    school: 'ETH Zurich'
  },
  {
    rank: 8,
    username: 'quantum_dev',
    nickname: '量子开发者',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=quantum',
    rating: 2456,
    school: 'Caltech'
  },
  {
    rank: 9,
    username: 'graph_guru',
    nickname: '图论大师',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=graph',
    rating: 2412,
    school: 'MIT'
  },
  {
    rank: 10,
    username: 'dp_dynamite',
    nickname: 'DP达人',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=dp',
    rating: 2389,
    school: 'Waterloo'
  }
]
