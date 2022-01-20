import request from "@/utils/request";
import {
  getLocalStorage
} from "@/utils/storage";
import axios from "axios";

const { Secure } = require("mali-secure")
// ==== 修改 ====
var AppId = "20220119111341818363";
var secret = "22ade93aadcef15bb317e6e643ea053b";
var version = "1.0.1";
const secure = new Secure(AppId, 1, "", secret, version)
// ==== End ====

var routeMap = {
  List: "/api/list",
  Showlog: "/api/showlog",
  Deploy: "/api/deploy",
  Rollback: "/api/rollback",
};

// 获取映射 URL
function getUrl(url) {
  if (url.indexOf("http") >= 0) {
    return url
  }
  var uri = routeMap[url]
  return uri
}

export function ipost(path, data, headers = {}) {
  var uri = getUrl(path)
  uri = secure.getSign(uri)

  if (headers['Content-Type'] === undefined) {
    headers['Content-Type'] = "application/x-www-form-urlencoded"
  }

  return request({
    url: uri,
    method: "post",
    data,
    headers: headers
  });
}


export function iget(path, params) {
  let z = routeMap[path] || path;
  return request({
    url: z,
    method: "get",
    params: {
      params
    }
  });
}

export function upload(path, data, success, error) {
  let z = routeMap[path] || path;
  let token = getLocalStorage("token");

  let turl = process.env.VUE_APP_BASE_API
  if (location.protocol == 'http:') {
    turl = turl.replace("https", "http")
  }

  axios.post(turl + "/" + z, data, {
    timeout: 1800000,
    headers: {
      'Content-Type': 'multipart/form-data',
      'Authorization': "Bearer " + token,
    }
  }).then((data) => {
    success(data.data)
  }).catch(function (data) {
    error(data)
  });
}