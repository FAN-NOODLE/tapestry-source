const getKsDays = (startDate, allDay = 30) => {
  let dateList = []
  let endDate = new Date()
  endDate.setDate(startDate.getDate() + allDay)
  while ((endDate.getTime() - startDate.getTime()) >= 0) {
    let month = (startDate.getMonth() + 1).toString().length === 1 ? '0' + (startDate.getMonth() + 1).toString() : (
      startDate.getMonth() + 1)
    let day = startDate.getDate().toString().length === 1 ? '0' + startDate.getDate() : startDate.getDate()
    let value = `${month}月${day}日`
    let dateString = `${startDate.getFullYear()}/${month}/${day} 10:00:00`
    let obj = {
      value,
      select_name: `${startDate.getFullYear()}-${month}-${day}`,
      date_list: getTimeSlot(dateString)
    }
    dateList.push(obj)
    startDate.setDate(startDate.getDate() + 1)
  }
  return dateList
}
// 取10点-22点的时间段
const getTimeSlot = value => {
  const text = [
    '10:00 - 11:00',
    '11:00 - 12:00',
    '12:00 - 13:00',
    '13:00 - 14:00',
    '14:00 - 15:00',
    '15:00 - 16:00',
    '16:00 - 17:00',
    '17:00 - 18:00',
    '18:00 - 19:00',
    '19:00 - 20:00',
    '20:00 - 21:00',
    '21:00 - 22:00'
  ]
  const nowDate = new Date().getTime()
  return text.map((item, index) => {
    const startTime = new Date(value).getTime() + 3600000 * index
    return {
      startTime: startTime / 1000,
      name: item,
      isFinish: startTime < nowDate - 3600000,
      inviteNum: 0
    }
  })
}
const getDays = (startDate, allDay = 30,timeSlot) => {
  let dateList = []
  let endDate = new Date(new Date().getTime() + 1000 * 60 * 60 * 24)
  endDate.setDate(startDate.getDate() + allDay)
  while ((endDate.getTime() - startDate.getTime()) >= 0) {
    let month = (startDate.getMonth() + 1).toString().length === 1 ? '0' + (startDate.getMonth() + 1).toString() : (
      startDate.getMonth() + 1)
    let day = startDate.getDate().toString().length === 1 ? '0' + startDate.getDate() : startDate.getDate()
    let value = `${month}月${day}日`
    let dateString = `${startDate.getFullYear()}/${month}/${day} 10:00:00`
    let obj = {
      value,
      select_name: `${startDate.getFullYear()}-${month}-${day}`,
      date_list: getCoachTimeSlot(dateString,timeSlot)
    }
    dateList.push(obj)
    startDate.setDate(startDate.getDate() + 1)
  }
  return dateList
}
//coach获取指定时间段
const getCoachTimeSlot = (value,timeSlot)=> {
  // 取出value中的年月日
  const year = value.split('/')[0]
  const month = value.split('/')[1]
  const day = value.split('/')[2].split(' ')[0]
  // 取出timeSlot中的时间段,start_time为开始时间,end_time为结束时间,时间格式为start_time - end_time
  let text
  if(timeSlot && timeSlot.length > 0){
    text = timeSlot.map(item => {
      return item.start_time + ' - ' + item.end_time
    })
  }else{
    text = [
      '10:00 - 11:00',
      '11:00 - 12:00',
      '12:00 - 13:00',
      '13:00 - 14:00',
      '14:00 - 15:00',
      '15:00 - 16:00',
      '16:00 - 17:00',
      '17:00 - 18:00',
      '18:00 - 19:00',
      '19:00 - 20:00',
      '20:00 - 21:00',
      '21:00 - 22:00'
    ]
  }

  const nowDate = new Date().getTime()
  return text.map((item, index) => {
    // 将value中的年月日和timeSlot中的时间段拼接成时间戳
    const startTime = new Date(`${year}/${month}/${day} ${item.split(' - ')[0]}`).getTime()
    return {
      startTime: startTime / 1000,
      name: item,
      isFinish: startTime < nowDate - 3600000,
      inviteNum: 0
    }
  })
}
export { getDays, getKsDays }

