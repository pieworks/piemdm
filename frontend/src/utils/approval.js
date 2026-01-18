import { useApprovalStore } from '@/pinia/modules/approval'
//  获取字典方法 使用示例 getDict('sex').then(res)  或者 async函数下 const res = await getDict('sex')
export const getDict = async(type) => {
  const approvalStore = useApprovalStore()
  await approvalStore.getApproval(type)
  return approvalStore.approvalMap[type]
}

//  字典文字展示方法
export const showDictLabel = (dict, code) => {
  if (!dict) {
    return ''
  }
  const dictMap = {}
  dict.forEach(item => {
    dictMap[item.value] = item.label
  })
  return dictMap[code]
}
