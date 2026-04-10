/* eslint-disable @typescript-eslint/no-explicit-any */
import { useClientEnv } from '@/stores/ClientEnv'
import { storeToRefs } from 'pinia'
const clientEnv = useClientEnv()
export const usePageSize = () => {
  const { innerWidth, innerHeight } = storeToRefs(clientEnv)
  return { innerWidth, innerHeight }
}
export const getRate = (rate: number) => '★★★★★☆☆☆☆☆'.slice(5 - rate, 10 - rate)
export const str = Math.random().toString(16).substring(2, 10)

