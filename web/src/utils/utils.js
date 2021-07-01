/**
 * 一些工具类的集合
 */

//dataURL to blob,
export function dataUrlToBlob(dataURI){
  let mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]; // mime类型
  let byteString = atob(dataURI.split(',')[1]); //base64 解码

  let arrayBuffer = new ArrayBuffer(byteString.length); //创建缓冲数组
  let intArray = new Uint8Array(arrayBuffer); //创建视图

  for (let i = 0; i < byteString.length; i++) {
    intArray[i] = byteString.charCodeAt(i);
  }
  return new Blob([intArray], {type: mimeString});
}

export  function format(time, format){
  var t = new Date(time);
  var tf = function(i){return (i < 10 ? '0' : '') + i};
  return format.replace(/yyyy|MM|dd|HH|mm|ss/g, function(a){
    switch(a){
      case 'yyyy':
        return tf(t.getFullYear());
        break;
      case 'MM':
        return tf(t.getMonth() + 1);
        break;
      case 'mm':
        return tf(t.getMinutes());
        break;
      case 'dd':
        return tf(t.getDate());
        break;
      case 'HH':
        return tf(t.getHours());
        break;
      case 'ss':
        return tf(t.getSeconds());
        break;
    }
  })
}

export  function getArchivesName(code){
  let arr = {
    '1': '身份证正面',
    '11': '身份证反面',
    '3': '行驶证主页',
    '10': '行驶证附页',
    '6': '车头45度照',
  }
  return arr[code] ? arr[code] : '--'
}