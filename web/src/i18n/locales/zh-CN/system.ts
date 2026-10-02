export default {
  system: {
    advertisement: {
      title: '广告列表',
      desc: '管理站点广告位与投放内容',
      create: '新建广告',
      delete: '删除广告',
      update: '更新广告',
      adTitle: '广告标题',
      adTitlePlaceholder: '请输入广告标题',
      adContent: '广告内容',
      adContentPlaceholder: '请输入广告内容',
      image: '图片',
      link: '链接',
      updateAction: '更新',
      confirmDeleteWithCount: '您已选中 [{count}] 项资源，删除后将无法恢复，是否确认删除？'
    },
    appConfig: {
      site: '网站配置',
      system: '系统配置',
      email: '邮箱配置',
      qiniu: '七牛云配置',
      jwt: 'JWT 配置',
      gaode: '高德配置'
    },
    feedback: {
      title: '反馈列表',
      desc: '查看并回复用户提交的反馈',
      delete: '删除反馈',
      user: '用户',
      reply: '回复',
      replyTitle: '回复反馈',
      confirmDeleteWithCount: '您已选中 [{count}] 项资源，删除后将无法恢复，是否确认删除？'
    },
    friendLink: {
      title: '友链列表',
      desc: '管理本站对外展示的友情链接',
      create: '新建友链',
      delete: '删除友链',
      update: '更新友链',
      name: '友链名称',
      namePlaceholder: '请输入友链名称',
      description: '友链描述',
      descriptionPlaceholder: '请输入友链描述',
      link: '链接',
      descCol: '描述',
      updateAction: '更新',
      confirmDeleteWithCount: '您已选中 [{count}] 项资源，删除后将无法恢复，是否确认删除？'
    },
    loginLogs: {
      title: '登录日志',
      desc: '查看用户登录记录与登录状态',
      uuidPlaceholder: '请输入用户UUID',
      user: '用户',
      username: '用户名',
      loginTime: '登录时间',
      loginMethod: '登录方式',
      loginAddress: '登录地址',
      os: '操作系统',
      deviceInfo: '设备信息',
      browserInfo: '浏览器信息',
      loginStatus: '登录状态'
    },
    config: {
      email: {
        title: '邮箱配置',
        host: '服务器地址',
        port: '服务器端口',
        from: '发件人邮箱',
        nickname: '发件人昵称',
        secret: '邮箱密钥',
        ssl: '使用SSL'
      },
      gaode: {
        title: '高德配置',
        enable: '是否开启',
        key: '高德密钥'
      },
      jwt: {
        title: 'JWT 配置',
        accessTokenSecret: '访问令牌密钥',
        accessTokenExpiry: '访问令牌过期时间',
        refreshTokenSecret: '刷新令牌密钥',
        refreshTokenExpiry: '刷新令牌过期时间',
        issuer: '签发者'
      },
      qiniu: {
        title: '七牛云配置',
        zone: '存储区域',
        bucket: '空间名称',
        accessKey: '密钥 AK',
        secretKey: '密钥 SK',
        cdnDomain: 'CDN加速域名',
        useCdn: '使用CDN上传加速',
        useHttps: '使用Https'
      },
      qq: {},
      site: {
        websiteInfo: '网站信息',
        logo: 'Logo图片',
        fullLogo: 'FullLogo图片',
        siteTitle: '网站标题',
        slogan: '网站标语',
        sloganEn: '英文标语',
        description: '网站描述',
        version: '网站版本',
        createdAt: '创建时间',
        icp: 'ICP备案',
        police: '公安备案',
        bilibili: 'bilibili链接',
        gitee: 'gitee链接',
        github: 'github链接',
        personalInfo: '个人信息',
        nickname: '昵称',
        job: '职业',
        address: '地址',
        email: '邮箱',
        qqImage: 'QQ图片',
        wechatImage: '微信图片',
        footerLink: '页脚链接',
        create: '新建',
        title: '标题',
        link: '链接',
        homeImage: '首页图片',
        imageUpload: '图片上传',
        clearLogoConfirm: '是否清空Logo图片？',
        clearFullLogoConfirm: '是否清空FullLogo图片？',
        clearQqImageConfirm: '是否清空QQ图片？',
        clearWechatImageConfirm: '是否清空微信图片？',
        removeCarouselConfirm: '是否移除该首页图片？',
        operationCancelled: '操作取消',
        warning: '警告',
        dragPrefix: '将文件拖到此处，或',
        dragClick: '点击上传',
        uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超过 20MB。'
      },
      system: {
        title: '系统配置',
        multipoint: '多地点登录拦截',
        sessionSecret: '会话密钥',
        ossType: '图片存储类型',
        selectPlaceholder: '请选择',
        ossLocal: '本地',
        ossQiniu: '七牛'
      }
    },
    users: {
      title: '用户列表',
      desc: '查看站点注册用户并管理冻结状态',
      uuidPlaceholder: '请输入用户UUID',
      registerSource: '注册来源',
      selectPlaceholder: '请选择',
      registerEmail: '邮箱',
      avatar: '头像',
      username: '用户名',
      address: '地址',
      registerTime: '注册时间',
      role: '角色',
      admin: '管理员',
      normalUser: '普通用户',
      freeze: '冻结',
      unfreeze: '解冻',
      freezeDialogTitle: '冻结/解冻用户',
      freezeConfirm: '是否{action}该用户：{username}'
    },
    userCenter: {
      comment: {
        title: '我的评论',
        desc: '查看我发表过的全部评论'
      },
      feedback: {
        title: '我的反馈',
        desc: '查看我提交的反馈及管理员回复',
        reply: '回复'
      },
      info: {
        title: '用户信息',
        userCard: '用户卡片',
        avatar: '头像',
        username: '用户名',
        address: '地址',
        signature: '签名',
        email: '邮箱',
        role: '用户权限',
        registerSource: '注册来源',
        normalUser: '普通用户',
        admin: '管理员',
        operation: '操作',
        changePassword: '修改密码',
        usernameMax: '用户名长度不应大于20位',
        addressMax: '地址长度不应大于200位',
        signatureMax: '签名长度不应大于320位'
      },
      star: {
        title: '我的收藏',
        desc: '查看我收藏过的文章',
        category: '类别',
        abstract: '简介',
        publishTime: '发布时间',
        articleId: '文章id'
      }
    },
    login: {
      title: '登录'
    },
    error: {
      notFound: '抱歉，您访问的页面不存在或已被移除。',
      backHome: '返回首页'
    }
  },
  forms: {
    login: {
      email: '邮箱',
      password: '密码',
      captcha: '验证码',
      emailPlaceholder: '请输入邮箱',
      passwordPlaceholder: '请输入密码',
      captchaPlaceholder: '请输入验证码',
      submit: '登 录',
      emailFormat: '请输入正确的邮箱格式',
      passwordLength: '密码的长度应为8~20位',
      captchaLength: '请输入6位的验证码'
    },
    register: {
      username: '用户名',
      password: '密码',
      confirmPassword: '确认密码',
      email: '邮箱',
      emailCaptcha: '邮箱验证码',
      usernamePlaceholder: '请输入用户名',
      passwordPlaceholder: '请输入密码',
      confirmPasswordPlaceholder: '请再次输入密码',
      emailPlaceholder: '请输入邮箱',
      imageCaptchaPlaceholder: '请输入图片验证码',
      emailCaptchaPlaceholder: '请输入邮箱验证码',
      passwordMismatch: '两次密码不一致！',
      sendCode: '发送验证码',
      submit: '注册',
      usernameMax: '用户名长度不应大于20位',
      passwordLength: '密码的长度应为8~20位',
      emailFormat: '请输入正确的邮箱格式',
      codeLength: '请输入6位的验证码'
    },
    forgotPassword: {
      email: '邮箱',
      emailCaptcha: '邮箱验证码',
      password: '密码',
      confirmPassword: '确认密码',
      emailPlaceholder: '请输入邮箱',
      imageCaptchaPlaceholder: '请输入图片验证码',
      emailCaptchaPlaceholder: '请输入邮箱验证码',
      newPasswordPlaceholder: '请输入新密码',
      confirmNewPasswordPlaceholder: '请再次输入新密码',
      passwordMismatch: '两次密码不一致！',
      sendCode: '发送验证码',
      submit: '确定',
      emailFormat: '请输入正确的邮箱格式',
      codeLength: '请输入6位的验证码',
      passwordLength: '密码的长度应为8~20位'
    },
    passwordReset: {
      oldPassword: '原密码',
      newPassword: '新密码',
      oldPasswordPlaceholder: '请输入旧密码',
      newPasswordPlaceholder: '请输入新密码',
      passwordLength: '密码的长度应为8~20位'
    },
    articleCreate: {
      cover: '文章封面',
      title: '文章标题',
      category: '文章类别',
      tags: '文章标签',
      abstract: '文章简介',
      titlePlaceholder: '请输入文章标题',
      categoryPlaceholder: '选择分类',
      tagPlaceholder: '选择标签',
      tagPlaceholderNoCategory: '请先选择分类',
      abstractPlaceholder: '请输入文章简介',
      dragPrefix: '将文件拖到此处，或',
      dragClick: '点击上传',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超过 20MB。'
    },
    articleUpdate: {
      cover: '文章封面',
      title: '文章标题',
      category: '文章类别',
      tags: '文章标签',
      abstract: '文章简介',
      content: '文章内容',
      titlePlaceholder: '请输入文章标题',
      categoryPlaceholder: '选择分类',
      tagPlaceholder: '选择标签',
      tagPlaceholderNoCategory: '请先选择分类',
      abstractPlaceholder: '请输入文章简介',
      editContent: '编辑内容',
      drawerFooter: '点击上方X或外部任意区域即可退出编辑',
      dragPrefix: '将文件拖到此处，或',
      dragClick: '点击上传',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超过 20MB。'
    },
    advertisementCreate: {
      image: '广告图片',
      link: '广告链接',
      title: '广告标题',
      content: '广告内容',
      linkPlaceholder: '请输入广告链接',
      titlePlaceholder: '请输入广告标题',
      contentPlaceholder: '请输入广告内容',
      dragPrefix: '将文件拖到此处，或',
      dragClick: '点击上传',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超过 20MB。'
    },
    advertisementUpdate: {
      image: '广告图片',
      link: '广告链接',
      title: '广告标题',
      content: '广告内容',
      linkPlaceholder: '请输入广告链接',
      titlePlaceholder: '请输入广告标题',
      contentPlaceholder: '请输入广告内容'
    },
    friendLinkCreate: {
      logo: 'logo图片',
      link: '友链链接',
      name: '友链名称',
      description: '友链描述',
      linkPlaceholder: '请输入友链链接',
      namePlaceholder: '请输入友链名称',
      descriptionPlaceholder: '请输入友链描述',
      dragPrefix: '将文件拖到此处，或',
      dragClick: '点击上传',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp 格式，大小不超过 20MB。'
    },
    friendLinkUpdate: {
      logo: 'logo图片',
      link: '友链链接',
      name: '友链名称',
      description: '友链描述',
      linkPlaceholder: '请输入友链链接',
      namePlaceholder: '请输入友链名称',
      descriptionPlaceholder: '请输入友链描述'
    },
    feedbackReply: {
      reply: '反馈回复',
      placeholder: '请输入反馈回复'
    }
  }
}
