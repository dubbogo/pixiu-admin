/* jshint esversion:9 */
/**websocket请求常量定义 */
const wsConsts = {
    /**websocket服务地址 */
    // url: 'ws://10.26.180.146:15479',
        url:'ws://127.0.0.1:15479',
   // url:'ws://192.168.10.156:15479',

    // url: 'ws://127.0.0.1:15479',
    camUrl: 'ws://127.0.0.1:15480',
    // camUrl: 'ws://192.168.10.156:15480',
    // url:'ws://192.168.10.176:15479',

    /**心跳间隔 */
    hbInterval: 10000,
    /**websocket接口映射列表 */
    methods: {
        /** 状态心跳 */
        heartbeat: 'zjetc.desktop.status',
        /** CPU卡号读取 */
        readCardId: 'zjetc.desktop.read-card-id',
        /** 公务卡信息读取 */
        readOfficialCard: 'zjetc.desktop.read-official-card',
        /** CPU信息读取 */
        readCpuInfo: 'zjetc.desktop.read-cpu-info',
        /** 读取卡内流水信息 */
        readCardJour: 'zjetc.desktop.read-card-jour',
        /** 卡片发行 */
        cpuIssue: 'zjetc.desktop.cpu-issue',
        /** 卡片注销 */
        cpuCancel: 'zjetc.desktop.cpu-cancel',
        /** 卡片解锁 */
        cpuUnlock: 'zjetc.desktop.cpu-unlock',
        /** OBU信息读取 */
        readObuInfo: 'zjetc.desktop.read-obu-info',
        /** OBU发行 */
        obuIssue: 'zjetc.desktop.obu-issue',
        /** 圈存 */
        cpuLoad: 'zjetc.desktop.cpu-load',
        /** 圈存异常处理 */
        cpuLoadAbnormal: 'zjetc.desktop.cpu-load-abnormal',
        /** 打开摄像头 */
        cameraOpen: 'zjetc.desktop.camera.open',
        /** 关闭摄像头 */
        cameraClose: 'zjetc.desktop.camera.close',
        /** 获取摄像头数量 */
        cameraCounts: 'zjetc.desktop.camera.counts',
        /** 拍照 */
        cameraTakePicture: 'zjetc.desktop.camera.take-picture',
        /** 换卡销卡 */
        cardReplaceCancel: 'zjetc.desktop.card-replace.cancel',
        /** 补领换卡新发行 */
        cardReplaceIssue: 'zjetc.desktop.card-replace.issue',
        /** 补领换卡新卡校验 */
        cardReplaceCheck: 'zjetc.desktop.card-replace.check',
        /**读取配置*/
        readConfig: 'zjetc.desktop.read-config',
        /**保存配置*/
        saveConfig: 'zjetc.desktop.save-config',
        /**打开cpu读写器*/
        openReader:'zjetc.desktop.cpu.open-reader',
        /**关闭读写器*/
        closeReader:'zjetc.desktop.cpu.close-reader',
        /**打开obu读写器*/
        opeRsu:'zjetc.desktop.open-rsu',
        /**关闭obu读写器*/
        closeRsu:'zjetc.desktop.close-rsu',
        posRecharge:'zjetc.desktop.pos-recharge',
        posCorrect:'zjetc.desktop.pos-correct'
        
    },
    /**到期时间 */
    expireTime: 3600000
}
export default wsConsts