// 关于错误处理和错误码的原则：
// 当你明确知道应该展示哪个错误码给客户端的时候，使用 output.Error(err,code,c)
// 当你获得一个错误，却不知道错误码，使用output.Auto(err,c)
//
// library和model应当返回兼容error的badcase错误，output库会判断如果是badcase error，自动找到错误码，记录日志和返回给客户端
// badcase新建错误方法:
//     import "sdk.look.360.cn/application/library/badcase"
//     err := badcase.New(错误信息 string, 错误码 int[，详细追溯数据 string])
// 其中追溯数据会记录到日志，不返回给客户端以保证敏感信息不外泄，客户端看到的是简单“错误信息 (错误码)”
//
// 谁写log谁报警？
// 如上，controller不需要显式调用log和报警模块，由output()统一控制策略
// library负责报警和记录日志，model只负责处理和获取数据并返回结果或badcase

package errno

const CHECK_UID = 2001     //check uid is error;
const CHECK_SIGN = 2002    //check sign is error;
const CHECK_CHANNEL = 2003 //check channel is error;
const CHECK_NUMBER = 2004  //check Number is error.value;
const CHECK_WORD = 2005    //check var is not word.value;
const CHECK_EMPTY = 2006   //check var is empty.value;
const CHECK_REQUEST = 2007 //check Request is error.value;
const CHECK_URLS = 2008    //check urls is error.value;
const CHECK_URL = 2009     //check url is error.value;
const CHECK_INLIST = 2010  //check param is in map/array

const REDIS_CALL_EMPTY = 1124 //redis call method got empty return data

const WEATHER_CACHE_ERROR = 1100 // bad weather cache data
const WEATHER_CITY_ERROR = 2100  // bad weather result due to city name
const WEATHER_DATA_ERROR = 2101  // bad weather result api data

const TABS_ENGINE_DATA = 1231 //tab from engine is error
const TABS_FILE_CACHE = 1235  //tab local file is error
const TABS_NOUPDATE = 2048    //tabs data no need to update

const LIST_ENGINE_DATA = 1218 //getData result is false.  c=youlike,channels except video
const LIST_VIDEO_DATA = 1221  //getData result is false.  c=video
const LIST_REDIS_CACHE = 1232 //data from redis is error
const LIST_TMPFS_CACHE = 1233 //data from tmpfs is error
const LIST_FILE_CACHE = 1236  //data from file is error

const LIST_ENGINE_URLCONF = 1237 //engine url conf is error

const SDK_LOCAL_CURL = 1201  //sdk local curl is error
const SDK_LOCAL_EMPTY = 3201 //sdk local curl is error
const SDK_LOCAL_ENC = 2301   //sdk local enc decode is error

const CITY_CODE_NOT_FOUND = 1058    //et city code error
const TOUTIAO_DATA_NOT_FOUND = 1060 //toutiao error
const CITY_LIST_NOT_FOUND = 1064    //toutiao error

const TOUTIAO_DATA_ERROR = 1318 //get toutiao list is error

const KUAISHIPIN_DATA_ERROR = 1320 //get kuaishipin video data is error
