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
export const safeJsonParse = (
  text: string,
  reviver?: (this: any, key: string, value: any) => any,
) => {
  let value: any
  let err: null | Error = null
  const isVaild = () => {
    return err === null
  }
  try {
    if (typeof text !== 'string' || text === 'null') {
      throw new Error('Invaild Json value')
    }
    value = JSON.parse(text, reviver)
  } catch (error: any) {
    err = error
  }
  return { value, err, isVaild }
}
