export default {
  system: {
    advertisement: {
      title: '廣告列表',
      desc: '管理站台廣告位與投放內容',
      create: '新建廣告',
      delete: '刪除廣告',
      update: '更新廣告',
      adTitle: '廣告標題',
      adTitlePlaceholder: '請輸入廣告標題',
      adContent: '廣告內容',
      adContentPlaceholder: '請輸入廣告內容',
      image: '圖片',
      link: '連結',
      updateAction: '更新',
      confirmDeleteWithCount: '您已選取 [{count}] 項資源，刪除後將無法復原，是否確認刪除？'
    },
    appConfig: {
      site: '網站設定',
      system: '系統設定',
      email: '信箱設定',
      qiniu: '七牛雲設定',
      jwt: 'JWT 設定',
      gaode: '高德設定'
    },
    feedback: {
      title: '意見回饋列表',
      desc: '檢視並回覆使用者提交的回饋',
      delete: '刪除回饋',
      user: '使用者',
      reply: '回覆',
      replyTitle: '回覆回饋',
      confirmDeleteWithCount: '您已選取 [{count}] 項資源，刪除後將無法復原，是否確認刪除？'
    },
    friendLink: {
      title: '友鏈列表',
      desc: '管理本站對外展示的友情連結',
      create: '新建友鏈',
      delete: '刪除友鏈',
      update: '更新友鏈',
      name: '友鏈名稱',
      namePlaceholder: '請輸入友鏈名稱',
      description: '友鏈描述',
      descriptionPlaceholder: '請輸入友鏈描述',
      link: '連結',
      descCol: '描述',
      updateAction: '更新',
      confirmDeleteWithCount: '您已選取 [{count}] 項資源，刪除後將無法復原，是否確認刪除？'
    },
    loginLogs: {
      title: '登入紀錄',
      desc: '檢視使用者登入紀錄與登入狀態',
      uuidPlaceholder: '請輸入使用者 UUID',
      user: '使用者',
      username: '使用者名稱',
      loginTime: '登入時間',
      loginMethod: '登入方式',
      loginAddress: '登入位置',
      os: '作業系統',
      deviceInfo: '裝置資訊',
      browserInfo: '瀏覽器資訊',
      loginStatus: '登入狀態'
    },
    config: {
      email: {
        title: '信箱設定',
        host: '伺服器位置',
        port: '伺服器埠號',
        from: '寄件者信箱',
        nickname: '寄件者暱稱',
        secret: '信箱金鑰',
        ssl: '使用 SSL'
      },
      gaode: {
        title: '高德設定',
        enable: '是否開啟',
        key: '高德金鑰'
      },
      jwt: {
        title: 'JWT 設定',
        accessTokenSecret: '存取權杖金鑰',
        accessTokenExpiry: '存取權杖到期時間',
        refreshTokenSecret: '重新整理權杖金鑰',
        refreshTokenExpiry: '重新整理權杖到期時間',
        issuer: '簽發者'
      },
      qiniu: {
        title: '七牛雲設定',
        zone: '儲存區域',
        bucket: '空間名稱',
        accessKey: '金鑰 AK',
        secretKey: '金鑰 SK',
        cdnDomain: 'CDN 加速網域',
        useCdn: '使用 CDN 上傳加速',
        useHttps: '使用 Https'
      },
      qq: {},
      site: {
        websiteInfo: '網站資訊',
        logo: 'Logo 圖片',
        fullLogo: 'FullLogo 圖片',
        siteTitle: '網站標題',
        slogan: '網站標語',
        sloganEn: '英文標語',
        description: '網站描述',
        version: '網站版本',
        createdAt: '建立時間',
        icp: 'ICP 備案',
        police: '公安備案',
        bilibili: 'bilibili 連結',
        gitee: 'gitee 連結',
        github: 'github 連結',
        personalInfo: '個人資訊',
        nickname: '暱稱',
        job: '職業',
        address: '地址',
        email: '信箱',
        qqImage: 'QQ 圖片',
        wechatImage: '微信圖片',
        footerLink: '頁尾連結',
        create: '新建',
        title: '標題',
        link: '連結',
        homeImage: '首頁圖片',
        imageUpload: '圖片上傳',
        clearLogoConfirm: '是否清空 Logo 圖片？',
        clearFullLogoConfirm: '是否清空 FullLogo 圖片？',
        clearQqImageConfirm: '是否清空 QQ 圖片？',
        clearWechatImageConfirm: '是否清空微信圖片？',
        removeCarouselConfirm: '是否移除該首頁圖片？',
        operationCancelled: '已取消操作',
        warning: '警告',
        dragPrefix: '將檔案拖曳到此處，或',
        dragClick: '點擊上傳',
        uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超過 20MB。'
      },
      system: {
        title: '系統設定',
        multipoint: '多地點登入攔截',
        sessionSecret: '工作階段金鑰',
        ossType: '圖片儲存類型',
        selectPlaceholder: '請選擇',
        ossLocal: '本機',
        ossQiniu: '七牛'
      }
    },
    users: {
      title: '使用者列表',
      desc: '檢視站台註冊使用者並管理凍結狀態',
      uuidPlaceholder: '請輸入使用者 UUID',
      registerSource: '註冊來源',
      selectPlaceholder: '請選擇',
      registerEmail: '信箱',
      avatar: '頭像',
      username: '使用者名稱',
      address: '地址',
      registerTime: '註冊時間',
      role: '角色',
      admin: '管理員',
      normalUser: '一般使用者',
      freeze: '凍結',
      unfreeze: '解凍',
      freezeDialogTitle: '凍結/解凍使用者',
      freezeConfirm: '是否{action}該使用者：{username}'
    },
    userCenter: {
      comment: {
        title: '我的留言',
        desc: '檢視我發表過的全部留言'
      },
      feedback: {
        title: '我的意見回饋',
        desc: '檢視我提交的回饋及管理員回覆',
        reply: '回覆'
      },
      info: {
        title: '使用者資訊',
        userCard: '使用者卡片',
        avatar: '頭像',
        username: '使用者名稱',
        address: '地址',
        signature: '個人簡介',
        email: '信箱',
        role: '使用者權限',
        registerSource: '註冊來源',
        normalUser: '一般使用者',
        admin: '管理員',
        operation: '操作',
        changePassword: '修改密碼',
        usernameMax: '使用者名稱長度不應超過 20 字元',
        addressMax: '地址長度不應超過 200 字元',
        signatureMax: '個人簡介長度不應超過 320 字元'
      },
      star: {
        title: '我的收藏',
        desc: '檢視我收藏過的文章',
        category: '類別',
        abstract: '簡介',
        publishTime: '發布時間',
        articleId: '文章 id'
      }
    },
    login: {
      title: '登入'
    },
    error: {
      notFound: '抱歉，您造訪的頁面不存在或已被移除。',
      backHome: '返回首頁'
    }
  },
  forms: {
    login: {
      email: '信箱',
      password: '密碼',
      captcha: '驗證碼',
      emailPlaceholder: '請輸入信箱',
      passwordPlaceholder: '請輸入密碼',
      captchaPlaceholder: '請輸入驗證碼',
      submit: '登 入',
      emailFormat: '請輸入正確的信箱格式',
      passwordLength: '密碼長度應為 8~20 字元',
      captchaLength: '請輸入 6 位驗證碼'
    },
    register: {
      username: '使用者名稱',
      password: '密碼',
      confirmPassword: '確認密碼',
      email: '信箱',
      emailCaptcha: '信箱驗證碼',
      usernamePlaceholder: '請輸入使用者名稱',
      passwordPlaceholder: '請輸入密碼',
      confirmPasswordPlaceholder: '請再次輸入密碼',
      emailPlaceholder: '請輸入信箱',
      imageCaptchaPlaceholder: '請輸入圖片驗證碼',
      emailCaptchaPlaceholder: '請輸入信箱驗證碼',
      passwordMismatch: '兩次輸入的密碼不一致！',
      sendCode: '發送驗證碼',
      submit: '註冊',
      usernameMax: '使用者名稱長度不應超過 20 字元',
      passwordLength: '密碼長度應為 8~20 字元',
      emailFormat: '請輸入正確的信箱格式',
      codeLength: '請輸入 6 位驗證碼'
    },
    forgotPassword: {
      email: '信箱',
      emailCaptcha: '信箱驗證碼',
      password: '密碼',
      confirmPassword: '確認密碼',
      emailPlaceholder: '請輸入信箱',
      imageCaptchaPlaceholder: '請輸入圖片驗證碼',
      emailCaptchaPlaceholder: '請輸入信箱驗證碼',
      newPasswordPlaceholder: '請輸入新密碼',
      confirmNewPasswordPlaceholder: '請再次輸入新密碼',
      passwordMismatch: '兩次輸入的密碼不一致！',
      sendCode: '發送驗證碼',
      submit: '確定',
      emailFormat: '請輸入正確的信箱格式',
      codeLength: '請輸入 6 位驗證碼',
      passwordLength: '密碼長度應為 8~20 字元'
    },
    passwordReset: {
      oldPassword: '原密碼',
      newPassword: '新密碼',
      oldPasswordPlaceholder: '請輸入舊密碼',
      newPasswordPlaceholder: '請輸入新密碼',
      passwordLength: '密碼長度應為 8~20 字元'
    },
    articleCreate: {
      cover: '文章封面',
      title: '文章標題',
      category: '文章類別',
      tags: '文章標籤',
      abstract: '文章簡介',
      titlePlaceholder: '請輸入文章標題',
      categoryPlaceholder: '選擇分類',
      tagPlaceholder: '選擇標籤',
      tagPlaceholderNoCategory: '請先選擇分類',
      abstractPlaceholder: '請輸入文章簡介',
      dragPrefix: '將檔案拖曳到此處，或',
      dragClick: '點擊上傳',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超過 20MB。'
    },
    articleUpdate: {
      cover: '文章封面',
      title: '文章標題',
      category: '文章類別',
      tags: '文章標籤',
      abstract: '文章簡介',
      content: '文章內容',
      titlePlaceholder: '請輸入文章標題',
      categoryPlaceholder: '選擇分類',
      tagPlaceholder: '選擇標籤',
      tagPlaceholderNoCategory: '請先選擇分類',
      abstractPlaceholder: '請輸入文章簡介',
      editContent: '編輯內容',
      drawerFooter: '點擊上方 X 或外部任意區域即可結束編輯',
      dragPrefix: '將檔案拖曳到此處，或',
      dragClick: '點擊上傳',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超過 20MB。'
    },
    advertisementCreate: {
      image: '廣告圖片',
      link: '廣告連結',
      title: '廣告標題',
      content: '廣告內容',
      linkPlaceholder: '請輸入廣告連結',
      titlePlaceholder: '請輸入廣告標題',
      contentPlaceholder: '請輸入廣告內容',
      dragPrefix: '將檔案拖曳到此處，或',
      dragClick: '點擊上傳',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超過 20MB。'
    },
    advertisementUpdate: {
      image: '廣告圖片',
      link: '廣告連結',
      title: '廣告標題',
      content: '廣告內容',
      linkPlaceholder: '請輸入廣告連結',
      titlePlaceholder: '請輸入廣告標題',
      contentPlaceholder: '請輸入廣告內容'
    },
    friendLinkCreate: {
      logo: 'logo 圖片',
      link: '友鏈連結',
      name: '友鏈名稱',
      description: '友鏈描述',
      linkPlaceholder: '請輸入友鏈連結',
      namePlaceholder: '請輸入友鏈名稱',
      descriptionPlaceholder: '請輸入友鏈描述',
      dragPrefix: '將檔案拖曳到此處，或',
      dragClick: '點擊上傳',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超過 20MB。'
    },
    friendLinkUpdate: {
      logo: 'logo 圖片',
      link: '友鏈連結',
      name: '友鏈名稱',
      description: '友鏈描述',
      linkPlaceholder: '請輸入友鏈連結',
      namePlaceholder: '請輸入友鏈名稱',
      descriptionPlaceholder: '請輸入友鏈描述'
    },
    feedbackReply: {
      reply: '回饋回覆',
      placeholder: '請輸入回饋回覆'
    }
  }
}
