export default {
  system: {
    advertisement: {
      title: 'Advertisements',
      desc: 'Manage site ad slots and placements',
      create: 'New Ad',
      delete: 'Delete Ad',
      update: 'Update Ad',
      adTitle: 'Ad Title',
      adTitlePlaceholder: 'Enter ad title',
      adContent: 'Ad Content',
      adContentPlaceholder: 'Enter ad content',
      image: 'Image',
      link: 'Link',
      updateAction: 'Update',
      confirmDeleteWithCount: 'You have selected [{count}] item(s). This cannot be undone. Confirm deletion?'
    },
    appConfig: {
      site: 'Website',
      system: 'System',
      email: 'Email',
      qiniu: 'Qiniu',
      jwt: 'JWT',
      gaode: 'Amap'
    },
    feedback: {
      title: 'Feedback',
      desc: 'View and reply to user feedback',
      delete: 'Delete Feedback',
      user: 'User',
      reply: 'Reply',
      replyTitle: 'Reply to Feedback',
      confirmDeleteWithCount: 'You have selected [{count}] item(s). This cannot be undone. Confirm deletion?'
    },
    friendLink: {
      title: 'Friend Links',
      desc: "Manage outbound friend links shown on this site",
      create: 'New Link',
      delete: 'Delete Link',
      update: 'Update Link',
      name: 'Link Name',
      namePlaceholder: 'Enter link name',
      description: 'Description',
      descriptionPlaceholder: 'Enter link description',
      link: 'Link',
      descCol: 'Description',
      updateAction: 'Update',
      confirmDeleteWithCount: 'You have selected [{count}] item(s). This cannot be undone. Confirm deletion?'
    },
    loginLogs: {
      title: 'Login Logs',
      desc: 'View user login records and login status',
      uuidPlaceholder: 'Enter user UUID',
      user: 'User',
      username: 'Username',
      loginTime: 'Login Time',
      loginMethod: 'Method',
      loginAddress: 'IP Location',
      os: 'OS',
      deviceInfo: 'Device',
      browserInfo: 'Browser',
      loginStatus: 'Status'
    },
    config: {
      email: {
        title: 'Email',
        host: 'SMTP Host',
        port: 'SMTP Port',
        from: 'Sender Email',
        nickname: 'Sender Name',
        secret: 'SMTP Secret',
        ssl: 'Use SSL'
      },
      gaode: {
        title: 'Amap',
        enable: 'Enable',
        key: 'Amap Key'
      },
      jwt: {
        title: 'JWT',
        accessTokenSecret: 'Access Token Secret',
        accessTokenExpiry: 'Access Token Expiry',
        refreshTokenSecret: 'Refresh Token Secret',
        refreshTokenExpiry: 'Refresh Token Expiry',
        issuer: 'Issuer'
      },
      qiniu: {
        title: 'Qiniu',
        zone: 'Zone',
        bucket: 'Bucket',
        accessKey: 'Access Key (AK)',
        secretKey: 'Secret Key (SK)',
        cdnDomain: 'CDN Domain',
        useCdn: 'Use CDN Upload Acceleration',
        useHttps: 'Use HTTPS'
      },
      qq: {},
      site: {
        websiteInfo: 'Website',
        logo: 'Logo',
        fullLogo: 'Full Logo',
        siteTitle: 'Site Title',
        slogan: 'Slogan',
        sloganEn: 'English Slogan',
        description: 'Description',
        version: 'Version',
        createdAt: 'Created At',
        icp: 'ICP Filing',
        police: 'Public Security Filing',
        bilibili: 'Bilibili URL',
        gitee: 'Gitee URL',
        github: 'GitHub URL',
        personalInfo: 'Profile',
        nickname: 'Nickname',
        job: 'Job',
        address: 'Location',
        email: 'Email',
        qqImage: 'QQ Image',
        wechatImage: 'WeChat Image',
        footerLink: 'Footer Links',
        create: 'New',
        title: 'Title',
        link: 'Link',
        homeImage: 'Home Banner',
        imageUpload: 'Upload Image',
        clearLogoConfirm: 'Clear the Logo image?',
        clearFullLogoConfirm: 'Clear the Full Logo image?',
        clearQqImageConfirm: 'Clear the QQ image?',
        clearWechatImageConfirm: 'Clear the WeChat image?',
        removeCarouselConfirm: 'Remove this home banner?',
        operationCancelled: 'Operation cancelled',
        warning: 'Warning',
        dragPrefix: 'Drop file here or ',
        dragClick: 'click to upload',
        uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp files, max size 20MB.'
      },
      system: {
        title: 'System',
        multipoint: 'Multi-session Login Block',
        sessionSecret: 'Session Secret',
        ossType: 'Image Storage',
        selectPlaceholder: 'Select',
        ossLocal: 'Local',
        ossQiniu: 'Qiniu'
      }
    },
    users: {
      title: 'Users',
      desc: 'View registered users and manage freeze status',
      uuidPlaceholder: 'Enter user UUID',
      registerSource: 'Register Source',
      selectPlaceholder: 'Select',
      registerEmail: 'Email',
      avatar: 'Avatar',
      username: 'Username',
      address: 'Location',
      registerTime: 'Register Time',
      role: 'Role',
      admin: 'Admin',
      normalUser: 'Normal User',
      freeze: 'Freeze',
      unfreeze: 'Unfreeze',
      freezeDialogTitle: 'Freeze / Unfreeze User',
      freezeConfirm: '{action} this user: {username}'
    },
    userCenter: {
      comment: {
        title: 'My Comments',
        desc: 'View all comments I have posted'
      },
      feedback: {
        title: 'My Feedback',
        desc: 'View my feedback and admin replies',
        reply: 'Reply'
      },
      info: {
        title: 'Profile',
        userCard: 'User Card',
        avatar: 'Avatar',
        username: 'Username',
        address: 'Location',
        signature: 'Bio',
        email: 'Email',
        role: 'Role',
        registerSource: 'Register Source',
        normalUser: 'Normal User',
        admin: 'Admin',
        operation: 'Actions',
        changePassword: 'Change Password',
        usernameMax: 'Username must be at most 20 characters',
        addressMax: 'Location must be at most 200 characters',
        signatureMax: 'Bio must be at most 320 characters'
      },
      star: {
        title: 'Favorites',
        desc: 'View articles I have bookmarked',
        category: 'Category',
        abstract: 'Summary',
        publishTime: 'Published At',
        articleId: 'Article ID'
      }
    },
    login: {
      title: 'Login'
    },
    error: {
      notFound: 'Sorry, the page you visited does not exist or has been removed.',
      backHome: 'Back to Home'
    }
  },
  forms: {
    login: {
      email: 'Email',
      password: 'Password',
      captcha: 'Captcha',
      emailPlaceholder: 'Enter your email',
      passwordPlaceholder: 'Enter your password',
      captchaPlaceholder: 'Enter the captcha',
      submit: 'Login',
      emailFormat: 'Please enter a valid email address',
      passwordLength: 'Password must be 8-20 characters',
      captchaLength: 'Please enter the 6-digit captcha'
    },
    register: {
      username: 'Username',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      email: 'Email',
      emailCaptcha: 'Email Code',
      usernamePlaceholder: 'Enter username',
      passwordPlaceholder: 'Enter password',
      confirmPasswordPlaceholder: 'Re-enter password',
      emailPlaceholder: 'Enter your email',
      imageCaptchaPlaceholder: 'Enter image captcha',
      emailCaptchaPlaceholder: 'Enter email verification code',
      passwordMismatch: 'Passwords do not match!',
      sendCode: 'Send Code',
      submit: 'Register',
      usernameMax: 'Username must be at most 20 characters',
      passwordLength: 'Password must be 8-20 characters',
      emailFormat: 'Please enter a valid email address',
      codeLength: 'Please enter the 6-digit code'
    },
    forgotPassword: {
      email: 'Email',
      emailCaptcha: 'Email Code',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      emailPlaceholder: 'Enter your email',
      imageCaptchaPlaceholder: 'Enter image captcha',
      emailCaptchaPlaceholder: 'Enter email verification code',
      newPasswordPlaceholder: 'Enter new password',
      confirmNewPasswordPlaceholder: 'Re-enter new password',
      passwordMismatch: 'Passwords do not match!',
      sendCode: 'Send Code',
      submit: 'Confirm',
      emailFormat: 'Please enter a valid email address',
      codeLength: 'Please enter the 6-digit code',
      passwordLength: 'Password must be 8-20 characters'
    },
    passwordReset: {
      oldPassword: 'Current Password',
      newPassword: 'New Password',
      oldPasswordPlaceholder: 'Enter current password',
      newPasswordPlaceholder: 'Enter new password',
      passwordLength: 'Password must be 8-20 characters'
    },
    articleCreate: {
      cover: 'Cover',
      title: 'Title',
      category: 'Category',
      tags: 'Tags',
      abstract: 'Summary',
      titlePlaceholder: 'Enter article title',
      categoryPlaceholder: 'Select category',
      tagPlaceholder: 'Select tags',
      tagPlaceholderNoCategory: 'Select a category first',
      abstractPlaceholder: 'Enter article summary',
      dragPrefix: 'Drop file here or ',
      dragClick: 'click to upload',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp files, max size 20MB.'
    },
    articleUpdate: {
      cover: 'Cover',
      title: 'Title',
      category: 'Category',
      tags: 'Tags',
      abstract: 'Summary',
      content: 'Content',
      titlePlaceholder: 'Enter article title',
      categoryPlaceholder: 'Select category',
      tagPlaceholder: 'Select tags',
      tagPlaceholderNoCategory: 'Select a category first',
      abstractPlaceholder: 'Enter article summary',
      editContent: 'Edit Content',
      drawerFooter: 'Click the X above or anywhere outside to exit editing',
      dragPrefix: 'Drop file here or ',
      dragClick: 'click to upload',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp files, max size 20MB.'
    },
    advertisementCreate: {
      image: 'Ad Image',
      link: 'Ad Link',
      title: 'Ad Title',
      content: 'Ad Content',
      linkPlaceholder: 'Enter ad link',
      titlePlaceholder: 'Enter ad title',
      contentPlaceholder: 'Enter ad content',
      dragPrefix: 'Drop file here or ',
      dragClick: 'click to upload',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp files, max size 20MB.'
    },
    advertisementUpdate: {
      image: 'Ad Image',
      link: 'Ad Link',
      title: 'Ad Title',
      content: 'Ad Content',
      linkPlaceholder: 'Enter ad link',
      titlePlaceholder: 'Enter ad title',
      contentPlaceholder: 'Enter ad content'
    },
    friendLinkCreate: {
      logo: 'Logo',
      link: 'Link',
      name: 'Link Name',
      description: 'Description',
      linkPlaceholder: 'Enter link URL',
      namePlaceholder: 'Enter link name',
      descriptionPlaceholder: 'Enter link description',
      dragPrefix: 'Drop file here or ',
      dragClick: 'click to upload',
      uploadTip: 'jpg/png/jpeg/ico/tiff/gif/svg/webp files, max size 20MB.'
    },
    friendLinkUpdate: {
      logo: 'Logo',
      link: 'Link',
      name: 'Link Name',
      description: 'Description',
      linkPlaceholder: 'Enter link URL',
      namePlaceholder: 'Enter link name',
      descriptionPlaceholder: 'Enter link description'
    },
    feedbackReply: {
      reply: 'Reply',
      placeholder: 'Enter your reply'
    }
  }
}
