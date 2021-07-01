export const menuList = [{
  name: '网关配置',
  id: 'Gateway',
  children: [{
    name: '映射配置',
    id: 'Overview',
    componentName: '/Overview'
  }, 
  // {
  //   name: '映射配置',
  //   id: 'Mapping',
  //   componentName: '/Mapping'
  // }, 
  {
    name: '插件配置',
    id: 'Plug',
    componentName: 'Plug'
  }]
}, {
  name: '限流配置',
  id: 'Flow',
  children: [{
    name: '限流配置',
    id: 'RateLimiter',
    componentName: '/RateLimiter'
  }]
}]
