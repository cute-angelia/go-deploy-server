import moment from 'moment';

// 时间格式化，时间戳， 返回时间字符串
export function TimeFormatTimestamp(timestamp, format = "YYYY-MM-DD HH:mm:ss") {
  if (timestamp && (timestamp + "").length != 10) {
    return moment(timestamp).format(format)
  }

  return moment.unix(timestamp).format(format)
}


/** 根据分类获取标签， 用于展示如：表格
 * 
 * 
 vue:
 <el-tag
  :type="
    getEleTag(scope.row.initiative).color
  "
>
  {{
    getEleTag(scope.row.initiative).label
  }}
</el-tag>

function:
getEleTag(value) {
      return Util.getEleTag(value, {
        key: "value",
        value: "label"
        }, this.initiativeOptions, "");
    },
 
 data:
  initiativeOptions: [
        { value: 0, label: "一对多" },
        { value: 1, label: "1对1" }
      ],
*/
export function getEleTag(inValue, inField = {
  key: "id",
  value: "name"
}, options, index) {
  var rlabel = ""
  var rvalue = 0
  var rcolor = 0

  var colors = [
    'success',
    'warning',
    'info',
    'danger'
  ]

  var l = colors.length

  for (let z = 0; z < options.length; z++) {
    const element = options[z];
    if (element[inField.key] + "" === inValue + "") {
      rlabel = element[inField.value]
      rvalue = element[inField.key] + ""
      rcolor = colors[z % l]
      break;
    }
  }

  return {
    color: rcolor,
    label: rlabel,
    value: rvalue
  }
}

// 用于展示在表格， 根据键值对获取
export function getObjectByValueInOptions(inValue, inField = {
  key: "id",
  value: "name"
}, options) {
  var resp = inValue
  for (let z = 0; z < options.length; z++) {
    const element = options[z];
    if (element[inField.key] * 1 === inValue * 1) {
      resp = element[inField.value]
      break;
    }
  }
  return resp
}

// 标签编辑处理，前后端的处理不一致，前端select 后端 按，分隔
export function formatTagsToArray(tags) {
  var t = [];
  if (typeof tags == "object") {
    return tags;
  }
  // tags == ""
  // if (tags.length === 0) {
  //   t = ["添加"];
  //   return t;
  // }

  // tags 有逗号
  if (tags && tags.length > 0) {
    t = tags.split(",");
  }

  return trimSpace(t);
}

// 排除空字符数组
export function trimSpace(array) {
  for (var i = 0; i < array.length; i++) {
    if (array[i] == "" || array[i] == " " || array[i] == null || typeof (array[i]) == "undefined") {
      array.splice(i, 1);
      i = i - 1;

    }
  }
  return array;
}

export function getExt(uri) {
  return uri.split('.').pop().toLowerCase().split("?").shift()
}

export function isMedia(uri) {
  var ext = getExt(uri)
  if (ext == "mp4" || ext == "avi" || ext == "mov") {
    return true
  } else {
    return false
  }
}

export function isImage(uri) {
  var ext = getExt(uri)
  if (ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "webp" || ext == "svg") {
    return true
  } else {
    return false
  }
}

export function isNovie(uri) {
  var ext = getExt(uri)
  if (ext == "epub" || ext == "txt") {
    return true
  } else {
    return false
  }
}


export function downloadUrlFile(url, filename) {
  const xhr = new XMLHttpRequest();
  xhr.open('GET', url, true);
  xhr.responseType = 'blob'
  //xhr.setRequestHeader('Authorization'', 'Basic a2VybWl0Omtlcm1pdA =='');
  xhr.onload = () => {
    if (xhr.status === 200) {
      // 获取图片blob数据并保存
      saveAs(xhr.response, filename);
    }
  }; xhr.send();
}

function saveAs(blob, filename) {
  if (window.navigator.msSaveOrOpenBlob) {
    navigator.msSaveBlob(blob, filename);
  } else {
    var link = document.createElement('a');
    var body = document.querySelector('body');

    link.href = window.URL.createObjectURL(blob);
    link.download = filename;

    // fix Firefox
    link.style.display = 'none';
    body.appendChild(link);

    link.click();
    body.removeChild(link);

    window.URL.revokeObjectURL(link.href);
  }
}