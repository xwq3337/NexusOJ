// Utils Package
// Export shared utility functions here
export const safeJsonParse = (
  text: string,
  reviver?: (this: any, key: string, value: any) => any,
) => {
  let data: any
  let err: null | Error = null
  const isVaild = () => {
    return err === null
  }
  try {
    if (typeof text !== 'string' || text === 'null') {
      throw new Error('Invaild Json value')
    }
    data = JSON.parse(text, reviver)
  } catch (error: any) {
    err = error
  }
  return { data, err, isVaild }
}
export {

}
