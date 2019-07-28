文件夹 PATH 列表
卷序列号为 BAB0-9040
C:.
│  go.mod
│  go.sum
│  readme.md
│  tree.md
│  
├─cache
│  └─redis
│          conn.go
│          
├─common
│      code.go
│      
├─config
│      db.go
│      oss.go
│      rabiitmq.go
│      redis.go
│      service.go
│      user.go
│      
├─config-example
│      db.go
│      oss.go
│      rabiitmq.go
│      redis.go
│      service.go
│      user.go
│      
├─db
│  │  file.go
│  │  user.go
│  │  userfile.go
│  │  
│  └─mysql
│          conn.go
│          
├─doc
│      ci-cd.png
│      structure.png
│      table.sql
│      
├─handler
│  │  auth.go
│  │  mpupload.go
│  │  upload.go
│  │  user.go
│  │  
│  └─Gin-handler
│          auth.go
│          mpupload.go
│          upload.go
│          user.go
│          
├─meta
│      filemeta.go
│      
├─mq
│      consumer.go
│      define.go
│      producer.go
│      
├─route
│      router.go
│      
├─service
│  ├─Gin
│  │      main.go
│  │      
│  ├─Microservice
│  │  ├─account
│  │  │  │  main.go
│  │  │  │  
│  │  │  ├─handler
│  │  │  │      user.go
│  │  │  │      
│  │  │  └─proto
│  │  │          user.micro.go
│  │  │          user.pb.go
│  │  │          user.proto
│  │  │          
│  │  └─apigw
│  │      └─handler
│  │              user.go
│  │              
│  └─normal
│      ├─transfer
│      │      main.go
│      │      
│      └─upload
│              main.go
│              
├─static
│  ├─css
│  │      bootstrap.min.css
│  │      fileinput.min.css
│  │      
│  ├─img
│  │      avatar.jpeg
│  │      loading.gif
│  │      
│  ├─js
│  │      auth.js
│  │      bootstrap.min.js
│  │      fileinput.min.js
│  │      FileSaver.js
│  │      jquery-3.2.1.min.js
│  │      layui.js
│  │      piexif.min.js
│  │      polyfill.min.js
│  │      popper.min.js
│  │      purify.min.js
│  │      sortable.min.js
│  │      StreamSaver.js
│  │      sw.js
│  │      theme.js
│  │      
│  ├─log
│  └─view
│          download.html
│          home.html
│          signin.html
│          signup.html
│          upload.html
│          
├─store
│  └─oss
│          oss_conn.go
│          
├─tempFiles
│      readme.md
│      TIM图片20190704175044.jpg
│      
├─test
│      test_ceph.go
│      test_mpupload.go
│      
└─util
        resp.go
        util.go
        
