
function initWsClient(nickname) {
  nickname =  nickname === "" ? randomName() : nickname
  let wsClient = new WebSocket('ws://' + window.location.hostname + ':9998/ws');
  window.wsSingleton = wsClient
  wsClient.onopen = () => {
    console.log("connected");
    wsClient.send(JSON.stringify({
      data: JSON.stringify({
        ID: new Date().getTime(),
        Name: nickname
      })
    }));
    wsClient.onmessage = (e) => {
      rsp = JSON.parse(e.data)
      $("#send-input").val("")
      switch (rsp.data) {
        case "INTERACTIVE_SIGNAL_START":
          $("#send-input").attr("disabled",false).focus()
          break;
        case "INTERACTIVE_SIGNAL_STOP":
          $("#send-input").attr("disabled",true)
          break;
        default:
          $('<li>').html(rsp.data.replaceAll("\n","<br>")).appendTo($ul);
      }
    } 
  }
  var $ul = $('#msg-list');
}

function sendWsMsg(data){
  window.wsSingleton.send(JSON.stringify({data: String(data)}))
}

function randomName(prefix = "", randomLength = 4) {
  // 兼容更低版本的默认值写法
  prefix === undefined ? prefix = "" : prefix;
  randomLength === undefined ? randomLength = 8 : randomLength;

  // 设置随机用户名
  // 用户名随机词典数组
  let nameArr = [
      [1, 2, 3, 4, 5, 6, 7, 8, 9, 0],
      ["a", "b", "c", "d", "e", "f", "g", "h", "i", "g", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"]
  ]
  // 随机名字字符串
  let name = prefix;
  // 循环遍历从用户词典中随机抽出一个
  for (var i = 0; i < randomLength; i++) {
      // 随机生成index
      let index = Math.floor(Math.random() * 2);
      let zm = nameArr[index][Math.floor(Math.random() * nameArr[index].length)];
      // 如果随机出的是英文字母
      if (index === 1) {
          // 则百分之50的概率变为大写
          if (Math.floor(Math.random() * 2) === 1) {
              zm = zm.toUpperCase();
          }
      }
      // 拼接进名字变量中
      name += zm;
  }
  // 将随机生成的名字返回
  return name;
}
