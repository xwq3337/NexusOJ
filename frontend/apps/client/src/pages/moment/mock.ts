// ============ 类型定义 ============
export interface MomentUser {
  id: number
  username: string
  nickname: string
  avatar: string
}

export interface MomentMedia {
  type: 'image' | 'video' | 'audio'
  url: string
  thumbnail?: string
  duration?: string
}

export interface MomentComment {
  id: number
  user: MomentUser
  content: string
  createdAt: string
  likes: number
}

export interface Moment {
  id: number
  user: MomentUser
  content: string
  tags: string[]
  media: MomentMedia[]
  likes: number
  comments: number
  shares: number
  bookmarks: number
  isLiked: boolean
  isBookmarked: boolean
  createdAt: string
  commentList: MomentComment[]
  type: 'recommend' | 'following' | 'latest'
}

// ============ Mock 用户 ============
const MOCK_USERS: MomentUser[] = [
  { id: 1, username: 'algorithm_master', nickname: '算法大师', avatar: 'https://picsum.photos/seed/u1/100/100' },
  { id: 2, username: 'acm_champion', nickname: 'ACM 冠军', avatar: 'https://picsum.photos/seed/u2/100/100' },
  { id: 3, username: 'code_ninja', nickname: '代码忍者', avatar: 'https://picsum.photos/seed/u3/100/100' },
  { id: 4, username: 'bug_hunter', nickname: 'Bug 猎人', avatar: 'https://picsum.photos/seed/u4/100/100' },
  { id: 5, username: 'stack_overflow', nickname: '栈溢出', avatar: 'https://picsum.photos/seed/u5/100/100' },
  { id: 6, username: 'dp_lover', nickname: '动态规划爱好者', avatar: 'https://picsum.photos/seed/u6/100/100' },
  { id: 7, username: 'binary_wizard', nickname: '二进制巫师', avatar: 'https://picsum.photos/seed/u7/100/100' },
]

// ============ Mock 动态数据 ============
export const MOCK_MOMENTS: Moment[] = [
  // ---- 视频 ----
  {
    id: 1,
    user: MOCK_USERS[1],
    content: '刚刚打完区域赛，分享一下最后一题的线段树优化思路！这道题当时全场只有 3 支队伍做出来，关键在于把区间操作转化为懒惰标记的合并 🏆',
    tags: ['ACM', '线段树', '区域赛'],
    media: [
      {
        type: 'video',
        url: 'https://cdn.plyr.io/static/demo/View_From_A_Blue_Moon_Trailer-576p.mp4',
        thumbnail: 'https://picsum.photos/seed/v1/800/450',
        duration: '2:30',
      },
    ],
    likes: 3842,
    comments: 256,
    shares: 128,
    bookmarks: 512,
    isLiked: false,
    isBookmarked: false,
    createdAt: '10 分钟前',
    commentList: [
      { id: 1, user: MOCK_USERS[0], content: '太强了！标记合并那步我完全没想到', createdAt: '8 分钟前', likes: 42 },
      { id: 2, user: MOCK_USERS[2], content: '请问这个能用树状数组替代吗？', createdAt: '5 分钟前', likes: 12 },
    ],
    type: 'recommend',
  },
  {
    id: 2,
    user: MOCK_USERS[0],
    content: '今天的每日一题，一道非常优雅的图论题。用 BFS 求最短路 + 状态压缩 DP，代码不到 50 行就搞定了！',
    tags: ['图论', 'BFS', '状态压缩'],
    media: [
      {
        type: 'video',
        url: 'https://www.w3schools.com/html/mov_bbb.mp4',
        thumbnail: 'https://picsum.photos/seed/v2/800/450',
        duration: '1:15',
      },
    ],
    likes: 1280,
    comments: 89,
    shares: 45,
    bookmarks: 230,
    isLiked: true,
    isBookmarked: true,
    createdAt: '30 分钟前',
    commentList: [
      { id: 3, user: MOCK_USERS[3], content: '状态压缩总是想不起来用，谢谢分享！', createdAt: '20 分钟前', likes: 28 },
    ],
    type: 'recommend',
  },

  // ---- 图片 ----
  {
    id: 3,
    user: MOCK_USERS[4],
    content: '实验室白板讨论了一个下午，终于把这个 NP-hard 问题的近似算法搞明白了。有时候回到纸笔推导比直接写代码效率高多了 ✍️',
    tags: ['算法', '近似算法'],
    media: [
      { type: 'image', url: 'https://picsum.photos/seed/m3a/800/600' },
      { type: 'image', url: 'https://picsum.photos/seed/m3b/800/600' },
      { type: 'image', url: 'https://picsum.photos/seed/m3c/800/600' },
    ],
    likes: 2156,
    comments: 178,
    shares: 92,
    bookmarks: 340,
    isLiked: false,
    isBookmarked: false,
    createdAt: '1 小时前',
    commentList: [
      { id: 4, user: MOCK_USERS[5], content: '白板推导是yyds！', createdAt: '45 分钟前', likes: 56 },
      { id: 5, user: MOCK_USERS[6], content: '这笔记做得比我论文还工整', createdAt: '30 分钟前', likes: 34 },
    ],
    type: 'latest',
  },
  {
    id: 4,
    user: MOCK_USERS[2],
    content: '从 AC 率 5% 到 100% 的进化史 📈 第一个版本超时 17 次，最后一个版本 0ms 通过，算法优化的魅力就在于此！',
    tags: ['优化', '性能'],
    media: [
      { type: 'image', url: 'https://picsum.photos/seed/m4a/600/800' },
      { type: 'image', url: 'https://picsum.photos/seed/m4b/600/800' },
    ],
    likes: 5672,
    comments: 423,
    shares: 289,
    bookmarks: 890,
    isLiked: true,
    isBookmarked: false,
    createdAt: '2 小时前',
    commentList: [
      { id: 6, user: MOCK_USERS[1], content: '从17次TLE到0ms，这才是真正的成长！', createdAt: '1.5 小时前', likes: 128 },
    ],
    type: 'recommend',
  },
  {
    id: 5,
    user: MOCK_USERS[3],
    content: '周末算法训练营结课啦！感谢所有参与者，大家的学习热情真的太高了 🔥 下一期我们将聚焦动态规划专题，敬请期待！',
    tags: ['训练营', '动态规划'],
    media: [
      { type: 'image', url: 'https://picsum.photos/seed/m5a/800/600' },
      { type: 'image', url: 'https://picsum.photos/seed/m5b/800/600' },
      { type: 'image', url: 'https://picsum.photos/seed/m5c/800/600' },
      { type: 'image', url: 'https://picsum.photos/seed/m5d/800/600' },
    ],
    likes: 8901,
    comments: 567,
    shares: 345,
    bookmarks: 1234,
    isLiked: false,
    isBookmarked: true,
    createdAt: '3 小时前',
    commentList: [
      { id: 7, user: MOCK_USERS[0], content: 'DP 专题安排上！我第一个报名', createdAt: '2 小时前', likes: 89 },
      { id: 8, user: MOCK_USERS[4], content: '上一期的回溯讲得太好了', createdAt: '1.5 小时前', likes: 45 },
    ],
    type: 'following',
  },
  {
    id: 6,
    user: MOCK_USERS[6],
    content: '今天在 NexusOJ 上刷到了一道神仙题，用到了 CDQ 分治 + 树状数组，分享一下我的解题思路和代码。这题难度评级 9/10，但想通了其实很优美。',
    tags: ['CDQ分治', '树状数组', 'NexusOJ'],
    media: [
      { type: 'image', url: 'https://picsum.photos/seed/m6/800/1000' },
    ],
    likes: 3402,
    comments: 234,
    shares: 167,
    bookmarks: 567,
    isLiked: false,
    isBookmarked: false,
    createdAt: '4 小时前',
    commentList: [
      { id: 9, user: MOCK_USERS[5], content: 'CDQ 分治每次看都觉得神奇', createdAt: '3 小时前', likes: 67 },
    ],
    type: 'latest',
  },

  // ---- 纯文本 ----
  {
    id: 7,
    user: MOCK_USERS[5],
    content: '给初学者的建议：不要一上来就刷困难题。先从简单题开始，把基础数据结构（栈、队列、堆、哈希表）搞明白，再逐步进阶到 DP 和图论。循序渐进才是正道 💡\n\n推荐刷题顺序：\n1. 数组/字符串基础\n2. 链表/栈/队列\n3. 二叉树/DFS/BFS\n4. 二分查找/滑动窗口\n5. 动态规划入门\n6. 图论基础',
    tags: ['经验分享', '刷题指南'],
    media: [],
    likes: 12450,
    comments: 892,
    shares: 2340,
    bookmarks: 5670,
    isLiked: true,
    isBookmarked: true,
    createdAt: '5 小时前',
    commentList: [
      { id: 10, user: MOCK_USERS[3], content: '收藏了！正好在迷茫该从哪里开始', createdAt: '4 小时前', likes: 234 },
      { id: 11, user: MOCK_USERS[0], content: '补充一下：刷题重在理解，不要背答案', createdAt: '3.5 小时前', likes: 156 },
    ],
    type: 'recommend',
  },
  {
    id: 8,
    user: MOCK_USERS[1],
    content: '明天就是 ICPC 亚洲区域赛了！我们队准备了整整三个月，从省赛铜牌到现在的目标——至少拿个银牌。不管结果如何，这段和大家一起刷题训练的日子是最宝贵的。加油！💪',
    tags: ['ICPC', '区域赛'],
    media: [],
    likes: 7892,
    comments: 567,
    shares: 890,
    bookmarks: 234,
    isLiked: false,
    isBookmarked: false,
    createdAt: '6 小时前',
    commentList: [
      { id: 12, user: MOCK_USERS[2], content: '加油！等你的好消息 🎉', createdAt: '5 小时前', likes: 345 },
    ],
    type: 'following',
  },

  // ---- 音频 ----
  {
    id: 9,
    user: MOCK_USERS[3],
    content: '🎙️ 算法小电台 EP.12：聊聊贪心算法——什么时候该用，什么时候不该用。很多人一看到最优解就上贪心，但其实贪心需要证明正确性！这期详细讲解了几个经典反例。',
    tags: ['播客', '贪心算法'],
    media: [
      {
        type: 'audio',
        url: 'https://www.w3schools.com/html/horse.mp3',
        duration: '15:32',
      },
    ],
    likes: 4521,
    comments: 312,
    shares: 178,
    bookmarks: 890,
    isLiked: false,
    isBookmarked: false,
    createdAt: '8 小时前',
    commentList: [
      { id: 13, user: MOCK_USERS[4], content: '音频讲得好清楚，终于理解为什么我的贪心总是WA了', createdAt: '7 小时前', likes: 89 },
    ],
    type: 'recommend',
  },
  {
    id: 10,
    user: MOCK_USERS[0],
    content: '🎵 周末刷题时的编程歌单分享。听音乐写代码真的能提高专注度，尤其是做那种需要深度思考的 DP 题时，来点 Lo-fi 效率拉满！',
    tags: ['音乐', '编程效率'],
    media: [
      {
        type: 'audio',
        url: 'https://www.w3schools.com/html/horse.mp3',
        duration: '3:45',
      },
    ],
    likes: 2345,
    comments: 167,
    shares: 89,
    bookmarks: 456,
    isLiked: true,
    isBookmarked: false,
    createdAt: '10 小时前',
    commentList: [
      { id: 14, user: MOCK_USERS[6], content: 'Lo-fi + DP = 完美组合', createdAt: '9 小时前', likes: 45 },
    ],
    type: 'latest',
  },

  // ---- 九宫格图片 ----
  {
    id: 11,
    user: MOCK_USERS[4],
    content: '校园算法社团纳新活动圆满结束！今年的新生质量真不错，好多同学连平衡树都会了 🌲 明天开始新生培训，从排序算法讲起～',
    tags: ['社团', '纳新'],
    media: [
      { type: 'image', url: 'https://picsum.photos/seed/m11a/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11b/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11c/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11d/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11e/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11f/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11g/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11h/400/400' },
      { type: 'image', url: 'https://picsum.photos/seed/m11i/400/400' },
    ],
    likes: 6780,
    comments: 445,
    shares: 234,
    bookmarks: 567,
    isLiked: false,
    isBookmarked: true,
    createdAt: '12 小时前',
    commentList: [
      { id: 15, user: MOCK_USERS[1], content: '新生就会平衡树？太卷了 😂', createdAt: '11 小时前', likes: 234 },
    ],
    type: 'following',
  },
  {
    id: 12,
    user: MOCK_USERS[6],
    content: '分享一个超赞的算法可视化网站，把红黑树的旋转操作做成了动画，一目了然！链接放评论区了 👇',
    tags: ['工具推荐', '红黑树', '可视化'],
    media: [
      {
        type: 'video',
        url: 'https://www.w3schools.com/html/movie.mp4',
        thumbnail: 'https://picsum.photos/seed/v12/800/450',
        duration: '0:45',
      },
    ],
    likes: 9123,
    comments: 678,
    shares: 4560,
    bookmarks: 3456,
    isLiked: false,
    isBookmarked: false,
    createdAt: '1 天前',
    commentList: [
      { id: 16, user: MOCK_USERS[0], content: '这个太实用了，终于搞懂红黑树了', createdAt: '23 小时前', likes: 567 },
      { id: 17, user: MOCK_USERS[2], content: '有 AVL 树的吗？', createdAt: '22 小时前', likes: 89 },
    ],
    type: 'recommend',
  },
]
