/**
 * 简单封装一下
 */
 var rsaUtil = {
  //RSA 位数，这里要跟后端对应
  bits: 1024,

  //当前JSEncrypted对象
  thisKeyPair: {},

  //生成密钥对(公钥和私钥)
  genKeyPair: function (bits = rsaUtil.bits) {
      let genKeyPair = {};
      rsaUtil.thisKeyPair = new JSEncrypt({default_key_size: bits});

      //获取私钥
      genKeyPair.privateKey = rsaUtil.thisKeyPair.getPrivateKey();

      //获取公钥
      genKeyPair.publicKey = rsaUtil.thisKeyPair.getPublicKey();

      return genKeyPair;
  },

  //公钥加密
  encrypt: function (plaintext, publicKey) {
      if (plaintext instanceof Object) {
          //1、JSON.stringify
          plaintext = JSON.stringify(plaintext)
      }
      publicKey && rsaUtil.thisKeyPair.setPublicKey(publicKey);
      return rsaUtil.thisKeyPair.encrypt(plaintext);
  },

  //私钥解密
  decrypt: function (ciphertext, privateKey) {
      privateKey && rsaUtil.thisKeyPair.setPrivateKey(privateKey);
      let decString = rsaUtil.thisKeyPair.decrypt(ciphertext);
      if(decString.charAt(0) === "{" || decString.charAt(0) === "[" ){
          //JSON.parse
          decString = JSON.parse(decString);
      }
      return decString;
  }
};

rsaUtil