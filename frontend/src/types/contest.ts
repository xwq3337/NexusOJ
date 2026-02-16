export interface Contest {
  id: string
  title: string
  startTime: string
  duration: string
  registered: number
  type: 'Weekly' | 'Biweekly' | 'Cup'
  status: 'Live' | 'Upcoming' | 'Ended'
}