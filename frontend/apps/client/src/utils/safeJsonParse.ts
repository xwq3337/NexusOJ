export const safeJsonParse = (
  text: string,
  reviver?: (this: any, key: string, value: any) => any
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
