import Cookies from 'js-cookie'

const TokenKey = 'ticket'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getLocalStorage (name) {
  
  var val = localStorage.getItem(name);
  return val ? JSON.parse(localStorage.getItem(name)) : ''
}

export function setLocalStorage (name,val){
  localStorage.setItem(name,JSON.stringify(val))
}

export function addLocalStorage (name,addVal){
  
  let oldVal = getLocalStorage(name) ? getLocalStorage(name) : [];
  let newVal = oldVal.concat(addVal);

  setLocalStorage(name,newVal);
}
export function deleteLocalStorage (name, deleteVal){
  let oldVal = getLocalStorage(name) ? getLocalStorage(name) : [];
  let _index = oldVal.findIndex((val) => val.time == deleteVal);

  if(_index > -1){
    oldVal.splice(_index,1);
  }
  setLocalStorage(name, oldVal);
}
export function clearLocalStorage(name){
  
  localStorage.removeItem(name);
}


//深拷贝
export const deepcopy = function (source) {
  if (!source) {
    return source;
  }
  let sourceCopy = source instanceof Array ? [] : {};
  for (let item in source) {
    sourceCopy[item] = typeof source[item] === 'object' ? deepcopy(source[item]) : source[item];
  }
  return sourceCopy;
};